package rpc

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"

	"dfl/svc/auth/server/app"
)

func GetPublicCert(a *app.App, w http.ResponseWriter, r *http.Request) error {
	pk := a.SK.Public().(*ecdsa.PublicKey)

	bytes, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return err
	}

	w.WriteHeader(200)
	err = pem.Encode(w, &pem.Block{Type: "PUBLIC KEY", Bytes: bytes})
	if err != nil {
		return err
	}

	return nil
}
