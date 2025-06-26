package v1

import (
	"net/http"
	"time"

	"suda-backend/internal/core/device"
)

// =====================
// Response Type Structs
// =====================

type CPUInfoResponse struct {
	ModelName string `json:"model_name"`
	Cores     int    `json:"cores"`
}

type CpuDetailed struct {
	Percentage  []float64 `json:"percentage"`
	Temperature []int     `json:"temperature"`
}

type RamInfo struct {
	Total       int `json:"total"`
	Used        int `json:"used"`
	Temperature int `json:"temperature"`
}

type SwapInfo struct {
	Total int `json:"total"`
	Used  int `json:"used"`
}

// ======================
// Handler Functions
// ======================

func GetCPUInfo(w http.ResponseWriter, r *http.Request) {

	info, err := device.ReadCpuBasicInfo()
	if err != nil {
		writeError(w, "Failed to read CPU info", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, CPUInfoResponse{
		ModelName: info.ModelName,
		Cores:     info.Cores,
	})
}

func GetCpuDetail(w http.ResponseWriter, r *http.Request) {
	info, err := device.GetCpuInfo(1 * time.Second)
	if err != nil {
		writeError(w, "Failed to read detailed CPU info", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, CpuDetailed{
		Percentage:  info.Percentage,
		Temperature: info.Temperature,
	})
}

func GetRam(w http.ResponseWriter, r *http.Request) {
	info, err := device.GetRamInfo()
	if err != nil {
		writeError(w, "Failed to read RAM info", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, RamInfo{
		Total:       info.Total,
		Used:        info.Used,
		Temperature: info.Temperature,
	})
}

func GetSwapInfo(w http.ResponseWriter, r *http.Request) {
	info, err := device.GetSwapInfo()
	if err != nil {
		writeError(w, "Failed to read Swap info", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, SwapInfo{
		Total: info.Total,
		Used:  info.Used,
	})
}

func GetTmuxSessions(w http.ResponseWriter, r *http.Request) {
	info, err := device.GetRamInfo()
	if err != nil {
		writeError(w, "Failed to read Swap info", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, SwapInfo{
		Total: info.Total,
		Used:  info.Used,
	})
}
