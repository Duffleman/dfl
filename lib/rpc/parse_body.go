package rpc

import (
	"encoding/json"
	"net/http"
)

func ParseBody(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
