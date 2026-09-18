#!/bin/bash

set -e

MODULE="modulo.ko"
PRIVATE_KEY="keys/MOK.priv"
PUBLIC_KEY="keys/MOK.der"
SIGN_FILE="/usr/src/linux-headers-$(uname -r)/scripts/sign-file"

echo "=== Firmando modulo del kernel ==="

if [ ! -f "$MODULE" ]; then
    echo "Error: No se encontro $MODULE"
    exit 1
fi

if [ ! -f "$PRIVATE_KEY" ] || [ ! -f "$PUBLIC_KEY" ]; then
    echo "Error: No se encontraron las claves MOK"
    exit 1
fi

sudo "$SIGN_FILE" sha256 "$PRIVATE_KEY" "$PUBLIC_KEY" "$MODULE"

echo "Modulo firmado correctamente."
echo ""
modinfo "./$MODULE" | grep -E "signer|sig_key|sig_hashalgo"