package api

import (
	"PR_project/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

func (h *Handler) handleTeamAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only POST is allowed")
		return
	}

	var req TeamDTO

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if req.TeamName == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "team_name is required")
		return
	}

	for _, member := range req.Members {
		if member.UserID == "" {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "member.user_id is required")
			return
		}
	}

	membersInput := make([]service.TeamMemberInput, 0, len(req.Members))
	for _, member := range req.Members {
		membersInput = append(membersInput, service.TeamMemberInput{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	ctx := r.Context()
	team, users, err := h.TeamService.CreateTeam(ctx, req.TeamName, membersInput)
	if errors.Is(err, service.ErrTeamAlreadyExists) {
		writeError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", "internal server error")
		return
	}

	membersDTO := make([]TeamMemberDTO, 0, len(users))
	for _, u := range users {
		membersDTO = append(membersDTO, TeamMemberDTO{
			UserID:   u.ID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	respTeam := TeamDTO{
		TeamName: team.Name,
		Members:  membersDTO,
	}

	resp := TeamAddResponse{
		Team: respTeam,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleTeamGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only GET is allowed")
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "team_name is required")
		return
	}

	ctx := r.Context()
	team, users, err := h.TeamService.GetTeam(ctx, teamName)
	if errors.Is(err, service.ErrTeamNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "team not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", "internal server error")
		return
	}

	membersDTO := make([]TeamMemberDTO, 0, len(users))
	for _, u := range users {
		membersDTO = append(membersDTO, TeamMemberDTO{
			UserID:   u.ID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	respTeam := TeamDTO{
		TeamName: team.Name,
		Members:  membersDTO,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respTeam)
}
