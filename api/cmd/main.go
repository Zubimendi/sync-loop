package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/Zubimendi/sync-loop/api/internal/auth"
	"github.com/Zubimendi/sync-loop/api/internal/handler"
	"github.com/Zubimendi/sync-loop/api/internal/middleware"
	"github.com/Zubimendi/sync-loop/api/internal/repo"
	"github.com/Zubimendi/sync-loop/api/internal/connector"
	"github.com/Zubimendi/sync-loop/api/internal/job"
	"github.com/Zubimendi/sync-loop/api/internal/temporal"
	"github.com/rs/cors"
)

func main() {
	log.Info().Str("JWT_SECRET", os.Getenv("JWT_SECRET")).Msg("env check")
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("db connect")
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("db ping")
	}
	log.Info().Msg("postgres connected")

	userRepo := repo.NewUserRepo(db)
	authSvc  := auth.NewService(userRepo)
	authH    := handler.NewAuthHandler(authSvc)
	
	connRepo := connector.NewRepo(db)
	connSvc  := connector.NewService(connRepo)
	connH    := connector.NewHandler(connSvc)

	c := cors.New(cors.Options{
	AllowedOrigins:   []string{"http://localhost:3000"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
})
	
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(c.Handler)
	authMw := middleware.Auth(authSvc)

	if err := temporal.Init(); err != nil {
		log.Fatal().Err(err).Msg("temporal init")
	}
	defer temporal.Close()
	jobH := job.NewHandler(temporal.DefaultClient)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", authH.Register)
		r.Post("/login", authH.Login)
		r.Post("/logout", authH.Logout)

		r.Group(func(r chi.Router) {
			r.Use(authMw)
			r.Get("/me", authH.Me)
			r.Get("/connectors", connH.List)
			r.Post("/connectors", connH.Create)
			r.Get("/jobs", jobH.List)
			r.Post("/jobs/run-now", jobH.RunNow)
			r.Post("/jobs/cancel", jobH.Cancel)
			r.Post("/jobs/terminate-all", jobH.TerminateAll)

			// NEW FEATURES:
			r.Post("/jobs/retry", jobH.Retry)                    // Retry failed workflow
			r.Post("/jobs/schedule", jobH.CreateSchedule)        // Create/update schedule
			r.Post("/jobs/schedule/toggle", jobH.ToggleSchedule) // Pause/unpause schedule
			r.Get("/jobs/{id}/status", jobH.GetJobStatus)
		})
	})

	log.Info().Msg("api listening :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal().Err(err).Msg("api failed")
	}
}
