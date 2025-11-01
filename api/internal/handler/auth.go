package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Zubimendi/sync-loop/api/internal/auth"
	"github.com/Zubimendi/sync-loop/api/internal/middleware"
)

type AuthHandler struct {
	svc *auth.Service
}

func NewAuthHandler(svc *auth.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type registerReq struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	WorkspaceName string `json:"workspace_name"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	u, wsp, token, err := h.svc.Register(r.Context(), req.Email, req.Password, req.WorkspaceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // true in prod with HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   86400,
	})
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":      u,
		"workspace": wsp,
	})
}

type loginReq struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

func(h * AuthHandler) Login(w http.ResponseWriter, r * http.Request) {
    var req loginReq
    if err := json.NewDecoder(r.Body).Decode( & req);
    err != nil {
        http.Error(w, "invalid body", http.StatusBadRequest)
        return
    }
    u, wsp, token, err := h.svc.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    http.SetCookie(w, & http.Cookie {
        Name: "token",
        Value: token,
        HttpOnly: true,
        Secure: false,
        SameSite: http.SameSiteLaxMode,
        Path: "/",
        MaxAge: 86400,
    })
    json.NewEncoder(w).Encode(map[string] interface {} {
        "user": u,
        "workspace": wsp,
    })
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.CtxUserID).(string)
	wid := r.Context().Value(middleware.CtxWorkspaceID).(string)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":      uid,
		"workspace_id": wid,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusNoContent)
}