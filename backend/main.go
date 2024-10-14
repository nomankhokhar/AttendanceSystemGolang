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

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-ClusterQueueLenght, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

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

	r.Use(CORS())

	r.GET("/ping", helloPoint)

	// User Handlers
	r.POST("/signup", userController.SignUp)
	r.POST("/login", userController.Login)

	r.POST("/forgot-password", userController.ForgotPassword)
	r.POST("/reset-password", userController.ResetPassword)

	// User Attendance Handlers
	r.POST("/attendance/insert", attendanceController.InsertAttendance)
	r.GET("/attendance/getattendance", attendanceController.GetAttendanceByEmail)
	r.DELETE("/attendance/deleteattendance", attendanceController.DeleteAttendanceByEmailAndID)
	r.POST("/attendance/updateattendance", attendanceController.UpdateAttendanceByEmailAndID)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func helloPoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Pong",
	})
}
