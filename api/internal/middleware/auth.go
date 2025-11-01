package middleware
import (
    "context"
    "net/http"
    "github.com/Zubimendi/sync-loop/api/internal/auth"
)
type ctxKey string
const CtxUserID = ctxKey("userID")
const CtxWorkspaceID = ctxKey("workspaceID")
func Auth(svc * auth.Service) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {
            c, err := r.Cookie("token")
            if err != nil {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            uid, wid, err := svc.ValidateToken(c.Value)
            if err != nil {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            ctx := context.WithValue(r.Context(), CtxUserID, uid)
            ctx = context.WithValue(ctx, CtxWorkspaceID, wid)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

