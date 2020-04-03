package controllers

import (
	"github.com/gin-gonic/gin"
	"hinze.dev/home/config"
	"hinze.dev/home/models"
	"hinze.dev/home/services"
	"log"
	"net/http"
)

func GetEmail(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if len(token) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	recaptcha, recaptchaError := services.Recaptcha.SiteVerify(token, c.ClientIP())
	if recaptchaError != nil {
		log.Println(recaptchaError)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if recaptcha != nil && recaptcha.Success {
		c.JSON(http.StatusOK, models.GetEmailResponse{
			Email: config.Email,
		})
	} else {
		c.AbortWithStatus(http.StatusForbidden)
	}
}
