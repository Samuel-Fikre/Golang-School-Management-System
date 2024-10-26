package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Teacher struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `json:"name" validate:"required,min=2,max=100"`    // Required: Teacher's name, 2-100 characters
	Age        int                `json:"age" validate:"required,min=21,max=100"`    // Required: Age, valid range from 21 to 100
	Subject    string             `json:"subject" validate:"required,min=2,max=50"`  // Required: Main subject taught, 2-50 characters
	Email      *string            `json:"email,omitempty" validate:"email"`          // Optional: Teacher's email, validated if present
	Phone      string             `json:"phone,omitempty" validate:"omitempty,e164"` // Optional: Contact phone number
	Classes    *[]string          `json:"classes,omitempty"`                         // Optional: List of class codes or IDs the teacher is assigned to
	DateHired  *time.Time         `json:"date_hired,omitempty"`                      // Optional: Date the teacher was hired
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`              // Automatically set on creation
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`              // Automatically updated on modification
	Teacher_ID string             `json:"teacher_id"`
	Class_ID   string             `json:"class_id"`
}
