build:
	docker build -t secretaria-bot .

run:
	docker run --name secretaria-bot --env-file secret.env -v secretaria-bot_log:/app/logs -p 8082:8082 -d secretaria-bot

start:
	docker start secretaria-bot

stop:
	docker stop secretaria-bot

rm:
	docker rm secretaria-bot

protoc:
	protoc internal/adapters/driven/financeservice/proto/finance.proto --go_out=internal/adapters/driven/financeservice --go-grpc_out=internal/adapters/driven/financeservice

.PHONY: build run start stop rm protoc