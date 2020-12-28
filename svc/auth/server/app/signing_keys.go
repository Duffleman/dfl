package app

import (
	"crypto/x509"
	"encoding/pem"
)

type SigningKeys struct {
	private interface{}
	public  interface{}
}

func (sk SigningKeys) Public() interface{} {
	return sk.public
}

func (sk SigningKeys) Private() interface{} {
	return sk.private
}

type SigningKeysInput struct {
	Private string
	Public  string
}

func ParseKeys(in *SigningKeysInput) (*SigningKeys, error) {
	var err error
	sk := &SigningKeys{}

	privateBlock, _ := pem.Decode([]byte(in.Private))
	x509EncodedPriv := privateBlock.Bytes

	sk.private, err = x509.ParsePKCS8PrivateKey(x509EncodedPriv)
	if err != nil {
		return nil, err
	}

	publicBlock, _ := pem.Decode([]byte(in.Public))
	x509EncodedPub := publicBlock.Bytes

	sk.public, err = x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, err
	}

	return sk, nil
}
