package teamhandler

import (
	"encoding/json"
	"net/http"
	"AvitoProject/utils"
	modelteam "AvitoProject/internal/models/team"
	"AvitoProject/internal/usecase/errors"
)

type teamUsecaseI interface{
	CreateTeam(t *modelteam.Team) error
	GetTeam(name string) (*modelteam.Team, error)
}

type Handler struct{
	teamUsecase teamUsecaseI
}

func New(t teamUsecaseI) *Handler{
	return &Handler{
		teamUsecase: t,
	}
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req modelteam.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "BAD_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.teamUsecase.CreateTeam(&req)
	if err != nil {
		if err == errors.ErrTeamExists {
			utils.WriteError(w, "TEAM_EXISTS", "team_name already exists", http.StatusBadRequest)
			return
		}
		utils.WriteError(w, "INTERNAL", "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"team": req})
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("team_name")
	if name == "" {
		utils.WriteError(w, "BAD_REQUEST", "team_name is required", http.StatusBadRequest)
		return
	}

	t, err := h.teamUsecase.GetTeam(name)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.WriteError(w, "NOT_FOUND", "resource not found", http.StatusNotFound)
			return
		}
		utils.WriteError(w, "INTERNAL", "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}