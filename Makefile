build:
	docker build -t go-secretaria-bot_app .

run:
	docker run --name go-secretaria-bot -p 8082:8082 -d go-secretaria-bot_app

.PHONY: build run