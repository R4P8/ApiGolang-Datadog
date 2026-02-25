FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o go-apitask ./cmd

FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /app/go-apitask .

EXPOSE 8000

CMD ["./go-apitask"]