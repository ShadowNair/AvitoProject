package userhandler

import (
	"encoding/json"
	"net/http"
	modeluser "AvitoProject/internal/models/user"
	modelpr "AvitoProject/internal/models/pr"
	"AvitoProject/internal/usecase/errors"
	"AvitoProject/utils"
)

type userUsecaseI interface {
	SetUserActive(userID modeluser.ID, active bool) (*modeluser.User, error)
	GetReviewPRs(userID modeluser.ID) ([]modelpr.PullRequestShort, error)
}

type Handler struct {
	userUsecase userUsecaseI
}

func New(u userUsecaseI) *Handler {
	return &Handler{
		userUsecase: u,
	}
}

func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "BAD_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	usr, err := h.userUsecase.SetUserActive(modeluser.ID(req.UserID), req.IsActive)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.WriteError(w, "NOT_FOUND", "user not found", http.StatusNotFound)
			return
		}
		utils.WriteError(w, "INTERNAL", "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"user": usr})
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		utils.WriteError(w, "BAD_REQUEST", "user_id is required", http.StatusBadRequest)
		return
	}

	prs, err := h.userUsecase.GetReviewPRs(modeluser.ID(userID))
	if err != nil {
		utils.WriteError(w, "INTERNAL", "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}