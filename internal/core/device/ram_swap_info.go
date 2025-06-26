package device

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

type ramInfo struct {
	Total       int
	Used        int
	Temperature int
}
type SwapInfo struct {
	Total int
	Used  int
}

func GetRamInfo() (ramInfo, error) {
	var info ramInfo
	cmd := exec.Command("free", "-m")
	output, err := cmd.Output()
	if err != nil {
		return info, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 3 {
		return info, errors.New("unexpected output from free")
	}

	memFields := strings.Fields(lines[1])
	if len(memFields) < 3 {
		return info, errors.New("could not parse RAM info")
	}
	info.Total, err = strconv.Atoi(memFields[1])
	if err != nil {
		return info, err
	}
	info.Used, err = strconv.Atoi(memFields[2])
	if err != nil {
		return info, err
	}
	tempCmd := exec.Command("sensors")
	var tempOut bytes.Buffer
	tempCmd.Stdout = &tempOut
	err = tempCmd.Run()
	if err != nil {
		info.Temperature = -1
	} else {
		info.Temperature = parseRamTemp(tempOut.String())
	}

	return info, nil
}

func GetSwapInfo() (SwapInfo, error) {
	var s SwapInfo

	cmd := exec.Command("free", "-m")
	output, err := cmd.Output()
	if err != nil {
		return s, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 3 {
		return s, errors.New("unexpected output from free")
	}

	swapFields := strings.Fields(lines[2])
	if len(swapFields) < 3 {
		return s, errors.New("could not parse swap info")
	}
	s.Total, err = strconv.Atoi(swapFields[1])
	if err != nil {
		return s, err
	}
	s.Used, err = strconv.Atoi(swapFields[2])
	if err != nil {
		return s, err
	}

	return s, nil
}

func parseRamTemp(s string) int {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "DIMM") && strings.Contains(line, "°C") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.Contains(part, "°C") {
					part = strings.TrimPrefix(part, "+")
					part = strings.TrimSuffix(part, "°C")
					if temp, err := strconv.ParseFloat(part, 64); err == nil {
						return int(temp)
					}
				}
			}
		}
	}
	return -1
}
