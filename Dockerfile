FROM golang:1.19-alpine

ENV CGO_ENABLED=0
EXPOSE 8082

WORKDIR /app
COPY . .

CMD ["go", "run", "cmd/main.go"]