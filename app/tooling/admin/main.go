package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY0LTdlMTItNDFhNi1hZjNlLTdkYTQwMmMzNDEiLCJ0eXAiOiJKV1QifQ.eyJleHAiOjE3MTI0NTc5NDksImlhdCI6MTY4MDgxMzk0OSwiaXNzIjoic2FsZXMgYXBpIHByb2plY3QiLCJzdWIiOiIxMjM0NTYiLCJSb2xlcyI6WyJ1c2VyIiwiYWRtaW4iXX0.jE0sOtaWnB5SMlWB5_PwcaWb_qfmrxs7H2tEZUywLGYAvwXXNN00pO3Xqu4t3nF7QTSggJLLd8Q0M65QZVDmg_wmRjTo_2NdrcoqFdNgtsLyDP5zdd6MbWKmjS-yjLXZJ_5Q-106VrD5g7BRCx-PcB1M61aaxyqdUJbgID9FrRPCY6Wk7Gud2rSD_agdWk-6WtGyW5Om-GwzDZlrU8YYdPQZAzSUeeQ7L-b9nOZLAwG6_Vd6cRtSdKSl18kBBNVEidWPjABqs28ulVy7Q-5HFBQCyay4bc1vLyK0sRttkc8w_yf_Ju6Icn7P9GwlNbBDiU-RKQXOev9LctJemJyouA"
	isValid := ValidateToken(token, privateKey)
	fmt.Println("Token Valid:", isValid)
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
	token.Header["kid"] = "54bb2164-7e12-41a6-af3e-7da402c341"

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

// Validate Token recreats the Claims that were used to generate token. It verifies that the token was signed using our key
func ValidateToken(tokenStr string, privateKey *rsa.PrivateKey) error {
	parser := jwt.Parser{
		ValidMethods: []string{"RS256"},
	}

	var parsedClaims struct {
		jwt.StandardClaims
		Roles []string
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing key id in token header")
		}

		kidID, ok := kid.(string)
		if !ok {
			return nil, errors.New("user token key id must be string")
		}
		fmt.Println("KID", kidID)
		return &privateKey.PublicKey, nil
	}

	parsedToken, err := parser.ParseWithClaims(tokenStr, &parsedClaims, keyFunc)
	if err != nil {
		return fmt.Errorf("parse with claims to validate token: %w", err)
	}

	if !parsedToken.Valid {
		return errors.New("invalid token")
	}

	return nil

}

func main() {
	err := GenKey()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
