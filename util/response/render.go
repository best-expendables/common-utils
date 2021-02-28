package response

import (
	"encoding/json"
	"net/http"
)

func RenderJson(w http.ResponseWriter, res ApiResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	for key, value := range res.Headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(res.Code)
	if res.Code == http.StatusNoContent {
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
