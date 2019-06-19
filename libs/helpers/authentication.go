package helpers

import (
	"crypto/rsa"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// TODO: move all of these functions to the user services

// GenerateToken returns a fresh token based on the given private RSA key
func GenerateToken(key *rsa.PrivateKey) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	return token.SignedString(key)
}

// ParsePublicKey returns the RSA public key contained in the given file
func ParsePublicKey(path string) *rsa.PublicKey {
	var (
		err error
		key *rsa.PublicKey
		str []byte
	)

	str, err = ioutil.ReadFile(path)
	HandleError(err)
	key, err = jwt.ParseRSAPublicKeyFromPEM(str)
	HandleError(err)

	return key
}

// ParsePrivateKey returns the RSA private key contained in the given file
func ParsePrivateKey(path string) *rsa.PrivateKey {
	var (
		err error
		key *rsa.PrivateKey
		str []byte
	)

	str, err = ioutil.ReadFile(path)
	HandleError(err)
	key, err = jwt.ParseRSAPrivateKeyFromPEM(str)
	HandleError(err)

	return key
}
