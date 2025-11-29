build:
	docker build -t secretaria-bot .

run:
	docker run --name secretaria-bot --env-file secret.env -v secretaria-bot_log:/app/logs -p 80:80 -d secretaria-bot

start:
	docker start secretaria-bot

stop:
	docker stop secretaria-bot && docker rm secretaria-bot

rm:
	docker rmi secretaria-bot:latest

gripmock-start:
	docker run -d --name gripmock -p 4770:4770 -p 4771:4771 -v ./proto/finance/stubs:/stubs:ro -v ./proto/finance:/proto:ro bavix/gripmock --stub=/stubs /proto/finance.proto

gripmock-stop:
	docker stop gripmock && docker rm gripmock

protoc:
	protoc proto/finance.proto --go_out=internal/adapters/driven/financeservice --go-grpc_out=internal/adapters/driven/financeservice

.PHONY: build run start stop rm protoc