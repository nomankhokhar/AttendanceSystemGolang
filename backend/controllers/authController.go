package controllers

import (
	"AttendanceSystem/db"
	"AttendanceSystem/models"
	"context"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func SignUp(c *gin.Context){
	var user models.User

	// Bind the request body to the user variable
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "provide complete data", "error_details": err.Error()})
		return
	}

	// Hash the password before saving it to the database
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "error_details": err.Error()})
		return
	}

	// Create the user
	_, err := models.CreateUser(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "error_details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully","user": user})
}

func Login(c *gin.Context){
	var loginData struct {
		Email string `json:"email"`
		Password string  `json:"password"`
	} 

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{"error": "provide complete data",  "error_details": err.Error()})
		return
	}


	// Find the user by email
	user , err := models.FindUserByEmail(loginData.Email)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email not found", "error_details": err.Error()})
		return
	}

	// Compare the password from the request with the hashed password in the database
	if !user.ComparePassword(loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password", "error_details": err.Error()})
		return
	}

	// If the password is correct, return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Login successful","user": user})
}

func ForgotPassword(c *gin.Context){
	var input struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "provide complete data",  "error_details": err.Error()})
		return
	}

	// Find the user by email
	user := models.User{}
	collection := db.GetDB().Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email not found", "error_details": err.Error()})
	}

	// Generate a reset token and expiry time
	resetToken := uuid.New().String()
	expiry := time.Now().Add(time.Hour * 1)


	// Update the user with the reset token and expiry time
	user.ResetToken = resetToken
	user.TokenExpiry = expiry

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"reset_token": resetToken, "token_expiry": expiry}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "error_details": err.Error()})
		return
	}
	
	resetLink := "http://localhost:8080/reset-password?token=" + resetToken

	go sendEmail(input.Email, resetLink) // send the reset token to the user's email

	c.JSON(http.StatusOK, gin.H{"message": "Reset token sent to your email"})
}

func sendEmail(Email string, resetLink string){
	// owner email
	from := ""
	password := ""

	// Sender data to user.
	to := Email

	subject := "Password Reset of your Attendance System Account"
	body := "Click the link below to reset your password\n" + resetLink

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" + body

	err := smtp.SendMail("smtp.gmail.com:587",
	smtp.PlainAuth("", from, password, "smtp.gmail.com"),
	from, []string{to}, []byte(msg))

	if err != nil {
		log.Fatal(err)
	}
}

func ResetPassword(c *gin.Context){
	var input struct {
		Token string `json:"token"`
		NewPassword string `json:"new_password"`
	} 

	// Bind the request body to the input variable
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input",  "error_details": err.Error()})
		return
	}


	// Retrieve the user with the reset token
	user := models.User{}
	collection := db.GetDB().Collection("users")

	// Find the user by the reset token
	err := collection.FindOne(context.Background(), bson.M{"reset_token": input.Token}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "error_details": err.Error()})
		return
	}

	// Check if the token has expired
	if time.Now().After(user.TokenExpiry) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token expired"})
		return
	}

	// hash the new password
	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "error_details": err.Error()})
	}

	// Update the user's password and clear the reset token
	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.TokenExpiry = time.Time{} // clear the token expiry


	_ , err = collection.UpdateOne(context.Background(), bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"password": user.Password, "reset_token": user.ResetToken, "token_expiry": user.TokenExpiry}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "error_details": err.Error()})
		return
	}

	// Return the updated user
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful", "user": user})
}