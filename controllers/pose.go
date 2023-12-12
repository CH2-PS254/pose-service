package controllers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"

	"pose-service/db"
	"pose-service/models"
	"pose-service/utils"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func GetPoses(c *gin.Context) {
	var poses []models.Pose

	if err := db.GetDB().Find(&poses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(
		gin.H{
			"poses": poses,
		},
	))
}

func GetPoseByID(c *gin.Context) {
	var pose models.Pose

	id := c.Param("id")

	if err := db.GetDB().First(&pose, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, utils.FormatError("The pose does not exist"))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(
		gin.H{
			"pose": pose,
		},
	))
}

func CreatePose(c *gin.Context) {
	var pose models.CreatePoseInput

	if err := c.BindJSON(&pose); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var out = make(map[string]string, len(ve))
			for _, fe := range ve {
				field, _ := reflect.TypeOf(models.CreatePoseInput{}).FieldByName(fe.Field())
				jsonName := string(field.Tag.Get("json"))
				out[jsonName] = utils.GetErrorMsg(fe)
			}

			c.JSON(http.StatusBadRequest, utils.FormatFail(out))
			return
		}

		c.JSON(http.StatusBadRequest, utils.FormatError("Validation errors in your request"))
		return
	}

	newPose := models.Pose{Name: pose.Name, Description: pose.Description}

	if err := db.GetDB().Create(&newPose).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	c.JSON(http.StatusCreated, utils.FormatSuccess(
		gin.H{
			"pose": newPose,
		},
	))
}

func UpdatePose(c *gin.Context) {
	var pose models.Pose

	id := c.Param("id")

	if err := db.GetDB().First(&pose, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, utils.FormatError("The pose does not exist"))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	var updatedPose models.UpdatePoseInput

	if err := c.BindJSON(&updatedPose); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var out = make(map[string]string, len(ve))
			for _, fe := range ve {
				field, _ := reflect.TypeOf(models.UpdatePoseInput{}).FieldByName(fe.Field())
				jsonName := string(field.Tag.Get("json"))
				out[jsonName] = utils.GetErrorMsg(fe)
			}

			c.JSON(http.StatusBadRequest, utils.FormatFail(out))
			return
		}

		c.JSON(http.StatusBadRequest, utils.FormatError("Validation errors in your request"))
		return
	}

	if err := db.GetDB().Model(&pose).Updates(models.Pose{Name: updatedPose.Name, Description: updatedPose.Description}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(
		gin.H{
			"pose": pose,
		},
	))
}

func DeletePose(c *gin.Context) {
	var pose models.Pose

	id := c.Param("id")

	if err := db.GetDB().First(&pose, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, utils.FormatError("The pose does not exist"))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	if err := db.GetDB().Delete(&pose).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(nil))
}

func UploadImage(c *gin.Context) {
	var pose models.Pose

	id := c.Param("id")

	if err := db.GetDB().First(&pose, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, utils.FormatError("The pose does not exist"))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatFail("Validation errors in your request"))
		return
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}
	defer client.Close()

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}
	defer f.Close()

	bucketName := "pose-service"

	wc := client.Bucket(bucketName).Object(file.Filename).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}
	if err := wc.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	publicURL := "https://storage.googleapis.com/" + bucketName + "/" + file.Filename

	if err := db.GetDB().Model(&pose).Update("image", publicURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(
		gin.H{
			"image": publicURL,
		},
	))
}
