savedcmd_modulo.mod := printf '%s\n'   modulo.o | awk '!x[$$0]++ { print("./"$$0) }' > modulo.mod
