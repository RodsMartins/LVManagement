.PHONY: tailwind-watch
tailwind-watch:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: templ-watch
templ-watch:
	templ generate --watch
	
.PHONY: dev
dev:
	make containers
	go build -o ./tmp/$(APP_NAME) ./cmd/$(APP_NAME)/main.go && air

.PHONY: containers
containers:
	docker compose up -d

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.Environment=production" -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

PHONY: db-init
db-init:
	docker compose exec -T db psql -Upostgres -dlvm < schema.sql

.PHONY: db-reset
db-reset:
	docker compose exec -T db psql -Upostgres -dlvm -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	make db-init

.PHONY: test
test:
	go test -race -v -timeout 30s ./...