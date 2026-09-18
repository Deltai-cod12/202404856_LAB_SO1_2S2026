package main

import (
	"fmt"
	"time"
)

const memoryMB = 200

func main() {
	fmt.Printf(
		"SO1 Go Client - reservando aproximadamente %d MB de RAM\n",
		memoryMB,
	)

	memory := make([]byte, memoryMB*1024*1024)

	// Fuerza la asignacion fisica inicial de las paginas.
	for i := 0; i < len(memory); i += 4096 {
		memory[i] = 1
	}

	fmt.Println("Memoria reservada correctamente.")

	for {
		// Volvemos a tocar periodicamente cada pagina para
		// mantener el conjunto de memoria activamente utilizado.
		for i := 0; i < len(memory); i += 4096 {
			memory[i]++
		}

		time.Sleep(5 * time.Second)
	}
}