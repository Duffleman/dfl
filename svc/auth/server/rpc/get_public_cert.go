package rpc

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"

	"dfl/svc/auth/server/app"
	"dfl/svc/auth/server/lib/rpc"
)

func GetPublicCert(a *app.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pk := a.SK.Public().(*ecdsa.PublicKey)

		bytes, err := x509.MarshalPKIXPublicKey(pk)
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}

		w.WriteHeader(200)
		err = pem.Encode(w, &pem.Block{Type: "PUBLIC KEY", Bytes: bytes})
		if err != nil {
			rpc.HandleError(w, r, err)
			return
		}
	}
}
