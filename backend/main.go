package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"AttendanceSystem/controllers/attendanceController"
	"AttendanceSystem/controllers/userController"
	"AttendanceSystem/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
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

	// User Handlers
	r.POST("/signup", userController.SignUp)
	r.POST("/login", userController.Login)
	r.POST("/forgot-password", userController.ForgotPassword)
	r.POST("/reset-password", userController.ResetPassword)


	// User Attendance Handlers

	r.POST("/attendance", attendanceController.InsertAttendance) 
	r.GET("/attendance", attendanceController.GetAttendanceByEmail) 

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func helloPoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Pong",
	})
}
