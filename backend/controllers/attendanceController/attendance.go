package attendanceController

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Attendance struct represents the structure for attendance records
type Attendance struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `json:"email" binding:"required"` // email of the user
	Date           time.Time          `json:"date" binding:"required"`
	StartTime      string             `json:"start_time" binding:"required"`
	FinishTime     string             `json:"finish_time" binding:"required"`
	HoursNotWorked float64            `json:"hours_not_worked"`
	Reason         string             `json:"reason"`
	Authorized     string             `json:"authorized" binding:"required"`
	TimeToCatchUp  float64            `json:"time_to_catch_up"`
	CaughtUp       bool               `json:"caught_up"`
	Sick           bool               `json:"sick"`
	TotalHours     float64            `json:"total_hours" binding:"required"`
	Task           string             `json:"task"`
}
