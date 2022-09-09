package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func handleError(c *gin.Context, err error) {
	logrus.Errorf("Error occured [error: %s]", err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": http.StatusInternalServerError,
		"msg":    "An error occured internally",
	})
}

func handleInvalidJsonStructure(c *gin.Context, err error) {
	logrus.Errorf("Error occured [error: %s]", err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": http.StatusBadRequest,
		"msg":    "Invalid JSON structure",
	})
}
