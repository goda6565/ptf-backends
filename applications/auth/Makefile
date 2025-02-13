.PHONY: check
check:
	go fmt ./...
	golangci-lint run

.PHONY: db-up
db-up:
	docker-compose up --build -d


# atlas関連
.PHONY: inspect
inspect:
	atlas schema inspect --env gorm --url env://src -w

.PHONY: migrate
migrate:
	atlas migrate diff --env gorm

.PHONY: apply
apply:
	atlas migrate apply --url postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable