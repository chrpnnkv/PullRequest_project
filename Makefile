run:
	go run ./cmd/app

build:
	go build -o server ./cmd/app

docker-up:
	docker compose up --build

docker-down:
	docker compose down
