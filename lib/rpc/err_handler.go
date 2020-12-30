package rpc

import (
	"encoding/json"
	"net/http"

	"dfl/lib/cher"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	if e, ok := err.(cher.E); ok {
		switch e.Code {
		case cher.NotFound:
			w.WriteHeader(404)
		case cher.AccessDenied:
			w.WriteHeader(403)
		case cher.Unauthorized:
			w.WriteHeader(401)
		default:
			w.WriteHeader(500)
		}

		json.NewEncoder(w).Encode(e)
		return
	}

	w.WriteHeader(500)
	json.NewEncoder(w).Encode(cher.New("unknown", cher.M{"error": err.Error()}))
}
