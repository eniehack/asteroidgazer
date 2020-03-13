package rsax

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func ConvertPublicKeyToString(publickey *rsa.PublicKey) string {
	pempublickey := new(pem.Block)
	pempublickey.Type = "RSA PUBLIC KEY"
	pempublickey.Bytes = x509.MarshalPKCS1PublicKey(publickey)
	encodedkey := pem.EncodeToMemory(pempublickey)
	return string(encodedkey)
}
