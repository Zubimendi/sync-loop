import (
"net/http"
"github.com/go-chi/chi/v5"
"github.com/go-chi/chi/v5/middleware"
"github.com/rs/zerolog/log"
)
func main() {
r := chi.NewRouter()
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
})

log.Info().Msg("api listening :8080")
http.ListenAndServe(":8080", r)
}