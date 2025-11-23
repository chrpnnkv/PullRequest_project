package api

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/team/add", h.handleTeamAdd)
	mux.HandleFunc("/team/get", h.handleTeamGet)

	mux.HandleFunc("/users/setIsActive", h.handleUserSetIsActive)
	mux.HandleFunc("/users/getReview", h.handleGetReviews)

	mux.HandleFunc("/pullRequest/create", h.handlePrAdd)
	mux.HandleFunc("/pullRequest/merge", h.handleMerge)
	mux.HandleFunc("/pullRequest/reassign", h.handleReassign)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
