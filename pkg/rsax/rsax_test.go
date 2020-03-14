package rsax

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestReadPrivateKey(t *testing.T) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
	}

	pkcs1key := x509.MarshalPKCS1PrivateKey(privatekey)
	pkcs8key, err := x509.MarshalPKCS8PrivateKey(privatekey)

	tests := []struct {
		Block    pem.Block
		Equalnil bool
		Name     string
	}{
		{
			Name:     "Vaild PKCS1 key",
			Equalnil: true,
			Block: pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: pkcs1key,
			},
		},
		{
			Name:     "Vaild PKCS8 key",
			Equalnil: true,
			Block: pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: pkcs8key,
			},
		},
		{
			Name:     "invaild pem.Block.Type(PKCS1)",
			Equalnil: false,
			Block: pem.Block{
				Type:  "RSAPRIVATEKEY",
				Bytes: pkcs1key,
			},
		},
		{
			Name:     "invaild pem.Block.Type(PKCS8)",
			Equalnil: false,
			Block: pem.Block{
				Type:  "PRIVATEKEY",
				Bytes: pkcs8key,
			},
		},
		{
			Name:     "invaild pem.Block.Bytes(PKCS1)",
			Equalnil: false,
			Block: pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: pkcs8key,
			},
		},
		{
			Name:     "invaild pem.Block.Bytes(PKCS8)",
			Equalnil: false,
			Block: pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: pkcs1key,
			},
		},
	}

	for _, test := range tests {
		pemfile := pem.EncodeToMemory(&test.Block)
		_, err := ReadPrivateKey(pemfile)
		switch err {
		case nil:
			if !test.Equalnil {
				t.Errorf("%s: %v", test.Name, err)
			}
		default:
			if test.Equalnil {
				t.Errorf("%s: %v", test.Name, err)
			}
		}
	}
}
