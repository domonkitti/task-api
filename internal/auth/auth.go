package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string, secret string) (string, error) {
	var t *jwt.Token

	t = jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{username},//เราสามารถเก็บค่าไว้มากกว่านี้เยอะจะตั้ง user id บลาๆ
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Second)),
	})

	signedToken, err := t.SignedString([]byte(secret))
	if err != nil {
		log.Println("error signing key")
		return signedToken, err
	}

	return signedToken, nil
}
