# Etapa de compilacion
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiar dependencias y código fuente
COPY go.mod ./
COPY main.go ./

# Compilar la aplicacion estaticamente
RUN go build -o api2 main.go

# Etapa final ligera
FROM alpine:latest

WORKDIR /root/

# Copiar el binario compilado desde la etapa anterior
COPY --from=builder /app/api2 .

EXPOSE 8082

CMD ["./api2"]