package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Daemon Proyecto 2 SO1 ===")

	valkey := newValkeyClient()

	defer valkey.close()

	if err := valkey.ping(); err != nil {
		fmt.Println("Error conectando con Valkey:", err)
		return
	}

	if err := configureEBPFPermissions(); err != nil {
		fmt.Println(err)
		return
	}

	ebpfMonitor, err := newEBPFMonitor()

	if err != nil {
		fmt.Println("Error iniciando eBPF:", err)
		return
	}

	defer ebpfMonitor.close()

	go ebpfMonitor.readEvents()


	kernelInfo, err := readKernelInfo()

	if err != nil {
		fmt.Println("Error leyendo informacion del kernel:", err)
		return
	}

	if err := valkey.saveRAM(kernelInfo.RAM); err != nil {
		fmt.Println("Error guardando RAM en Valkey:", err)
	} else {
		fmt.Println("Metricas de RAM almacenadas en Valkey.")
	}

	fmt.Println()
	fmt.Println("=== INFORMACION DE RAM ===")

	fmt.Printf(
		"RAM Total: %d KB\n",
		kernelInfo.RAM.TotalKB,
	)

	fmt.Printf(
		"RAM Libre: %d KB\n",
		kernelInfo.RAM.FreeKB,
	)

	fmt.Printf(
		"RAM Utilizada: %d KB\n",
		kernelInfo.RAM.UsedKB,
	)

	fmt.Printf(
		"Procesos recibidos: %d\n",
		len(kernelInfo.Processes),
	)

	containers, err := getDockerContainers()

	if err != nil {
		fmt.Printf(
			"Error al obtener contenedores: %v\n",
			err,
		)

		return
	}

	fmt.Println()
	fmt.Println("=== CONTENEDORES DOCKER ===")

	for i := range containers {

		containers[i].Process = findProcessByPID(
			kernelInfo.Processes,
			containers[i].PID,
		)

		container := containers[i]

		fmt.Println()
		fmt.Printf(
			"Contenedor: %s\n",
			container.Name,
		)

		fmt.Printf(
			"ID: %s\n",
			container.ID,
		)

		fmt.Printf(
			"PID: %d\n",
			container.PID,
		)

		if container.Process == nil {
			fmt.Println(
				"Metricas no encontradas.",
			)

			continue
		}

		process := container.Process

		fmt.Printf(
			"CMD: %s\n",
			process.Cmd,
		)

		fmt.Printf(
			"VSZ: %d KB\n",
			process.VSZKB,
		)

		fmt.Printf(
			"RSS: %d KB\n",
			process.RSSKB,
		)

		fmt.Printf(
			"MEM: %.2f%%\n",
			float64(process.MemPercentX100)/100.0,
		)

		fmt.Printf(
			"CPU: %.2f%%\n",
			float64(process.CPUPercentX100)/100.0,
		)
	}

	if err := valkey.saveContainers(containers); err != nil {
		fmt.Println(
			"Error guardando contenedores en Valkey:",
			err,
		)
	} else {
		fmt.Printf(
			"Metricas de %d contenedores almacenadas en Valkey.\n",
			len(containers),
		)
	}

	fmt.Println()
	fmt.Println("=== CLASIFICACION DE CONTENEDORES ===")

	var highContainers []ContainerInfo
	var lowContainers []ContainerInfo

	for _, container := range containers {

		switch container.Profile {

		case "high":
			highContainers = append(
				highContainers,
				container,
			)

		case "low":
			lowContainers = append(
				lowContainers,
				container,
			)
		}
	}

	fmt.Printf(
		"Alto consumo: %d\n",
		len(highContainers),
	)

	fmt.Printf(
		"Bajo consumo: %d\n",
		len(lowContainers),
	)

	fmt.Println()
	fmt.Println("--- ALTO CONSUMO ---")

	for _, container := range highContainers {
		fmt.Printf(
			"%s | Recurso: %s | PID: %d\n",
			container.Name,
			container.Resource,
			container.PID,
		)
	}

	fmt.Println()
	fmt.Println("--- BAJO CONSUMO ---")

	for _, container := range lowContainers {
		fmt.Printf(
			"%s | Recurso: %s | PID: %d\n",
			container.Name,
			container.Resource,
			container.PID,
		)
	}

	fmt.Println()
	fmt.Println("=== POLITICA DE ADMINISTRACION ===")

	highDecision := selectHighContainers(
		highContainers,
	)

	lowDecision := selectContainers(
		lowContainers,
		3,
	)

	fmt.Println()
	fmt.Println("--- CONSERVAR HIGH ---")

	for _, container := range highDecision.Keep {

		process := container.Process

		if process == nil {
			continue
		}

		fmt.Printf(
			"[KEEP] %s | Recurso: %s | RSS: %d KB | VSZ: %d KB | MEM: %.2f%% | CPU: %.2f%%\n",
			container.Name,
			container.Resource,
			process.RSSKB,
			process.VSZKB,
			float64(process.MemPercentX100)/100.0,
			float64(process.CPUPercentX100)/100.0,
		)
	}

	fmt.Println()
	fmt.Println("--- ELIMINAR HIGH ---")

	for _, container := range highDecision.Delete {

		process := container.Process

		if process == nil {
			continue
		}

		fmt.Printf(
			"[DELETE] %s | Recurso: %s | RSS: %d KB | VSZ: %d KB | MEM: %.2f%% | CPU: %.2f%%\n",
			container.Name,
			container.Resource,
			process.RSSKB,
			process.VSZKB,
			float64(process.MemPercentX100)/100.0,
			float64(process.CPUPercentX100)/100.0,
		)
	}

	fmt.Println()
	fmt.Println("--- CONSERVAR LOW ---")

	for _, container := range lowDecision.Keep {

		score := calculateScores(
			lowContainers,
		)[container.ID]

		fmt.Printf(
			"[KEEP] %s | Score: %d\n",
			container.Name,
			score,
		)
	}

	fmt.Println()
	fmt.Println("=== MONITOR eBPF ===")
	fmt.Println("Esperando eventos durante 30 segundos...")

	go func() {

		for event := range ebpfMonitor.events {
		
			container := findContainerByPID(
				containers,
				event.TargetPID,
			)
		
			if container == nil {
				continue
			}
		
			fmt.Printf(
				"[eBPF] Evento de contenedor | Nombre: %s | ID: %s | PID: %d | Signal: %d\n",
				container.Name,
				container.ID,
				container.PID,
				event.Signal,
			)
		}
	}()
	
	time.Sleep(30 * time.Second)

	fmt.Println()
	fmt.Println("--- ELIMINAR LOW ---")

	for _, container := range lowDecision.Delete {

		score := calculateScores(
			lowContainers,
		)[container.ID]

		fmt.Printf(
			"[DELETE] %s | Score: %d\n",
			container.Name,
			score,
		)
	}

}
