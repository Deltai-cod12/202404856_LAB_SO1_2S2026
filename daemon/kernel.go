package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const procFile = "/proc/continfo_pr2_so1_202404856"

func readKernelInfo() (*KernelInfo, error) {
	data, err := os.ReadFile(procFile)

	if err != nil {
		return nil, fmt.Errorf(
			"error al leer %s: %w",
			procFile,
			err,
		)
	}

	var kernelInfo KernelInfo

	err = json.Unmarshal(data, &kernelInfo)

	if err != nil {
		return nil, fmt.Errorf(
			"error al deserializar JSON: %w",
			err,
		)
	}

	return &kernelInfo, nil
}

func findProcessByPID(
	processes []ProcessInfo,
	pid int,
) *ProcessInfo {

	for i := range processes {
		if processes[i].PID == pid {
			return &processes[i]
		}
	}

	return nil
}
