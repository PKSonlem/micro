.PHONY: generate
generate:
	protoc --go_out=. --go-grpc_out=. ./proto/notification.proto
	go generate ./...


.PHONY: install
install:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

