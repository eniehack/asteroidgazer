package rsax

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	InvaildPEM       = errors.New("Invaild pem file format")
	NotRSAPrivateKey = errors.New("pem file is not RSA PRIVATE KEY.")
)

func ConvertPublicKeyToString(publickey *rsa.PublicKey) string {
	pempublickey := new(pem.Block)
	pempublickey.Type = "RSA PUBLIC KEY"
	pempublickey.Bytes = x509.MarshalPKCS1PublicKey(publickey)
	encodedkey := pem.EncodeToMemory(pempublickey)
	return string(encodedkey)
}

func ReadPrivateKey(file []byte) (*rsa.PrivateKey, error) {
	pem, _ := pem.Decode(file)
	if pem == nil {
		return nil, InvaildPEM
	}
	switch pem.Type {
	case "RSA PRIVATE KEY":
		privatekey, err := x509.ParsePKCS1PrivateKey(pem.Bytes)
		if err != nil {
			return nil, err
		}
		return privatekey, nil
	case "PRIVATE KEY":
		privatekeyinterface, err := x509.ParsePKCS8PrivateKey(pem.Bytes)
		if err != nil {
			return nil, err
		}
		privatekey, ok := privatekeyinterface.(*rsa.PrivateKey)
		if !ok {
			return nil, NotRSAPrivateKey
		}
		return privatekey, nil
	default:
		return nil, NotRSAPrivateKey
	}
}
