package main

import (
	"log"
	"os"
	"pose-service/controllers"
	"pose-service/db"
	"pose-service/middlewares"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db.Init()

	r := gin.Default()

	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.Login)

	authorized := r.Group("/")
	authorized.Use(middlewares.Middleware)
	{
		authorized.GET("/users", controllers.GetUsers)
		authorized.GET("/users/:id", controllers.GetUserByID)
		authorized.PUT("/users/:id", controllers.UpdateUser)
		authorized.DELETE("/users/:id", controllers.DeleteUser)

		authorized.GET("/poses", controllers.GetPoses)
		authorized.GET("/poses/:id", controllers.GetPoseByID)
		authorized.POST("/poses", controllers.CreatePose)
		authorized.PUT("/poses/:id", controllers.UpdatePose)
		authorized.DELETE("/poses/:id", controllers.DeletePose)
		authorized.POST("/poses/:id/image", controllers.UploadImage)
	}

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
