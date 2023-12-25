FROM golang:alpine3.19

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8082

CMD ["go", "run", "cmd/main.go"]