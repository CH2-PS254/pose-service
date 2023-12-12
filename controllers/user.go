package controllers

import (
	"net/http"
	"pose-service/db"
	"pose-service/models"
	"pose-service/utils"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User

	if err := db.GetDB().Find(&users).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	var user models.User

	id := c.Param("id")

	if err := db.GetDB().First(&user, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var input models.CreateUserInput

	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	user := models.User{Username: input.Username, Password: hashedPassword}

	if err := db.GetDB().Create(&user).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User created successfully",
	})
}

func UpdateUser(c *gin.Context) {
	var user models.UpdateUserInput

	id := c.Param("id")

	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var existingUser models.User

	if err := db.GetDB().First(&existingUser, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	existingUser.Username = user.Username
	existingUser.Password = user.Password

	if err := db.GetDB().Save(&existingUser).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, existingUser)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := db.GetDB().First(&user, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if err := db.GetDB().Delete(&user).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}
