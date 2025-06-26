// ======================
// Helper Functions
// ======================
package v1

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, err error, status int) {
	http.Error(w, msg+": "+err.Error(), status)
}
