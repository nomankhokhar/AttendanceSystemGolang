package adminController

import (
	"AttendanceSystem/auth"
	"AttendanceSystem/db"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Route to create a new user
func CreateUser(c *gin.Context) {
	email := c.Query("email")

	// Extract token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Remove "Bearer " prefix to get the actual token
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Validate the token and get claims
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "error_details": err.Error()})
		return
	}

	// Check if the email from the token matches the input email
	if claims.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error_details": "Email is invalid"})
		return
	}

	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if a user with the provided email already exists
	collection := db.GetDB().Collection("users")
	filter := bson.M{"email": user.Email}
	var existingUser struct {
		Email string `json:"email"`
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == nil {
		// User already exists
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		// Some other error occurred while querying the database
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Insert the new user
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Route to update a user's password based on email and name
func UpdateUser(c *gin.Context) {
	email := c.Query("email")

	// Extract token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Remove "Bearer " prefix to get the actual token
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Validate the token and get claims
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "error_details": err.Error()})
		return
	}
	if claims.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error_details": "Email is Invalid"})
		return
	}

	var requestData struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	collection := db.GetDB().Collection("users")
	filter := bson.M{"email": requestData.Email, "name": requestData.Name}

	// Check if user exists
	var existingUser bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user"})
		return
	}

	// Update the user's password
	update := bson.M{"$set": bson.M{"password": requestData.Password}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// Route to get all users
func GetAllUsers(c *gin.Context) {
	email := c.Query("email")

	// Extract token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Remove "Bearer " prefix to get the actual token
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Validate the token and get claims
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "error_details": err.Error()})
		return
	}
	if claims.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error_details": "Email is Invalid"})
		return
	}
	collection := db.GetDB().Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	defer cursor.Close(context.TODO())

	var users []bson.M
	if err := cursor.All(context.TODO(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Route to delete a user based on email and name
func DeleteUser(c *gin.Context) {
	email := c.Query("email")

	// Extract token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Remove "Bearer " prefix to get the actual token
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Validate the token and get claims
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "error_details": err.Error()})
		return
	}

	// Ensure the email in the token matches the email in the request
	if claims.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error_details": "Email is invalid"})
		return
	}

	var requestData struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	collection := db.GetDB().Collection("users")
	filter := bson.M{"email": requestData.Email, "name": requestData.Name}

	// Check if the user exists
	var existingUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	err = collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		// User does not exist
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	} else if err != nil {
		// Database error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// User exists, proceed to delete
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Route to get attendance details based on user's email
func GetUserAttendanceDetail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email required"})
		return
	}

	collection := db.GetDB().Collection("attendances")
	filter := bson.M{"email": email}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendance"})
		return
	}
	defer cursor.Close(context.TODO())

	var attendanceRecords []bson.M
	if err := cursor.All(context.TODO(), &attendanceRecords); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse attendance"})
		return
	}

	// Check if attendance records are empty
	if len(attendanceRecords) == 0 {
		c.JSON(http.StatusOK, gin.H{"attendance": nil})
		return
	}

	// Return attendance data if records exist
	c.JSON(http.StatusOK, gin.H{"attendance": attendanceRecords})
}

// Route to update attendance details based on given info
func UpdateUserAttendance(c *gin.Context) {
	email := c.Query("email")

	// Extract token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Remove "Bearer " prefix to get the actual token
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Validate the token and get claims
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "error_details": err.Error()})
		return
	}
	if claims.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error_details": "Email is Invalid"})
		return
	}

	var requestData struct {
		ID            string `json:"_id"`
		AdminEmail    string `json:"admin_email"`
		EmployeeEmail string `json:"employee_email"`
		Authorized    string `json:"authorized"`
		CaughtUp      bool   `json:"caught_up"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(requestData.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := db.GetDB().Collection("attendances")
	filter := bson.M{"_id": objectID, "email": requestData.EmployeeEmail}

	// Check if attendance exists
	var existingAttendance bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&existingAttendance)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Attendance not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check attendance"})
		return
	}

	// Update the attendance
	update := bson.M{
		"$set": bson.M{
			"authorized": requestData.Authorized,
			"caught_up":  requestData.CaughtUp,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance updated successfully", "attendance_id": requestData.ID})
}
