package controllers

import (
	"sms-system/database"

	"go.mongodb.org/mongo-driver/mongo"
)

// Database Collections

var studentCollection *mongo.Collection = database.OpenCollection(database.Client, "students")

var teacherCollection *mongo.Collection = database.OpenCollection(database.Client, "teachers")

var classCollection *mongo.Collection = database.OpenCollection(database.Client, "class")

var parentCollection *mongo.Collection = database.OpenCollection(database.Client, "parents")

var timetableCollection *mongo.Collection = database.OpenCollection(database.Client, "timetable")
