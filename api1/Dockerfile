# Etapa de compilación
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiar dependencias y código fuente
COPY go.mod ./
COPY main.go ./

# Compilar la aplicación estáticamente
RUN CGO_ENABLED=0 GOOS=linux go build -o api1 main.go

# Etapa final ligera
FROM alpine:latest

WORKDIR /root/

# Copiar el binario compilado desde la etapa anterior
COPY --from=builder /app/api1 .

EXPOSE 8081

CMD ["./api1"]