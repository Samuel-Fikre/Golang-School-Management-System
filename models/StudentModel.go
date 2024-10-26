package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID   primitive.ObjectID `bson:"_id, omitempty"`
	Name string             `json:"name" validate:"required,min=2,max=100"`
	Age  int                `json:"age" validate:"required,min=2,max=100"`
	// Pointer (*string): This makes Email optional, allowing it to be nil if the user doesn’t provide an email.
	// Validation (validate:"email"): This doesn’t make the field required but only validates the field if it is present. If Email is included, it must match a valid email pattern.
	Email      *string   `json:"email,omitempty" validate:"email"`
	Phone      *string   `json:"phone,omitempty"`
	Enrolled   bool      `json:"enrolled"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
	Student_ID string    `bson:"student_id"  json:"student_id"`
	Class_ID   string    `bson:"class_id" json:"class_id"`
}
