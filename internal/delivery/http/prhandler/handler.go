package prhandler

import (
	"encoding/json"
	"net/http"
	"AvitoProject/utils"
	modelpr "AvitoProject/internal/models/pr"
	modeluser "AvitoProject/internal/models/user"
	"AvitoProject/internal/usecase/errors"
)

type prUsecaseI interface {
	CreatePR(prID modelpr.ID, name string, authorID modeluser.ID) (*modelpr.PullRequest, error)
	MergePR(prID modelpr.ID) (*modelpr.PullRequest, error)
	ReassignReviewer(prID modelpr.ID, oldUserID modeluser.ID) (*modelpr.PullRequest, modeluser.ID, error)
}

type Handler struct {
	prUsecase prUsecaseI
}

func New(p prUsecaseI) *Handler {
	return &Handler{
		prUsecase: p,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "BAD_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	prObj, err := h.prUsecase.CreatePR(
		modelpr.ID(req.PullRequestID),
		req.PullRequestName,
		modeluser.ID(req.AuthorID),
	)
	if err != nil {
		switch err {
		case errors.ErrPRExists:
			utils.WriteError(w, "PR_EXISTS", "PR id already exists", http.StatusConflict)
			return
		case errors.ErrNotFound:
			utils.WriteError(w, "NOT_FOUND", "author or team not found", http.StatusNotFound)
			return
		default:
			utils.WriteError(w, "INTERNAL", "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"pr": prObj})
}

func (h *Handler) Merge(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "BAD_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	prObj, err := h.prUsecase.MergePR(modelpr.ID(req.PullRequestID))
	if err != nil {
		if err == errors.ErrNotFound {
			utils.WriteError(w, "NOT_FOUND", "PR not found", http.StatusNotFound)
			return
		}
		utils.WriteError(w, "INTERNAL", "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"pr": prObj})
}

func (h *Handler) Reassign(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldReviewerID string `json:"old_user_id"` // В OpenAPI: old_user_id
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, "BAD_REQUEST", "invalid JSON", http.StatusBadRequest)
		return
	}

	prObj, newReviewer, err := h.prUsecase.ReassignReviewer(
		modelpr.ID(req.PullRequestID),
		modeluser.ID(req.OldReviewerID),
	)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			utils.WriteError(w, "NOT_FOUND", "PR or user not found", http.StatusNotFound)
			return
		case errors.ErrPRMerged:
			utils.WriteError(w, "PR_MERGED", "cannot reassign on merged PR", http.StatusConflict)
			return
		case errors.ErrNotAssigned:
			utils.WriteError(w, "NOT_ASSIGNED", "reviewer is not assigned to this PR", http.StatusConflict)
			return
		case errors.ErrNoCandidate:
			utils.WriteError(w, "NO_CANDIDATE", "no active replacement candidate in team", http.StatusConflict)
			return
		default:
			utils.WriteError(w, "INTERNAL", "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr":          prObj,
		"replaced_by": newReviewer,
	})
}