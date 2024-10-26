package adminController

import (
	"AttendanceSystem/db"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Route to create a new user
func CreateUser(c *gin.Context) {
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

	collection := db.GetDB().Collection("attendances")
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Route to update a user's password based on email and name
func UpdateUser(c *gin.Context) {
	var requestData struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	collection := db.GetDB().Collection("attendances")
	filter := bson.M{"email": requestData.Email, "name": requestData.Name}
	update := bson.M{"$set": bson.M{"password": requestData.Password}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// Route to get all users
func GetAllUsers(c *gin.Context) {
	collection := db.GetDB().Collection("attendances")
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
	var requestData struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	collection := db.GetDB().Collection("attendances")
	filter := bson.M{"email": requestData.Email, "name": requestData.Name}

	_, err := collection.DeleteOne(context.TODO(), filter)
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

	c.JSON(http.StatusOK, attendanceRecords)
}

// Route to update attendance details based on given info
func UpdateUserAttendance(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "Attendance updated successfully"})
}
