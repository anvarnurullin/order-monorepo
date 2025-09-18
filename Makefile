# ==========================
# INFRASTRUCTURE
# ==========================
.PHONY: up down rebuild logs

# Поднять контейнеры (Postgres, Redis, Kafka, Zookeeper)
up:
	docker-compose up -d

# Остановить и удалить контейнеры
down:
	docker-compose down

# Пересобрать и перезапустить сервисы
rebuild:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Смотреть логи всех контейнеров
logs:
	docker-compose logs -f


# ==========================
# DATABASE MIGRATIONS
# ==========================
.PHONY: migrate-up migrate-down

# Применить миграции
migrate-up:
	docker-compose exec postgres psql -U app -d app -f /docker-entrypoint-initdb.d/0001_create_products.up.sql
	docker-compose exec postgres psql -U app -d app -f /docker-entrypoint-initdb.d/0002_create_orders.up.sql
	docker-compose exec postgres psql -U app -d app -f /docker-entrypoint-initdb.d/0003_create_users.up.sql
	docker-compose exec postgres psql -U app -d app -f /docker-entrypoint-initdb.d/seed_products.sql
	docker-compose exec postgres psql -U app -d app -f /docker-entrypoint-initdb.d/seed_orders.sql

# Откатить миграции
migrate-down:
	docker-compose exec postgres psql -U app -d app -f /docker-entrypoint-initdb.d/0003_create_users.down.sql
	docker-compose exec postgres psql -U app -d app -f /docker-entrypoint-initdb.d/0002_create_orders.down.sql
	docker-compose exec postgres psql -U app -d app -f /docker-entrypoint-initdb.d/0001_create_products.down.sql
