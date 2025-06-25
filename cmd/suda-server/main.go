package main

import (
	"net/http"
	v1 "suda-backend/internal/api/v1"
)

func main() {
	http.HandleFunc("/api/v1/device/cpu", v1.GetCPUInfo)

	http.ListenAndServe(":8080", nil)
}
