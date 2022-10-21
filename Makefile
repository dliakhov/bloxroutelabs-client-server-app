.PHONY: gen-protobuf run-demo-docker-compose stop-demo-docker-compose run-server run-client run-rabbit-mq run-tests

gen-protobuf:
	protoc --proto_path=models --go_out=models --go_opt=paths=source_relative models/command.proto

run-demo-docker-compose:
	docker compose up

stop-demo-docker-compose:
	docker compose down --remove-orphans

run-server:
	go run main.go server

run-client:
	go run main.go client

run-rabbit-mq:
	docker compose run rabbit

run-tests:
	go test ./...
