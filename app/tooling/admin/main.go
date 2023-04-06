package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenKey creates an x509 private/public key for auth token
func GenKey() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	publicKey := privateKey.PublicKey
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("creating private file: %w", err)
	}
	defer privateFile.Close()

	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding to private file: %w", err)
	}

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return fmt.Errorf("marshalling public key: %w", err)
	}

	publicFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("creating public file: %w", err)
	}
	defer publicFile.Close()

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public file: %w", err)
	}

	fmt.Println("public and private key generated")
	return nil
}

func GenToken() error {
	// the generated token expires in a year
	claims := struct {
		jwt.StandardClaims
		Roles []string
	}{
		jwt.StandardClaims{
			Issuer:    "sales api project",
			Subject:   "123456",
			ExpiresAt: time.Now().Add(8790 * time.Hour).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
		[]string{"user", "admin"},
	}

	method := jwt.GetSigningMethod("RS256")
	token := jwt.NewWithClaims(method, claims)

	f, err := os.Open("config/keys/private.pem")
	if err != nil {
		return fmt.Errorf("open key file: %w", err)
	}

	privatePEM, err := io.ReadAll(io.LimitReader(f, 1024*1024))
	if err != nil {
		return fmt.Errorf("read key file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("JWT parse RSA Private key from PEM: %w", err)
	}

	tokenstr, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("token:", tokenstr)
	return nil
}

func main() {
	err := GenToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
