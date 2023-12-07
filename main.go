package main

import (
	"pose-service/controllers"
	"pose-service/db"

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

	r.Run("localhost:8080")
}
