
DATE ?= $(shell date +%FT%T%z)
BUILD_VERSION ?= '1.0.0'
BINARY=build/main

MIGRATION_PATH=db/migrations

# =========================
# Migration commands
# =========================

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

# Create new migration (you provide name)
# Usage: make create-migration name=create_users_table
create-migration:
	migrate create -ext sql -dir "$(MIGRATION_PATH)" -seq $(name)

# =========================
# DB reset
# =========================

# Drop all tables
db-drop:
	migrate -database "$(DB_URL)" -path "$(MIGRATION_PATH)" force 0
	migrate -database "$(DB_URL)" -path "$(MIGRATION_PATH)" down

.PHONY: run build tidy
	
run:
	go run ./cmd/api/main.go

build:
	go build -o bin/server ./cmd/main.go


tidy:
	go mod tidy

.PHONY: configure
configure:
	go mod download && go mod tidy

.PHONY: build.alpine
build.alpine:
	go build -tags musl \
		-ldflags " \
		-X github.com/clientportal-api.BuildVersion=${BUILD_VERSION} \
		-X github.com/clientportal-api.BuildDate=${DATE}" \
		-o $(BINARY) cmd/main.go