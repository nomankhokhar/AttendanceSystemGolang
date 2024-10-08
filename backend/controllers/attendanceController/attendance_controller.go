package attendanceController

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"AttendanceSystem/auth"
	"AttendanceSystem/controllers/userController"
	"AttendanceSystem/db"
)


type UpdateAttendanceInput struct {
	ID            primitive.ObjectID `json:"_id" binding:"required"` // Attendance ID
	AdminEmail    string             `json:"admin_email" binding:"required"`
	EmployeeEmail  string            `json:"employee_email" binding:"required"`
	Authorized    string             `json:"authorized" binding:"required"`
	CaughtUp      bool               `json:"caught_up"`
}

// InsertAttendance handles the creation of a new attendance record
func InsertAttendance(c *gin.Context) {

	// Parse input data
	var input Attendance
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Check if the email from the token matches the email from the input
	if claims.Email != input.Email {
		c.JSON(http.StatusUnauthorized, gin.H{"error_details": "Email is Invalid"})
		return
	}

	// check the user is authorized and exists
	var user *userController.User
	user, err = userController.FindUserByEmail(input.Email)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email not found", "error_details": err.Error()})
		return
	}

	// Convert input date to time.Time
	date, err := time.Parse("2006-01-02", input.Date.Format("2006-01-02")) // Assuming date in format YYYY-MM-DD

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Create Attendance record
	newAttendance := Attendance{
		Email:          user.Email, // Use the logged-in user's email
		Date:           date,
		StartTime:      input.StartTime,
		FinishTime:     input.FinishTime,
		HoursNotWorked: input.HoursNotWorked,
		Reason:         input.Reason,
		Authorized:     input.Authorized,
		TimeToCatchUp:  input.TimeToCatchUp,
		CaughtUp:       input.CaughtUp,
		Sick:           input.Sick,
		TotalHours:     input.TotalHours,
		Task:           input.Task,
	}

	// Insert into database
	err = InsertAttendanceByEmail(&newAttendance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert attendance record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance record added successfully", "record": newAttendance})
}

// GetAttendanceByEmail retrieves attendance records by email
func GetAttendanceByEmail(c *gin.Context) {
	// Get the logged-in user's email from context
    email := c.Query("email")

	

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	

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

	// Check if the email from the token matches the email from the input
	if claims.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error_details": "Email is Invalid"})
		return
	}

	// Fetch attendance records by email
	attendances, err := GetAttendanceEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendance records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendances": attendances})
}

// InsertAttendanceByEmail inserts a new attendance record for the provided email
func InsertAttendanceByEmail(attendance *Attendance) error {
	collection := db.GetDB().Collection("attendances")

	// Inserting attendance into the collection
	_, err := collection.InsertOne(context.Background(), attendance)
	if err != nil {
		return err
	}
	return nil
}

// GetAttendanceEmail retrieves attendance records for a specific email
func GetAttendanceEmail(email string) ([]Attendance, error) {
	var attendances []Attendance
	collection := db.GetDB().Collection("attendances")

	// Query the database for attendance records by email
	cursor, err := collection.Find(context.Background(), bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode the cursor results into the attendances slice
	if err = cursor.All(context.Background(), &attendances); err != nil {
		return nil, err
	}

	return attendances, nil
}

// DeleteAttendanceByEmailAndID deletes an attendance record by email and ID
func DeleteAttendanceByEmailAndID(c *gin.Context){
	email := c.Query("email")
	id := c.Query("id")

	// Check if email and ID are provided
	if email == ""  || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and ID are required"})
		return
	}

		


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

	// Check if the email from the token matches the email from the input
	if claims.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error_details": "Email is Invalid"})
		return
	}
	

	// Delete attendance record by email and ID

	err = deleteAttendanceByEmailAndID(email , id)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Attendance record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendance record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance record deleted successfully"})	
}


// deleteAttendanceByEmailAndID deletes an attendance record by email and ID
func deleteAttendanceByEmailAndID(email string, id string) error {
	collection := db.GetDB().Collection("attendances")

	// Convert the ID from string to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	// Delete the attendance record by email and ID
	result, err := collection.DeleteOne(context.Background(), bson.M{"email": email, "_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete attendance record: %v", err)
	}

	// Check if any documents were deleted
	if result.DeletedCount == 0 {
		return fmt.Errorf("no record found with the given email and ID")
	}

	fmt.Println("Deleted attendance record:", result.DeletedCount)
	return nil
}


// only admin can update attendance
func UpdateAttendanceByEmailAndID(c *gin.Context) {
	// Extract token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Check if the token starts with "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token must start with 'Bearer '"})
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


	
	// Check if the user is admin
	isAdmin := claims.Role == "admin" // Assuming 'Role' is part of the claims

	var input UpdateAttendanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide complete data", "error_details": err.Error()})
		return
	}

	if claims.Email != input.AdminEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error_details": "Email is Invalid"})
		return
	}

	// Retrieve the existing attendance record by ID
	var existingAttendance Attendance
	err = db.GetDB().Collection("attendances").FindOne(c, bson.M{"_id": input.ID}).Decode(&existingAttendance)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance record not found"})
		return
	}

	// Only allow admin to update the attendance record
	if isAdmin {
		// Ensure that the employee email matches the existing attendance record
		if existingAttendance.Email != input.EmployeeEmail {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only update attendance for the specified employee"})
			return
		}

		// Update the Authorized and CaughtUp fields
		existingAttendance.Authorized = input.Authorized
		existingAttendance.CaughtUp = input.CaughtUp
	} else {
		// Regular users cannot update attendance records
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update attendance records"})
		return
	}

	// Update the attendance record in the database
	_, err = db.GetDB().Collection("attendances").UpdateOne(
		c,
		bson.M{"_id": input.ID},
		bson.M{"$set": existingAttendance},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance record updated successfully", "attendance": existingAttendance})
}