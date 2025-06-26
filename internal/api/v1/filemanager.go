package v1

import (
	"net/http"
	"suda-backend/internal/core/filemanager"
)

func GetDirectoryContents(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		writeError(w, "Path parametresi eksik", nil, http.StatusBadRequest)
		return
	}

	contents, err := filemanager.ListDirectory(path)
	if err != nil {
		writeError(w, "Klasör okunamadı", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, contents)
}
