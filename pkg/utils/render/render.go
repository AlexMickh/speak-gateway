package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, obj any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		return err
	}

	return nil
}
