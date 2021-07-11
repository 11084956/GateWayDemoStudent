package public

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SignKey = "my_sign_key"

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.StandardClaims
}

func Decode(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SignKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(MyCustomClaims)
	if !ok {
		return "", err
	}

	if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
		return "", errors.New("request expired")
	}

	if claims.Foo != "test" {
		return "", errors.New("sign foo error")
	}

	return claims.Foo, nil
}

func Encode(foo string) (string, error) {
	mySignKey := []byte(SignKey)

	claims := MyCustomClaims{
		Foo: foo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 20).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(mySignKey)
}
