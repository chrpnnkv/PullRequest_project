package api

import "PR_project/internal/service"

type Handler struct {
	TeamService *service.TeamService
	UserService *service.UserService
	PrService   *service.PrService
}

func NewHandler(
	teamService *service.TeamService,
	userService *service.UserService,
	prService *service.PrService,
) *Handler {
	return &Handler{
		TeamService: teamService,
		UserService: userService,
		PrService:   prService,
	}
}
