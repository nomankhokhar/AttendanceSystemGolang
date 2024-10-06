package main

import (
	"AttendanceSystem/controllers"
	"AttendanceSystem/db"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	mongoClient, err := db.InitDB()
	
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return
	}
	
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	r := gin.Default()

	r.GET("/ping", helloPoint)

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func helloPoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Pong",
	})
}
