# Etapa de compilacion
FROM docker.io/golang:1.22-alpine AS builder

WORKDIR /app

# Copiar dependencias y código fuente
COPY go.mod ./
COPY main.go ./

# Compilar la aplicacion estaticamente
RUN go build -o api3 main.go

# Etapa final ligera
FROM docker.io/alpine:latest

WORKDIR /root/

COPY --from=builder /app/api3 .

EXPOSE 8083

CMD ["./api3"]