#!/bin/sh

echo "========================================"
echo "SO1 - CONTENEDOR INTRUSO"
echo "Simulacion de lectura de archivos sensibles"
echo "========================================"

while true
do
    echo ""
    echo "[INTRUDER] Intentando leer /etc/passwd..."
    cat /etc/passwd > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "[INTRUDER] Lectura de /etc/passwd realizada."
    else
        echo "[INTRUDER] No fue posible leer /etc/passwd."
    fi

    echo "[INTRUDER] Intentando leer /etc/shadow..."
    cat /etc/shadow > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "[INTRUDER] Lectura de /etc/shadow realizada."
    else
        echo "[INTRUDER] Acceso a /etc/shadow denegado."
    fi

    echo "[INTRUDER] Intentando leer /etc/hostname..."
    cat /etc/hostname > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "[INTRUDER] Lectura de /etc/hostname realizada."
    fi

    sleep 10
done