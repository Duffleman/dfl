package rpc

import (
	"encoding/json"
	"net/http"
)

func WriteOut(w http.ResponseWriter, r *http.Request, res interface{}) {
	if res == nil {
		w.WriteHeader(204)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		HandleError(w, r, err)
	}
}
