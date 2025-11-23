package api

type UserReqDTO struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserAndPrDTO struct {
	UserID string       `json:"user_id"`
	Prs    []PrShortDTO `json:"pull_requests"`
}

type UserAddResponse struct {
	User UserDTO `json:"user"`
}
