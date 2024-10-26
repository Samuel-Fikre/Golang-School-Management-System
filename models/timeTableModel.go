package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Timetable struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`                                                                                                     // Unique identifier for the timetable entry
	Class_ID     string             `bson:"class_id" json:"class_id"`                                                                                          // Required: Reference to the class
	Teacher_ID   string             `bson:"teacher_id" json:"teacher_id"`                                                                                      // Required: Reference to the teacher conducting the class
	DayOfWeek    string             `bson:"day_of_week" json:"day_of_week" validate:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday Sunday"` // Required: Day of the week
	StartTime    time.Time          `bson:"start_time" json:"start_time" validate:"required"`                                                                  // Required: Start time of the class
	EndTime      time.Time          `bson:"end_time" json:"end_time" validate:"required"`                                                                      // Required: End time of the class
	Room         *string            `json:"room,omitempty" validate:"max=100"`                                                                                 // Optional: Room where the class is held, max 100 characters
	Comments     *string            `json:"comments,omitempty" validate:"max=500"`                                                                             // Optional: Additional notes or comments about the class, max 500 characters
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`                                                                                      // Automatically set on creation
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`                                                                                      // Automatically updated on modification
	Timetable_ID string             `bson:"timetable_id" json:"timetable_id"`
}
