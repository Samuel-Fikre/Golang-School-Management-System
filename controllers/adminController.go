package controllers

import (
	"context"
	"log"
	"net/http"
	"sms-system/database"
	"sms-system/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin" // Add this import
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Database Collections

var studentCollection *mongo.Collection = database.OpenCollection(database.Client, "students")

var teacherCollection *mongo.Collection = database.OpenCollection(database.Client, "teachers")

var classCollection *mongo.Collection = database.OpenCollection(database.Client, "class")

var parentCollection *mongo.Collection = database.OpenCollection(database.Client, "parents")

var timetableCollection *mongo.Collection = database.OpenCollection(database.Client, "timetable")

// Student Controllers
func GetStudents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "students_list", Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, err := studentCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching the students"})
			return
		}

		var allStudents []bson.M

		if err = result.All(ctx, &allStudents); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allStudents)
	}
}

func GetStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		studentID := c.Param("student_id")
		var student models.Student

		err := studentCollection.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&student)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the student you requested"})
		}

		c.JSON(http.StatusOK, student)

	}
}

func CreateStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var class models.Class
		var student models.Student
		var validate = validator.New()

		if err := c.BindJSON(&student); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(student)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := classCollection.FindOne(ctx, bson.M{"class_id":student.Class_ID}).Decode(&class)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu was not found"})
			return
		}

		student.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		student.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	student.ID = primitive.NewObjectID()
	}
}

func UpdateStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

// Teacher Controllers
func GetTeachers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "teacher_list", Value: bson.D{{Key: "$slice", Value: bson.A{
					"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, err := teacherCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})

		var allTeachers []bson.M

		if err = result.All(ctx, &allTeachers); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allTeachers)

	}
}

func GetTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var teacher models.Teacher
		teacherID := c.Param("teacher_id")

		err := teacherCollection.FindOne(ctx, bson.M{"teacher_id": teacherID}).Decode(&teacher)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error couldnt fetch the teacher"})
		}

		c.JSON(http.StatusOK, teacher)

	}
}

func CreateTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

func UpdateTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

// Class Controllers
func GetClasses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "class_list", Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, err := classCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "error fetching the class list"})
			return
		}

		var allClasses []bson.M

		if err = result.All(ctx, allClasses); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allClasses)
	}
}

func GetClass() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var class models.Class
		classID := c.Param("class_id")

		err := classCollection.FindOne(ctx, bson.M{"class_id": classID}).Decode(&class)

		if err != nil {
     c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while fetching the class"})
		}



		c.JSON(http.StatusOK, class )



	}
}

func CreateClass() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

func UpdateClass() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

// Parent Controllers
func GetParents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_count", Value: 1},
				{Key: "parent_lists", Value: bson.D{{Key: "$slice", Value: bson.A{startIndex, recordPerPage}}}},
			}},
		}

		result, err := parentCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the parents list"})
			return
		}

		var allParents []bson.M

		if err = result.All(ctx, &allParents); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allParents)
	}
}

func GetParent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var parent models.Parent
		parentID := c.Param("parent_id")

		err := classCollection.FindOne(ctx, bson.M{"parent_id": parentID}).Decode(&parent)

		if err != nil {
     c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while fetching the class"})
		}



		c.JSON(http.StatusOK, parent )
	}
}

func CreateParent() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

func UpdateParent() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

// Timetable Controllers
func GetTimeTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_count", Value: 1},
				{Key: "parent_lists", Value: bson.D{{Key: "$slice", Value: bson.A{startIndex, recordPerPage}}}},
			}},
		}

		result, err := timetableCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the timetables"})
			return
		}

		var allTimetables []bson.M

		if err = result.All(ctx, &allTimetables); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allTimetables)
	}
}

func GetTimeTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var timetable models.Timetable
		classID := c.Param("class_id")

		err := classCollection.FindOne(ctx, bson.M{"class_id": classID}).Decode(&)

		if err != nil {
     c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while fetching the timetable"})
		}



		c.JSON(http.StatusOK, parent )
	}
}

func CreateTimeTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}

func UpdateTimeTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation goes here
	}
}
