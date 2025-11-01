package connector

import (
	"encoding/json"
	"net/http"

	"github.com/Zubimendi/sync-loop/api/internal/middleware"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	wid := r.Context().Value(middleware.CtxWorkspaceID).(string)
	cc, err := h.svc.ListSources(r.Context(), wid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"connectors": cc})
}

type createReq struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	uid := r.Context().Value(middleware.CtxUserID).(string)
	wid := r.Context().Value(middleware.CtxWorkspaceID).(string)

	c, err := h.svc.CreateSource(r.Context(), req.Name, req.Type, req.Config, uid, wid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(c)
}