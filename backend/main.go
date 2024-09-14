package main

import (
	"AttendanceSystem/db"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	mongoClient := db.InitDB()
	defer mongoClient.Disconnect(context.Background())
	r := gin.Default()

	r.GET("/ping", helloPoint)
	r.Run(":8080")
}

func helloPoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Pong",
	})
}
