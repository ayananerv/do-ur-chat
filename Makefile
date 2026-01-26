.PHONY: run build test migrate proto

run:
	go run ./cmd/server

build:
	go build -o ./bin/server ./cmd/server

test:
	go test ./... -v

migrate:
	go run ./cmd/migrate/main.go

docker-build:
	docker build -t do-ur-chat .

proto:
	protoc --proto_path=api/proto \
		   --go_out=. --go_opt=module=github.com/ayananerv/do-ur-chat \
		   api/proto/*.proto