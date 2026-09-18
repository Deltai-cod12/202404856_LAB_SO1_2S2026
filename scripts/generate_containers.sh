#!/bin/bash

set -e

TOTAL_CONTAINERS=5

echo "========================================"
echo "Generador de contenedores - Proyecto 2"
echo "========================================"
echo "Creando $TOTAL_CONTAINERS contenedores..."
echo ""

for i in $(seq 1 $TOTAL_CONTAINERS)
do
    PROFILE=$((RANDOM % 4))

    TIMESTAMP=$(date +%s%N)
    NAME="so1-${TIMESTAMP}-${i}"

    case $PROFILE in

        0)
            echo "[$i/$TOTAL_CONTAINERS] LOW -> $NAME"

            docker run -d \
                --name "$NAME" \
                --label so1.profile=low \
                --label so1.resource=low \
                alpine \
                sleep 3600 \
                > /dev/null
            ;;

        1)
            echo "[$i/$TOTAL_CONTAINERS] HIGH CPU -> $NAME"

            docker run -d \
                --name "$NAME" \
                --label so1.profile=high \
                --label so1.resource=cpu \
                alpine \
                sh -c 'while true; do :; done' \
                > /dev/null
            ;;

        2)
            echo "[$i/$TOTAL_CONTAINERS] HIGH RAM -> $NAME"

            docker run -d \
                --name "$NAME" \
                --label so1.profile=high \
                --label so1.resource=ram \
                so1-go-client \
                > /dev/null
            ;;

        3)
            echo "[$i/$TOTAL_CONTAINERS] INTRUDER -> $NAME"

            docker run -d \
                --name "$NAME" \
                --label so1.profile=low \
                --label so1.resource=intruder \
                --label so1.intruder=true \
                so1-intruder \
                > /dev/null
            ;;

    esac
done

echo ""
echo "Generacion finalizada."