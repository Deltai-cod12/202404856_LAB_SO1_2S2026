package main

import "sort"

type ContainerDecision struct {
	Keep   []ContainerInfo
	Delete []ContainerInfo
}

func metricValue(
	container ContainerInfo,
	metric string,
) uint64 {

	if container.Process == nil {
		return 0
	}

	switch metric {
	case "rss":
		return container.Process.RSSKB

	case "vsz":
		return container.Process.VSZKB

	case "memory":
		return container.Process.MemPercentX100

	case "cpu":
		return container.Process.CPUPercentX100
	}

	return 0
}

func calculateScores(
	containers []ContainerInfo,
) map[string]int {

	scores := make(map[string]int)

	metrics := []string{
		"rss",
		"vsz",
		"memory",
		"cpu",
	}

	for _, metric := range metrics {

		sorted := append(
			[]ContainerInfo(nil),
			containers...,
		)

		sort.Slice(
			sorted,
			func(i, j int) bool {
				return metricValue(
					sorted[i],
					metric,
				) > metricValue(
					sorted[j],
					metric,
				)
			},
		)

		total := len(sorted)

		for position, container := range sorted {
			points := total - position

			scores[container.ID] += points
		}
	}

	return scores
}

func selectContainers(
	containers []ContainerInfo,
	keepCount int,
) ContainerDecision {

	var decision ContainerDecision

	if len(containers) <= keepCount {
		decision.Keep = append(
			decision.Keep,
			containers...,
		)

		return decision
	}

	scores := calculateScores(containers)

	sorted := append(
		[]ContainerInfo(nil),
		containers...,
	)

	sort.Slice(
		sorted,
		func(i, j int) bool {

			scoreI := scores[sorted[i].ID]
			scoreJ := scores[sorted[j].ID]

			if scoreI == scoreJ {
				return sorted[i].ID < sorted[j].ID
			}

			return scoreI > scoreJ
		},
	)

	decision.Keep = append(
		decision.Keep,
		sorted[:keepCount]...,
	)

	decision.Delete = append(
		decision.Delete,
		sorted[keepCount:]...,
	)

	return decision
}

func selectHighContainers(
	containers []ContainerInfo,
) ContainerDecision {

	var decision ContainerDecision

	var ramContainers []ContainerInfo
	var cpuContainers []ContainerInfo

	for _, container := range containers {

		switch container.Resource {

		case "ram":
			ramContainers = append(
				ramContainers,
				container,
			)

		case "cpu":
			cpuContainers = append(
				cpuContainers,
				container,
			)
		}
	}

	// Ordenar HIGH RAM.
	sort.Slice(
		ramContainers,
		func(i, j int) bool {

			a := ramContainers[i].Process
			b := ramContainers[j].Process

			if a == nil {
				return false
			}

			if b == nil {
				return true
			}

			if a.MemPercentX100 != b.MemPercentX100 {
				return a.MemPercentX100 >
					b.MemPercentX100
			}

			if a.RSSKB != b.RSSKB {
				return a.RSSKB > b.RSSKB
			}

			return a.VSZKB > b.VSZKB
		},
	)

	// Ordenar HIGH CPU.
	sort.Slice(
		cpuContainers,
		func(i, j int) bool {

			a := cpuContainers[i].Process
			b := cpuContainers[j].Process

			if a == nil {
				return false
			}

			if b == nil {
				return true
			}

			return a.CPUPercentX100 >
				b.CPUPercentX100
		},
	)

	// Caso 1:
	// Existen ambos perfiles.
	if len(ramContainers) > 0 &&
		len(cpuContainers) > 0 {

		decision.Keep = append(
			decision.Keep,
			ramContainers[0],
			cpuContainers[0],
		)

		if len(ramContainers) > 1 {
			decision.Delete = append(
				decision.Delete,
				ramContainers[1:]...,
			)
		}

		if len(cpuContainers) > 1 {
			decision.Delete = append(
				decision.Delete,
				cpuContainers[1:]...,
			)
		}

		return decision
	}

	// Caso 2:
	// Solo existen HIGH RAM.
	if len(ramContainers) > 0 {

		keepCount := 2

		if len(ramContainers) < keepCount {
			keepCount = len(ramContainers)
		}

		decision.Keep = append(
			decision.Keep,
			ramContainers[:keepCount]...,
		)

		if len(ramContainers) > keepCount {
			decision.Delete = append(
				decision.Delete,
				ramContainers[keepCount:]...,
			)
		}

		return decision
	}

	// Caso 3:
	// Solo existen HIGH CPU.
	if len(cpuContainers) > 0 {

		keepCount := 2

		if len(cpuContainers) < keepCount {
			keepCount = len(cpuContainers)
		}

		decision.Keep = append(
			decision.Keep,
			cpuContainers[:keepCount]...,
		)

		if len(cpuContainers) > keepCount {
			decision.Delete = append(
				decision.Delete,
				cpuContainers[keepCount:]...,
			)
		}
	}

	return decision
}
