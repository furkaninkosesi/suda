package device

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CPUStat struct {
	User, Nice, System, Idle, Iowait, Irq, SoftIrq, Steal, Guest, GuestNice uint64
}

type CpuInfo struct {
	Percentage  []float64
	Temperature []int
}

func parseCPUStatLine(line string) (CPUStat, error) {
	fields := strings.Fields(line)
	if len(fields) < 11 {
		return CPUStat{}, fmt.Errorf("geçersiz cpu istatistik satırı: yetersiz alan sayısı (%d)", len(fields))
	}

	var stats CPUStat
	var err error

	stats.User, err = strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("user field could not be read: %w", err)
	}
	stats.Nice, err = strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("nice field could not be read: %w", err)
	}
	stats.System, err = strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("system field could not be read: %w", err)
	}
	stats.Idle, err = strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("idle field could not be read: %w", err)
	}
	stats.Iowait, err = strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("iowait field could not be read: %w", err)
	}
	stats.Irq, err = strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("irq field could not be read: %w", err)
	}
	stats.SoftIrq, err = strconv.ParseUint(fields[7], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("softirq field could not be read: %w", err)
	}
	stats.Steal, err = strconv.ParseUint(fields[8], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("steal field could not be read: %w", err)
	}
	stats.Guest, err = strconv.ParseUint(fields[9], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("guest field could not be read: %w", err)
	}
	stats.GuestNice, err = strconv.ParseUint(fields[10], 10, 64)
	if err != nil {
		return CPUStat{}, fmt.Errorf("guestnice field could not be read: %w", err)
	}

	return stats, nil
}

func getCPUStats() (map[string]CPUStat, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, fmt.Errorf("could not open /proc/stat file: %w", err)
	}
	defer file.Close()

	stats := make(map[string]CPUStat)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			cpuID := fields[0]

			stat, err := parseCPUStatLine(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s Could not parse line: %v\n", cpuID, err)
				continue
			}
			stats[cpuID] = stat
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading file: %w", err)
	}

	return stats, nil
}

func calculateCPUUsage(stat1, stat2 CPUStat) float64 {
	idle1 := stat1.Idle + stat1.Iowait
	total1 := stat1.User + stat1.Nice + stat1.System + idle1 + stat1.Irq + stat1.SoftIrq + stat1.Steal

	idle2 := stat2.Idle + stat2.Iowait
	total2 := stat2.User + stat2.Nice + stat2.System + idle2 + stat2.Irq + stat2.SoftIrq + stat2.Steal

	totalDelta := float64(total2 - total1)
	idleDelta := float64(idle2 - idle1)

	if totalDelta == 0 {
		return 0.0
	}

	cpuUsage := (totalDelta - idleDelta) / totalDelta * 100.0
	return cpuUsage
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

func getCPUTemperatures() ([]int, error) {
	var temps []int
	matches, err := filepath.Glob("/sys/class/thermal/thermal_zone*")
	if err != nil {
		return nil, fmt.Errorf("termal bölge dosyaları bulunamadı: %w", err)
	}

	for _, path := range matches {
		typeData, err := ioutil.ReadFile(filepath.Join(path, "type"))
		if err != nil {
			continue
		}
		typeStr := strings.TrimSpace(string(typeData))

		if !strings.Contains(typeStr, "x86_pkg_temp") && !strings.Contains(typeStr, "cpu") {
			continue
		}

		tempData, err := ioutil.ReadFile(filepath.Join(path, "temp"))
		if err != nil {
			continue
		}

		tempMilli, err := strconv.Atoi(strings.TrimSpace(string(tempData)))
		if err != nil {
			continue
		}

		temps = append(temps, tempMilli/1000)
	}

	if len(temps) == 0 {

	}

	return temps, nil
}

func GetCpuInfo(duration time.Duration) (*CpuInfo, error) {
	usageMap, err := GetCPUUsagePerCore(duration)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage data: %w", err)
	}

	var coreKeys []string
	for key := range usageMap {
		if strings.HasPrefix(key, "cpu") && key != "cpu" {
			coreKeys = append(coreKeys, key)
		}
	}
	sort.Strings(coreKeys)

	var usageSlice []float64
	for _, key := range coreKeys {
		usageSlice = append(usageSlice, usageMap[key])
	}

	temps, err := getCPUTemperatures()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to get CPU temperatures: %v\n", err)
		temps = []int{}
	}

	return &CpuInfo{
		Percentage:  usageSlice,
		Temperature: temps,
	}, nil
}
