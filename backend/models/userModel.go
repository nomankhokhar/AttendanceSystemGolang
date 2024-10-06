package models

import (
	"AttendanceSystem/db"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`      
	Email     string             `bson:"email"`     
	Password  string             `bson:"password"`  
	CreatedAt time.Time          `bson:"created_at"`
	ResetToken    string         `bson:"reset_token,omitempty"`
	TokenExpiry   time.Time      `bson:"token_expiry,omitempty"`
}

func (u *User) HashPassword() error {
	hashPassword , err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashPassword)
	return nil
}

// ComparePassword compares the input password with the hashed password
func (u *User) ComparePassword(inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
	return err == nil
}

// CreateUser creates a new user
func CreateUser(user *User) (*mongo.InsertOneResult, error) {
	collection := db.GetDB().Collection("users")
	user.CreatedAt = time.Now()

	// Check if the email already exists
	exitingUser, _ := FindUserByEmail(user.Email)
	if exitingUser != nil {
		return nil, errors.New("User already exists")
	}

	// insert the user into the database
	return collection.InsertOne(context.Background(), user)
}


// FindUserByEmail finds a user by their email
func FindUserByEmail(email string) (*User, error) {
	var user User
	collection := db.GetDB().Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return &user, nil
}