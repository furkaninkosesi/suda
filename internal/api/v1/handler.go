package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"suda-backend/internal/core/device"
)

func GetCPUInfo(w http.ResponseWriter, r *http.Request) {
	type CPUInfoResponse struct {
		ModelName string `json:"model_name"`
		Cores     int    `json:"cores"`
	}
	info, err := device.ReadCpuBasicInfo()
	if err != nil {
		http.Error(w, "Failed to read CPU info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CPUInfoResponse{
		ModelName: info.ModelName,
		Cores:     info.Cores,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetCpuDetail(w http.ResponseWriter, r *http.Request) {
	type CpuDetailed struct {
		Percentage  []float64 `json:"percentage"`
		Temperature []int     `json:"temperature"`
	}

	info, err := device.GetCpuInfo(1 * time.Second)
	if err != nil {
		http.Error(w, "Failed to read CPU info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CpuDetailed{
		Percentage:  info.Percentage,
		Temperature: info.Temperature,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
