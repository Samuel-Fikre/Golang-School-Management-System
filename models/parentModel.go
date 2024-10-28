package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Parent struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `json:"name" validate:"required,min=2,max=100"`    // Required: Teacher's name, 2-100 characters
	Email *string            `json:"email,omitempty" validate:"email"`          // Optional: Teacher's email, validated if present
	Phone string             `json:"phone,omitempty" validate:"omitempty,e164"` // Optional: Contact phone number
	// Optional: Date the teacher was hired
	CreatedAt  time.Time `bson:"created_at" json:"created_at"` // Automatically set on creation
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"` // Automatically updated on modification
	Parent_ID  string    `json:"parent_id" bson:"parent_id"`
	Student_ID string    `json:"student_id" bson:"student_id"`
}
