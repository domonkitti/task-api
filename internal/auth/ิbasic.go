package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        u, p, ok := c.Request.BasicAuth()
        log.Println("---basic auth---")
        log.Println(u, p, ok)
        if !ok || u != "admin" || p != "1234" {
            c.Writer.Header().Set("WWW-Authenticate", "Basic")
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
    }
}