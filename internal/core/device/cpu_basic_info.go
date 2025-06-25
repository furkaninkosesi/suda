package device

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CpuBasicInfo struct {
	ModelName string
	Cores     int
}

func ReadCpuBasicInfo() (CpuBasicInfo, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return CpuBasicInfo{}, err
	}
	var model string
	var cores int

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				model = strings.TrimSpace(parts[1])
			}
		} else if strings.HasPrefix(line, "cpu cores") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				cores, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}
		}
	}

	if model == "" {
		return CpuBasicInfo{}, fmt.Errorf("could not find CPU model name")
	}
	if cores == 0 {
		return CpuBasicInfo{}, fmt.Errorf("could not find CPU cores count")
	}

	return CpuBasicInfo{
		ModelName: model,
		Cores:     cores,
	}, nil

}
