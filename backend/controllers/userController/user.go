package userController

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"AttendanceSystem/db"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	Role         string             `bson:"role"`
	CreatedAt    time.Time          `bson:"created_at"`
	ResetToken   string             `bson:"reset_token,omitempty"`
	TokenExpiry  time.Time          `bson:"token_expiry,omitempty"`
}

func (u *User) HashPassword() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashPassword)
	return nil
}

func (u *User) ComparePassword(inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
	return err == nil
}

func CreateUser(user *User) (*mongo.InsertOneResult, error) {
	collection := db.GetDB().Collection("users")
	user.CreatedAt = time.Now()

	// Check if the email already exists
	existingUser, _ := FindUserByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("User already exists")
	}

	// Insert the user into the database
	return collection.InsertOne(context.Background(), user)
}

func FindUserByEmail(email string) (*User, error) {
	var user User
	collection := db.GetDB().Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return &user, nil
}
