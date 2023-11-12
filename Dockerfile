FROM golang:1.21.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

FROM builder AS build

COPY . .

RUN go build -o ./../../api cmd/api/main.go

CMD ["/api"]