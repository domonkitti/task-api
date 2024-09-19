package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Guard(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token "Bearer xxx" from cookie
		auth, err := c.Cookie("token")
		if err != nil {
			log.Println("Token missing in cookie")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "กรุณา login"})
			c.Abort()
			return
		}

		// Remove prefix "Bearer " from auth token
		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := verifyToken(tokenString, secret)
		if err != nil {
			log.Printf("Token verification failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "กรุณา login"})
			c.Abort()
			return
		}

		log.Printf("Token verified successfully. Claims: %+v\n", token.Claims)
	}
}
func verifyToken(tokenString string, secret string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Return the verified token
	return token, nil
}
