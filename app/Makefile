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

.PHONY: generate
generate:
	go generate ./...
