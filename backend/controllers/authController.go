package controllers

import (
	"AttendanceSystem/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context){
	var user models.User

	// Bind the request body to the user variable
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "provide complete data"})
		return
	}

	// Hash the password before saving it to the database
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Create the user
	_, err := models.CreateUser(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context){
	var loginData struct {
		Email string `json:"email"`
		Password string  `json:"password"`
	} 

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{"error": "provide complete data"})
		return
	}


	// Find the user by email
	user , err := models.FindUserByEmail(loginData.Email)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email not found"})
		return
	}

	// Compare the password from the request with the hashed password in the database
	if !user.ComparePassword(loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// If the password is correct, return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}