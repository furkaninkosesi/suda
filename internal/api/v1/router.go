package v1

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/device/cpu_info", GetCPUInfo)
	mux.HandleFunc("/api/v1/device/cpu_detail", GetCpuDetail)
	mux.HandleFunc("/api/v1/device/ram_info", GetRam)
	mux.HandleFunc("/api/v1/device/swap_info", GetSwapInfo)
	mux.HandleFunc("/api/v1/device/tmux_info", GetTmuxInfo)
}
