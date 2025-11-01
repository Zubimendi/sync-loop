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
	"github.com/Zubimendi/sync-loop/api/internal/middleware" // our package
	"github.com/Zubimendi/sync-loop/api/internal/repo"
)

func main() {
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

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	authMw := middleware.Auth(authSvc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", authH.Register)
				r.Post("/login", authH.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMw)
			r.Get("/me", authH.Me) // we’ll add this next
		})
	})

	log.Info().Msg("api listening :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal().Err(err).Msg("api failed")
	}
}
