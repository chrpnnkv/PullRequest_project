# PR Reviewer

Сервис для назначения ревьюверов на Pull Request’ы на основе командной структуры.
Поддерживает команды, пользователей, автоматическое назначение ревьюверов, перенос ревью, деактивацию пользователей.

## Возможности сервиса Teams

* Создание команды с участниками
* Получение команды по имени
* Автоматическое создание/обновление пользователей при создании команды

## Возможности сервиса Users

* Деактивация/активация пользователя
* Получение PR-ов, где пользователь назначен ревьювером

## Возможности сервиса Pull Requests

* Создание PR, автоматическое назначение до 2 ревьюверов
* Идемпотентный merge
* Переназначение ревьювера внутри команды


## Запуск через Docker Compose

`docker compose up --build`


После запуска сервис доступен по адресу:

http://localhost:8080


## Миграции

При первом старте Postgres автоматически применяет:
`migrations/001_init.sql`

Это создаёт таблицы:

1. teams
2. users
3. pull_requests
4. pr_reviewers

## API Endpoints

Полностью соответствует OpenAPI (см. файл openapi.yml).

## Запуск без Docker
`go mod download` 
`go run ./cmd/app`

Требуется локальный PostgreSQL и корректный DB_DSN.

## Примеры запросов (PowerShell curl)

**Создать команду:**
curl -Method POST -Uri "http://localhost:8080/team/add" -ContentType "application/json" -Body '{"team_name":"backend","members":[{"user_id":"u1","username":"Alice","is_active":true}]}'
