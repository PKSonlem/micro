### general tips

## load envs
```sh
set -a
. ./.env
```

## install go migrate

[link](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)


## create migration files
```sh
migrate create -dir migrations -ext sql -seq <your migration name>
```

## запуск постгреса локально
```sh
docker compose up -d

docker compose stop
```
