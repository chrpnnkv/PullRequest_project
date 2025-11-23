package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"PR_project/internal/api"
	"PR_project/internal/repository"
	"PR_project/internal/service"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB is not reachable:", err)
	}

	teamRepo := repository.NewPostgresTeamRepository(db)
	userRepo := repository.NewPostgresUserRepository(db)
	prRepo := repository.NewPostgresPrRepository(db)

	teamService := &service.TeamService{TRepository: teamRepo}
	userService := &service.UserService{URepository: userRepo}
	prService := &service.PrService{
		PrRepository: prRepo,
		URepository:  userRepo,
	}

	handler := api.Handler{
		TeamService: teamService,
		UserService: userService,
		PrService:   prService,
	}

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
