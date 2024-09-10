package user

import (
	"fmt"
	"net/http"
	"task-api/internal/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	Service Service
}

func NewController(db *gorm.DB, secret string) Controller {
	return Controller{
		Service: NewService(db, secret),
	}
}

func (controller Controller) Login(ctx *gin.Context) {
	var (
		request model.RequestLogin
	)

	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	//เอามาจาก service ตัวlogin ที่จะคืนค่าออกมาเป็นtolem , nill
	token, err := controller.Service.Login(request)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "message": err.Error(),
        })
        return
    }
	//ส่วนของการตั้งค่า Token
    ctx.SetCookie(
        "token",
        fmt.Sprintf("Bearer %v", token), int(30*time.Second),
        "/",
        "localhost",
        false,
        true, // http-only cookie flag
    )

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login succeeed",
	})
}
