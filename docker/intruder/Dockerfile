FROM alpine:latest

WORKDIR /app

COPY intruder.sh .

RUN chmod +x intruder.sh

# Ejecutamos como usuario sin privilegios.
RUN adduser -D intruder

USER intruder

CMD ["./intruder.sh"]