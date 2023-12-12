package controllers

import (
	"errors"
	"net/http"
	"pose-service/db"
	"pose-service/models"
	"pose-service/utils"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.BindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var out = make(map[string]string, len(ve))
			for _, fe := range ve {
				field, _ := reflect.TypeOf(models.LoginInput{}).FieldByName(fe.Field())
				jsonName := string(field.Tag.Get("json"))
				out[jsonName] = utils.GetErrorMsg(fe)
			}

			c.JSON(http.StatusBadRequest, utils.FormatFail(out))
			return
		}

		c.JSON(http.StatusBadRequest, utils.FormatError("Validation errors in your request"))
		return
	}

	var user models.User

	if err := db.GetDB().Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatError("The user does not exist"))
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusBadRequest, utils.FormatError("Invalid password"))
		return
	}

	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.FormatError("Failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccess(
		gin.H{
			"username": user.Username,
			"token":    token,
		},
	))
}
