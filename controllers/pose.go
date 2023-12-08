package controllers

import (
	"context"
	"io"
	"net/http"

	"pose-service/db"
	"pose-service/models"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPoses(c *gin.Context) {
	var poses []models.Pose

	if err := db.GetDB().Find(&poses).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, poses)
}

func GetPoseByID(c *gin.Context) {
	var pose models.Pose

	id := c.Param("id")

	if err := db.GetDB().First(&pose, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Pose not found"})
			return
		}

		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, pose)
}

func CreatePose(c *gin.Context) {
	var pose models.CreatePoseInput

	if err := c.BindJSON(&pose); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	newPose := models.Pose{Name: pose.Name, Description: pose.Description}

	if err := db.GetDB().Create(&newPose).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newPose)
}

func UpdatePose(c *gin.Context) {
	var pose models.Pose

	id := c.Param("id")

	if err := db.GetDB().First(&pose, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Pose not found"})
			return
		}

		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var updatedPose models.UpdatePoseInput

	if err := c.BindJSON(&updatedPose); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := db.GetDB().Model(&pose).Updates(models.Pose{Name: updatedPose.Name, Description: updatedPose.Description}).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, pose)
}

func DeletePose(c *gin.Context) {
	var pose models.Pose

	id := c.Param("id")

	if err := db.GetDB().First(&pose, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Pose not found"})
			return
		}

		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err := db.GetDB().Delete(&pose).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Pose deleted successfully"})
}

func UploadImage(c *gin.Context) {
	var pose models.Pose

	id := c.Param("id")

	if err := db.GetDB().First(&pose, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Pose not found"})
			return
		}

		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create storage client"})
		return
	}
	defer client.Close()

	f, err := file.Open()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to open file"})
		return
	}
	defer f.Close()

	bucketName := "pose-service"

	wc := client.Bucket(bucketName).Object(file.Filename).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to copy file to bucket"})
		return
	}
	if err := wc.Close(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to close bucket writer"})
		return
	}

	publicURL := "https://storage.googleapis.com/" + bucketName + "/" + file.Filename

	if err := db.GetDB().Model(&pose).Update("image", publicURL).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}
