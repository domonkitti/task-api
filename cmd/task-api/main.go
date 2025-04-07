package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"task-api/internal/item"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// โหลด .env จาก root ของโปรเจค (ขึ้นไป 2 ชั้น)
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("FOO: ", os.Getenv("FOO"))

	// Connect database
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")))
	if err != nil {
		log.Panic(err)
	}

	// Controller สำหรับ items (หรือ projects)
	controller := item.NewController(db)

	// สร้าง router ด้วย Gin
	r := gin.Default()

	config := cors.Config{
        AllowAllOrigins:  true,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"*"}, 
        ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }
	r.Use(cors.New(config))

	// Endpoint สำหรับดึง version ของ DB
	r.GET("/version", func(c *gin.Context) {
		version, err := GetLatestDBVersion(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": version})
	})

	// Endpoint สำหรับ items/projects
	projects := r.Group("/projects")
	{
		projects.POST("/create", controller.CreateProject)
		projects.GET("/allrequested", controller.GetAllProjects)
		projects.GET("/", controller.GetSubTasksByProjectID)
		projects.GET("/subtasks/:project_id", controller.GetSubTasksByProjectID)
		projects.POST("/subtasks/addsubtask/:projectId", controller.CreateSubTask)
		projects.DELETE("/subtasks/deletesubtask/:projectId", controller.DeleteSubTask)
		projects.POST("/savebudgetrequests/:subProjectId", controller.CreateBudgetRequests)//หน้า detail
		projects.GET("/getbudgetrequests/:subProjectId", controller.GetBudgetRequestsBySubtaskId)//หน้า detail
		projects.GET("/getbudgetrequestsbyproject/:projectId", controller.GetBudgetRequestsByProjectID)//ทำ summa


	}

	// สร้าง HTTP server ด้วยค่า Port จาก environment variable "Port"
	srv := &http.Server{
		Addr:    os.Getenv("Port"),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // รอรับสัญญาณปิด

	log.Println("Shutting down server...")

	// ตั้ง timeout สำหรับการ shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

type GooseDBVersion struct {
	ID        int
	VersionID int
	IsApplied bool
	Tstamp    string
}

func (GooseDBVersion) TableName() string {
	return "goose_db_version"
}

func GetLatestDBVersion(db *gorm.DB) (int, error) {
	var version GooseDBVersion
	err := db.Order("version_id desc").Where("is_applied = ?", true).First(&version).Error
	if err != nil {
		return 0, err
	}
	return version.VersionID, nil
}
