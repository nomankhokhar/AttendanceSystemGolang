package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/", helloPoint)
	r.Run(":8080")
}

func helloPoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Attendance is Working",
	})
}
