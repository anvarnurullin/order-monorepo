# Order Management System

Система управления заказами.

## Запуск системы

1. Запуск всех сервисов:

```bash
make up
```

2. Применение миграций (после первого запуска):

```bash
make migrate-up
```

3. Пересборка после изменений:

```bash
make rebuild
```

4. Остановка системы:

```bash
make down
```

## Сервисы

- **Frontend**: http://localhost
- **Catalog API**: http://localhost/api/v1/products
- **Order API**: http://localhost/api/v1/orders
- **Auth API**: http://localhost/api/v1/auth

## Структура проекта

- `services/catalog/` - Сервис каталога товаров
- `services/order/` - Сервис заказов
- `services/auth/` - Сервис авторизации
- `frontend/` - React фронтенд
- `migrations/` - SQL миграции
- `docker-compose.yml` - Конфигурация контейнеров
