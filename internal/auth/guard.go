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
        tokenString, err := c.Cookie("token")
        if err != nil {
            log.Println("Token missing in cookie")
            c.JSON(http.StatusUnauthorized, gin.H{"message": "กรุณา login"})
            c.Abort()
            return
        }

        // Remove prefix "Bearer " from auth token
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        token, err := verifyToken(tokenString, secret)
        if err != nil {
            log.Printf("Token verification failed: %v\n", err)
            c.JSON(http.StatusUnauthorized, gin.H{"message": "กรุณา login"})
            c.Abort()
            return
        }
        //เอาclaim ออกมา
        claims := token.Claims.(jwt.MapClaims)

        //log.Printf("ค่าใน: %+v\n", claims)

        // ตรวจสอบและแปลงค่าของ "aud"
        var username string
        switch aud := claims["aud"].(type) {
        case string://ไม่เข้าใจตรงเช็ค เคสเนี่ยเหละ
            // กรณีที่ "aud" เป็น string เดี่ยว
            username = aud
        case []interface{}:
            // กรณีที่ "aud" เป็น slice ของ interface{} (กรณีที่มี audience หลายค่า)
            if len(aud) > 0 {
                username = aud[0].(string) // ใช้ค่าตัวแรกของ slice
                //log.Printf("Username: %s\n", username)
            }
        default:
            log.Println("Invalid 'aud' field in token")
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
            c.Abort()
            return
        }

        // เก็บ username ลงใน context ของ Gin
        c.Set("username", username)
        log.Printf("Username set in context: %s\n", username) // Log การเก็บ username

        // ส่งต่อไปยัง handler ถัดไป
        c.Next()
    }
}

// verifyToken ฟังก์ชันตรวจสอบ JWT ว่าถูกต้องหรือไม่
func verifyToken(tokenString string, secret string) (*jwt.Token, error) {
    log.Println("Verifying token...") // Log ว่าเริ่มตรวจสอบ JWT

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // ตรวจสอบวิธีการลงชื่อว่าเป็น HS256 หรือไม่
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            log.Printf("Unexpected signing method: %v\n", token.Header["alg"])
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        log.Println("Signing method is valid") // Log การตรวจสอบวิธีการลงชื่อ
        return []byte(secret), nil
    })

    if err != nil {
        log.Printf("Token verification error: %v\n", err)
    } else {
        log.Println("Token verified successfully") // Log เมื่อตรวจสอบ JWT สำเร็จ
    }

    return token, err
}
