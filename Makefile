build:
	docker build -t go-secretaria-bot_app .

run:
	docker run --name go-secretaria-bot --env-file secret.env -v go-secretaria-bot_log:/app/logs -p 8082:8082 -d go-secretaria-bot_app

start:
	docker start go-secretaria-bot

stop:
	docker stop go-secretaria-bot

rm:
	docker rm go-secretaria-bot

protoc:
	protoc internal/adapters/driven/financeservice/proto/finance.proto --go_out=internal/adapters/driven/financeservice --go-grpc_out=internal/adapters/driven/financeservice

.PHONY: build run start stop rm protoc