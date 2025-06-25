package v1

import (
	"encoding/json"
	"net/http"

	"suda-backend/internal/core/device"
)

type CPUInfoResponse struct {
	ModelName string `json:"model_name"`
	Cores     int    `json:"cores"`
}

func GetCPUInfo(w http.ResponseWriter, r *http.Request) {
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
