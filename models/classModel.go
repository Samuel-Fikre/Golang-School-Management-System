package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name" validate:"required,min=2, max=100"`
	Code        string             `json:"code" validate:"required,min=3,max=10"` // Required: Class code, 3-10 characters
	Description *string            `json:"description,omitempty" validate:"max=500"`
	Students    *[]string          `json:"students,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	ClassID     string             `json:"class_id" bson:"class_id"`
	TeacherID   string             `json:"teacher_id" bson:"teacher_id"`
	TimetableID string             `json:"timetable_id" bson:"timetable_id"`
}
