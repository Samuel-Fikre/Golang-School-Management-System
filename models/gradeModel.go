package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Grade struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	StudentID    string             `bson:"student_id" json:"student_id"`
	ClassID      string             `bson:"class_id" json:"class_id"`
	TeacherID    string             `bson:"teacher_id" json:"teacher_id"`
	Score        float64            `bson:"score" json:"score" validate:"required,min=0,max=100"`
	Comments     *string            `json:"comments,omitempty" validate:"max=500"`
	DateRecorded time.Time          `bson:"date_recorded" json:"date_recorded"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	GradeID      string             `bson:"grade_id"  json:"grade_id"`
}
