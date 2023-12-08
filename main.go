package main

import (
	"log"
	"os"
	"pose-service/controllers"
	"pose-service/db"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db.Init()

	r := gin.Default()

	r.GET("/poses", controllers.GetPoses)
	r.GET("/poses/:id", controllers.GetPoseByID)
	r.POST("/poses", controllers.CreatePose)
	r.PUT("/poses/:id", controllers.UpdatePose)
	r.DELETE("/poses/:id", controllers.DeletePose)
	r.POST("/poses/:id/image", controllers.UploadImage)

	var address string
	if runtime.GOOS == "windows" {
		address = "localhost"
	} else {
		address = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	r.Run(address + ":" + port)
}
