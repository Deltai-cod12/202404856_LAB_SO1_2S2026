package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func getDockerContainers() ([]ContainerInfo, error) {
	cmd := exec.Command(
		"docker",
		"ps",
		"--filter",
		"label=so1.profile",
		"--format",
		"{{.ID}}|{{.Names}}",
	)

	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	lines := strings.Split(
		strings.TrimSpace(string(output)),
		"\n",
	)

	var containers []ContainerInfo

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 2)

		if len(parts) != 2 {
			continue
		}

		containerID := parts[0]
		containerName := parts[1]

		inspectCmd := exec.Command(
			"docker",
			"inspect",
			"-f",
			"{{.State.Pid}}|{{index .Config.Labels \"so1.profile\"}}|{{index .Config.Labels \"so1.resource\"}}",
			containerID,
		)

		inspectOutput, err := inspectCmd.Output()

		if err != nil {
			continue
		}

		inspectData := strings.TrimSpace(
			string(inspectOutput),
		)

		inspectParts := strings.SplitN(
			inspectData,
			"|",
			3,
		)

		if len(inspectParts) != 3 {
			continue
		}

		pid, err := strconv.Atoi(inspectParts[0])

		if err != nil {
			continue
		}

		containers = append(
			containers,
			ContainerInfo{
				ID:       containerID,
				Name:     containerName,
				PID:      pid,
				Profile:  inspectParts[1],
				Resource: inspectParts[2],
			},
		)
	}

	return containers, nil
}

func findContainerByPID(
	containers []ContainerInfo,
	pid uint32,
) *ContainerInfo {

	for i := range containers {
		if containers[i].PID == int(pid) {
			return &containers[i]
		}
	}

	return nil
}