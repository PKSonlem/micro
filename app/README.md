## Настройка и запуск проекта

### Переменные окружения
Приложение загружает переменные из `deploy/local/.env`. Создайте этот файл со следующими переменными:

```bash
# PostgreSQL Configuration
POSTGRES_DB=mentorship_db
POSTGRES_HOST=db  # для Docker, или localhost для локальной разработки
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres_password
POSTGRES_SSL_MODE=disable
POSTGRES_MAX_IDLE_CONNECTIONS=10
POSTGRES_MAX_OPEN_CONNECTIONS=10

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

### Запуск через Docker Compose
```bash
cd app
make up  # или docker-compose -f deploy/local/docker-compose.yml up -d
```

### Интеграционные тесты
Для тестов скопируйте `deploy/local/.test.env.example` в `deploy/local/.test.env` и выполните:
```bash
make test-integration
```