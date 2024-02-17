FROM golang:1.22.0-alpine3.19

ENV POSTGRES_DSN=""
ENV API_PORT=""

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY /cmd ./cmd
COPY /internal ./internal

RUN go build -o ./bin/api ./cmd/

CMD ["./bin/api"]
