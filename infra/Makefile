# ==========================
# INFRASTRUCTURE
# ==========================
.PHONY: up down logs

# Поднять контейнеры (Postgres, Redis, Kafka, Zookeeper)
up:
	docker-compose up -d

# Остановить и удалить контейнеры + тома
down:
	docker-compose down -v

# Смотреть логи всех контейнеров
logs:
	docker-compose logs -f
