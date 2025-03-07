### general tips

## load envs
```sh
set -a
. ./.env
```

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


