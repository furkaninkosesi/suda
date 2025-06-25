package device

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type CPUStat struct {
	User, Nice, System, Idle, Iowait, Irq, SoftIrq, Steal, Guest, GuestNice uint64
}

func parseCPUStatLine(line string) (CPUStat, error) {
	fields := strings.Fields(line)
	if len(fields) < 11 {
		return CPUStat{}, fmt.Errorf("invalid cpu stat line")
	}

	var stats CPUStat
	var err error
	stats.User, err = strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return stats, err
	}
	stats.Nice, _ = strconv.ParseUint(fields[2], 10, 64)
	stats.System, _ = strconv.ParseUint(fields[3], 10, 64)
	stats.Idle, _ = strconv.ParseUint(fields[4], 10, 64)
	stats.Iowait, _ = strconv.ParseUint(fields[5], 10, 64)
	stats.Irq, _ = strconv.ParseUint(fields[6], 10, 64)
	stats.SoftIrq, _ = strconv.ParseUint(fields[7], 10, 64)
	stats.Steal, _ = strconv.ParseUint(fields[8], 10, 64)
	stats.Guest, _ = strconv.ParseUint(fields[9], 10, 64)
	stats.GuestNice, _ = strconv.ParseUint(fields[10], 10, 64)

	return stats, nil
}

func getCPUStats() (map[string]CPUStat, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := make(map[string]CPUStat)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") {
			stat, err := parseCPUStatLine(line)
			if err != nil {
				return nil, err
			}
			fields := strings.Fields(line)
			cpuID := fields[0] // cpu, cpu0, cpu1, ...
			stats[cpuID] = stat
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func calculateCPUUsage(stat1, stat2 CPUStat) float64 {
	idle1 := stat1.Idle + stat1.Iowait
	idle2 := stat2.Idle + stat2.Iowait

	total1 := stat1.User + stat1.Nice + stat1.System + stat1.Irq + stat1.SoftIrq + stat1.Steal + idle1
	total2 := stat2.User + stat2.Nice + stat2.System + stat2.Irq + stat2.SoftIrq + stat2.Steal + idle2

	totald := float64(total2 - total1)
	idled := float64(idle2 - idle1)

	if totald == 0 {
		return 0
	}

	return (totald - idled) / totald * 100
}

func GetCPUUsagePerCore(duration time.Duration) (map[string]float64, error) {
	stats1, err := getCPUStats()
	if err != nil {
		return nil, err
	}

	time.Sleep(duration)

	stats2, err := getCPUStats()
	if err != nil {
		return nil, err
	}

	usage := make(map[string]float64)
	for cpuID, stat1 := range stats1 {
		stat2, ok := stats2[cpuID]
		if !ok {
			continue
		}
		usage[cpuID] = calculateCPUUsage(stat1, stat2)
	}

	return usage, nil
}

type CpuInfo struct {
	Percentage  []float64
	Temperature []int
}

func getCPUTemperatures() ([]int, error) {
	temps := []int{}
	matches, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil {
		return nil, err
	}

	for _, path := range matches {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}
		strData := strings.TrimSpace(string(data))
		tempMilli, err := strconv.Atoi(strData)
		if err != nil {
			continue
		}
		temp := tempMilli / 1000
		temps = append(temps, temp)
	}

	return temps, nil
}

func GetCpuInfo(duration time.Duration) (*CpuInfo, error) {
	usageMap, err := GetCPUUsagePerCore(duration)
	if err != nil {
		return nil, err
	}

	var usageSlice []float64
	for i := 0; ; i++ {
		cpuID := fmt.Sprintf("cpu%d", i)
		val, ok := usageMap[cpuID]
		if !ok {
			break
		}
		usageSlice = append(usageSlice, val)
	}

	temps, _ := getCPUTemperatures()

	return &CpuInfo{
		Percentage:  usageSlice,
		Temperature: temps,
	}, nil
}
