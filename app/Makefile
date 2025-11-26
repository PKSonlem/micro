.PHONY: test
test:
	go test ./... -v

.PHONY: test-integration
test-integration:
	go test ./... -v --tags=integration


# устанавливает зависимости
.PHONY: install
install:
	go install go.uber.org/mock/mockgen@latest
	# Устанавливаем последнюю версию golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(HOME)/go/bin

.PHONY: generate
generate:
	go generate ./...

.PHONY: fmt
fmt:
	go fmt ./...
	golangci-lint run --fix ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

.PHONY: up
up:
	docker compose -f ./deploy/local/docker-compose.yml up -d --build --quiet-pull 

.PHONY: down
down:
	docker compose -f ./deploy/local/docker-compose.yml down 
