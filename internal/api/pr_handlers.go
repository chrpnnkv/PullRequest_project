package api

import (
	"PR_project/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func (h *Handler) handlePrAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only POST is allowed")
		return
	}

	var req PrReqDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if req.PrID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "pull_request_id is required")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "pull_request_name is required")
		return
	}
	if req.AuthorID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "author_id is required")
		return
	}

	ctx := r.Context()
	pr, err := h.PrService.CreatePRWithReviewers(ctx, req.PrID, req.Name, req.AuthorID)
	if errors.Is(err, service.ErrPrAlreadyExists) {
		writeError(w, http.StatusConflict, "PR_EXISTS", "PR id already exists")
		return
	}
	if errors.Is(err, service.ErrUserNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", "internal server error")
		return
	}

	respPr := PrDTO{
		PrID:         pr.ID,
		Name:         pr.Name,
		AuthorID:     pr.AuthorID,
		Status:       pr.Status,
		CreatedAt:    pr.CreatedAt,
		MergedAt:     pr.MergedAt,
		ReviewersIDs: pr.ReviewersIDs,
	}

	resp := PrAddResponse{
		Pr: respPr,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only POST is allowed")
		return
	}

	var req PrReqDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if req.PrID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "pull_request_id is required")
		return
	}

	ctx := r.Context()
	pr, err := h.PrService.Merge(ctx, req.PrID, time.Now())
	if errors.Is(err, service.ErrPrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "pr not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", "internal server error")
		return
	}

	respPr := PrDTO{
		PrID:         pr.ID,
		Name:         pr.Name,
		AuthorID:     pr.AuthorID,
		Status:       pr.Status,
		CreatedAt:    pr.CreatedAt,
		MergedAt:     pr.MergedAt,
		ReviewersIDs: pr.ReviewersIDs,
	}

	resp := PrAddResponse{
		Pr: respPr,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleReassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only POST is allowed")
		return
	}

	var req PrReassignReqDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if req.PrID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "pull_request_id is required")
		return
	}
	if req.OldReviewer == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "old_reviewer_id is required")
		return
	}

	ctx := r.Context()
	pr, replacedBy, err := h.PrService.Reassign(ctx, req.PrID, req.OldReviewer)
	if errors.Is(err, service.ErrPrNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "pr not found")
		return
	}
	if errors.Is(err, service.ErrPrAlreadyMerged) {
		writeError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
		return
	}
	if errors.Is(err, service.ErrUserIsNotReviewer) {
		writeError(w, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
		return
	}
	if errors.Is(err, service.ErrNoCandidate) {
		writeError(w, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", "internal server error")
		return
	}

	respPr := PrDTO{
		PrID:         pr.ID,
		Name:         pr.Name,
		AuthorID:     pr.AuthorID,
		Status:       pr.Status,
		CreatedAt:    pr.CreatedAt,
		MergedAt:     pr.MergedAt,
		ReviewersIDs: pr.ReviewersIDs,
	}

	resp := PrReassignedResponse{
		Pr:         respPr,
		ReplacedBy: replacedBy,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
