# ==========================
# INFRASTRUCTURE
# ==========================
.PHONY: up down logs

# Поднять контейнеры (Postgres, Redis, Kafka, Zookeeper)
up:
	docker-compose up -d

# Остановить и удалить контейнеры
down:
	docker-compose down

# Смотреть логи всех контейнеров
logs:
	docker-compose logs -f
