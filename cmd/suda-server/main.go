package main

import (
	"net/http"
	v1 "suda-backend/internal/api/v1"
)

func main() {
	http.HandleFunc("/api/v1/device/cpu_info", v1.GetCPUInfo)
	http.HandleFunc("/api/v1/device/cpu_detail", v1.GetCpuDetail)
	http.HandleFunc("/api/v1/device/ram_info", v1.GetRam)
	http.HandleFunc("/api/v1/device/swap_info", v1.GetSwapInfo)

	http.ListenAndServe(":8080", nil)
}
