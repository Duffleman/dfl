package app

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path"

	"dfl/lib/cher"
	"dfl/tools/certgen"

	"github.com/fatih/color"
	"github.com/tjarratt/babble"
)

func (a *App) VerifyKeyPair(name string) error {
	if err := a.checkForInit(); err != nil {
		return err
	}

	privFilePath := path.Join(
		a.RootDirectory,
		certgen.CertFolderMap[certgen.KeyPair],
		fmt.Sprintf("%s.private.pem", name),
	)

	pubFilePath := path.Join(
		a.RootDirectory,
		certgen.CertFolderMap[certgen.KeyPair],
		fmt.Sprintf("%s.public.pem", name),
	)

	privatePemBytes, err := ioutil.ReadFile(privFilePath)
	if err != nil {
		return err
	}

	publicPemBytes, err := ioutil.ReadFile(pubFilePath)
	if err != nil {
		return err
	}

	privateBytes, _ := pem.Decode(privatePemBytes)
	privBlock := privateBytes.Bytes

	publicBytes, _ := pem.Decode(publicPemBytes)
	pubBlock := publicBytes.Bytes

	private, err := x509.ParseECPrivateKey(privBlock)
	if err != nil {
		return err
	}

	publicInt, err := x509.ParsePKIXPublicKey(pubBlock)
	if err != nil {
		return err
	}

	var public *ecdsa.PublicKey

	if v, ok := publicInt.(*ecdsa.PublicKey); ok {
		public = v
	} else {
		return cher.New("public_key_malformed", nil)
	}

	message := makeMessage()
	hash := sha256.Sum256([]byte(message))

	signature, err := ecdsa.SignASN1(rand.Reader, private, hash[:])
	if err != nil {
		return err
	}

	c := color.New()

	valid := ecdsa.VerifyASN1(public, hash[:], signature)

	c.Printf("Message: %s\n", message)
	c.Print("Signature ")

	switch valid {
	case true:
		c.Add(color.BgGreen)
		c.Println("VERIFIED")
	default:
	case true:
		c.Add(color.BgRed)
		c.Println("MISMATCH")
	}

	return nil
}

func makeMessage() string {
	babbler := babble.NewBabbler()
	babbler.Count = 4

	return babbler.Babble()
}
