### general tips

## Quick Start

1. Создайте файл `deploy/local/.env` со следующим содержимым:
```
POSTGRES_DB=mentorship_db
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres_password
POSTGRES_SSL_MODE=disable
POSTGRES_MAX_IDLE_CONNECTIONS=10
POSTGRES_MAX_OPEN_CONNECTIONS=10
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

2. Запустите проект: `make up`

## Environment Variables

Приложение использует следующие переменные окружения:

### Required:
- `POSTGRES_DB` - имя базы данных
- `POSTGRES_HOST` - хост PostgreSQL (для Docker: `db`, для локальной разработки: `localhost`)
- `POSTGRES_USER` - пользователь PostgreSQL
- `POSTGRES_PASSWORD` - пароль PostgreSQL

### Optional (с дефолтными значениями):
- `POSTGRES_PORT` - порт PostgreSQL (default: 5432)
- `POSTGRES_SSL_MODE` - режим SSL (default: disable)
- `POSTGRES_MAX_IDLE_CONNECTIONS` - макс. idle соединений (default: 10)
- `POSTGRES_MAX_OPEN_CONNECTIONS` - макс. открытых соединений (default: 10)
- `SERVER_HOST` - хост сервера (default: localhost)
- `SERVER_PORT` - порт сервера (default: 8080)

Создайте `.test.env` файл в `app/deploy/local/` на основе примера переменных выше или скопировав `.test.env.example`.

.env - для корректной работы докера(хост базы данных указан по имени контейнера)
.test.env - для запуска интеграционных тестов

Для локального запуска нужно закомментировать запуск сервиса в компоузе и в .env файле поменять значение *POSTGRES_HOST* с db на localhost

## install go migrate

[link](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## install golang-cilint

[link](https://golangci-lint.run/welcome/install/)


## create migration files
```sh
migrate create -dir migrations -ext sql -seq <your migration name>
```

## запуск постгреса локально
```sh
docker compose up -d

docker compose stop
```


