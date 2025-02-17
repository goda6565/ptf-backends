.PHONY: check
check:
	go fmt ./...
	golangci-lint run

.PHONY: test
test:
	go test -v ./...

.PHONY: db-up
db-up:
	docker-compose up --build -d

.PHONY: generate
generate:
	oapi-codegen --config=./api/config.yaml ./api/openapi.yaml

# atlas関連
.PHONY: inspect
inspect:
	atlas schema inspect --env gorm --url env://src -w

.PHONY: migrate
migrate:
	atlas migrate diff --env gorm

.PHONY: apply
apply:
	atlas migrate apply --url postgres://postgres:postgres@localhost:5432/auth?sslmode=disable

.PHONT: psql
psql:
	docker-compose exec postgres psql -U postgres -d auth