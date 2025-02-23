ifneq (,$(wildcard .env))
    include .env
    export
endif

$(info POSTGRES_URL: $(POSTGRES_URL))

.PHONY: test coverage clean run compose down watch worker

test:
	ENV=test go test -coverprofile=cover.out -v ./...
coverage:
	go tool cover -html=cover.out
clean:
	rm main cover.out || true
	docker compose down --volumes --remove-orphans
down:
	docker compose down --remove-orphans
run:
	go run cmd/api/main.go
watch:
	air -c .air.toml
worker:
	go run cmd/worker/main.go
compose:
	docker compose -f compose.yml up -d --remove-orphans
start:
	docker compose -f app.compose.yml up -d --remove-orphans
stop:
	docker compose -f app.compose.yml down --remove-orphans
restart:
	docker compose -f app.compose.yml down --remove-orphans
	docker compose -f app.compose.yml up -d --remove-orphans
shutdown:
	docker compose -f app.compose.yml down --remove-orphans
	docker compose down --remove-orphans

#### DB ####
.PHONY: sqlc migrate-up migrate-down new-migration
sqlc:
	./scripts/sqlc-generate.sh

migrate-up:
	migrate -path pkg/db/migration -database "$(POSTGRES_URL)" -verbose up

migrate-down:
	migrate -path pkg/db/migration -database "$(POSTGRES_URL)" -verbose down
	
new-migration:
	migrate create -ext sql -dir pkg/db/migration -seq $(name)


