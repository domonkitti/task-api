package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-api/internal/auth"
	"task-api/internal/item"
	//"task-api/internal/mylog"
	"task-api/internal/user"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func Logger() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         t := time.Now()
//         c.Set("example", "12345")
// 		c.Set("example2", "1")
//         log.Println("---- Before ----")
//         c.Next()
//         log.Println("---- After ----")
//         latency := time.Since(t)
//         log.Print(latency)
//         status := c.Writer.Status()
//         log.Println(status)
//     }      
// }

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  fmt.Println("FOO: ",os.Getenv("FOO"))
	// Connect database
	db, err := gorm.Open(
		postgres.Open(
			os.Getenv("DATABASE_URL"),
		),
	)
	if err != nil {
		log.Panic(err)
	}

	// Controller
	controller := item.NewController(db)

	// Router
	r := gin.Default()

	config := cors.DefaultConfig()
	// frontend URL
	config.AllowOrigins = []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}
	config.AllowCredentials = true
	
	r.Use(cors.New(config))

	r.GET("/version", func(c *gin.Context) {
		version, err := GetLatestDBVersion(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": version})
	})
	r.GET("/test", func(ctx *gin.Context) {
		fmt.Println("------Test---------")
		value, exists := ctx.Get("example")
		log.Printf("default example? = %v, %T \n", value ,value)
        if exists {
            log.Println("example = ", value)
        } else {
            log.Println("example does not exists")
        }
		// for i :=0; i<20; i++{
		// 	fmt.Println(i)
		// 	time.Sleep(1 * time.Second)
		// }
        ctx.JSON(201, "test response")
    })
	// Register router
	userController := user.NewController(db,os.Getenv("JWT_SECRET"))
	r.POST("/login",userController.Login)
	r.GET("/items/", controller.FindItems)
    items := r.Group("/items")
	items.Use(auth.Guard(os.Getenv("JWT_SECRET"))) //ปิดแปปทำงายยาก
    {
        items.POST("/", controller.CreateItem)
        //items.GET("/", controller.FindItems)
        items.PATCH("/:id", controller.UpdateItemStatus)
		items.GET("/:id", controller.FindItemByID)
		items.PUT("/:id", controller.UpdateIteminfo)
		items.DELETE("/:id", controller.DeleteItem)

    }
	// Start server
	srv := &http.Server{
        Addr:    (os.Getenv("Port")),
        Handler: r.Handler(),
    }
    go func() {
        // service connections
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()
    // Wait for interrupt signal to gracefully shutdown the server with
    // a timeout of 5 seconds.
    quit := make(chan os.Signal, 1)
    // kill (no param) default send syscall.SIGTERM
    // kill -2 is syscall.SIGINT
    // kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutdown Server ...")
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server Shutdown:", err)
    }
    // catching ctx.Done(). timeout of 5 seconds.
    select {
    case <-ctx.Done():
        log.Println("timeout of 5 seconds.")
    }
    log.Println("Server exiting")
}
	
type GooseDBVersion struct {
	ID        int
	VersionID int
	IsApplied bool
	Tstamp    string
}

// TableName overrides the table name used by User to `profiles`
func (GooseDBVersion) TableName() string {
	return "goose_db_version"
}

// GetLatestDBVersion returns the latest applied version from the goose_db_version table.
func GetLatestDBVersion(db *gorm.DB) (int, error) {
	var version GooseDBVersion

	// Query to get the latest version applied
	err := db.Order("version_id desc").Where("is_applied = ?", true).First(&version).Error
	if err != nil {
		return 0, err
	}

	return version.VersionID, nil
}
