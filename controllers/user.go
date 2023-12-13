package controllers

import (
	"errors"
	"net/http"
	"pose-service/db"
	"pose-service/models"
	"pose-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUsers(c *gin.Context) {
	var users []models.User

	if err := db.GetDB().Find(&users).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.FormatError("No users found"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(
		gin.H{
			"users": users,
		},
	))
}

func GetUserByID(c *gin.Context) {
	var user models.User

	id := c.Param("id")

	if err := db.GetDB().First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.FormatError("The user does not exist"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(
		gin.H{
			"user": user,
		},
	))
}

func CreateUser(c *gin.Context) {
	var input models.CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatError("Validation errors in your request"))
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	user := models.User{Username: input.Username, Password: hashedPassword}
	if err := db.GetDB().Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, utils.FormatError("The user already exists"))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(
		gin.H{
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
			},
		},
	))
}

func UpdateUser(c *gin.Context) {
	var user models.UpdateUserInput

	id := c.Param("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatError("Validation errors in your request"))
		return
	}

	var existingUser models.User

	if err := db.GetDB().First(&existingUser, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.FormatError("The user does not exist"))
		return
	}

	existingUser.Username = user.Username
	existingUser.Password = user.Password

	if err := db.GetDB().Save(&existingUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(
		gin.H{
			"user": existingUser,
		},
	))
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := db.GetDB().First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.FormatError("The user does not exist"))
		return
	}

	if err := db.GetDB().Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Something is broken"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(nil))
}
