# Carrega variáveis do .env
include .env
export

.PHONY: up down logs

# Sobe a infra (banco + redis) em background
up:
	docker compose up -d

# Derruba tudo
down:
	docker compose down

# Vê logs dos containers
logs:
	docker compose logs -f

# Roda a aplicação Go
run:
	go run cmd/api/main.go