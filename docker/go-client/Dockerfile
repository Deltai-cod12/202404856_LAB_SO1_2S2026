FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY main.go .

RUN go build -o go-client .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go-client .

CMD ["./go-client"]