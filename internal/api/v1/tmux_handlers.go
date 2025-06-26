package v1

import (
	"net/http"
	"suda-backend/internal/core/tmux"
)

type TmuxInfo struct {
	Used    int    `json:"used"`
	Name    string `json:"name"`
	Windows string `json:"windows"`
	Created string `json:"created"`
	Size    string `json:"size"`
}

func GetTmuxInfo(w http.ResponseWriter, r *http.Request) {

	info, err := tmux.GetTmuxSessions()
	if err != nil {
		writeError(w, "Failed to read TMux info", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, info)
}
