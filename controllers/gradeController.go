package controllers

import (
	"context"
	"log"
	"net/http"
	"sms-system/database"
	"sms-system/models"
	"time"

	"github.com/gin-gonic/gin" // Add this import
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var gradeCollection *mongo.Collection = database.OpenCollection(database.Client, "grade")

func GetGrades() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := gradeCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the grades"})
		}
		var allGrades []bson.M
		if err = result.All(ctx, &allGrades); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allGrades)
	}
}

func GetGrade() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		GradeID := c.Param("grade_id")
		var grade models.Grade

		err := gradeCollection.FindOne(ctx, bson.M{"orderItem_id": GradeID}).Decode(&grade)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the grade"})
		}
		c.JSON(http.StatusOK, grade)
	}
}

func GetGradelistByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		classID := c.Param("class_id")

		allGradeLists, err := ItemsByGrade(classID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing grade items by order"})
		}

		c.JSON(http.StatusOK, allGradeLists)
	}
}
func ItemsByGrade(classID string) (classGrades []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "class_id", Value: classID}}}}

	// Lookup for Student data
	lookUpStudentStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "student"},
		{Key: "localField", Value: "student_id"},
		{Key: "foreignField", Value: "student_id"},
		{Key: "as", Value: "student"},
	}}}

	unwindStudentStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$student"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	// Lookup for Class data
	lookUpClassStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "class"},
		{Key: "localField", Value: "class_id"},
		{Key: "foreignField", Value: "class_id"},
		{Key: "as", Value: "class"},
	}}}

	unwindClassStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$class"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	// Project relevant fields, including student and class info
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "student_name", Value: "$student.name"},
		{Key: "class_name", Value: "$class.name"},
		{Key: "grade_id", Value: "$grade_id"},
		{Key: "score", Value: "$score"},
		{Key: "comments", Value: "$comments"},
		{Key: "date_recorded", Value: "$date_recorded"},
	}}}

	// Group grades by ClassID and accumulate student grade details
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$class_id"},
		{Key: "class_id", Value: bson.D{{Key: "$first", Value: "$class_id"}}},
		{Key: "class_name", Value: bson.D{{Key: "$first", Value: "$class.name"}}},
		{Key: "grades", Value: bson.D{{Key: "$push", Value: bson.D{
			{Key: "student_name", Value: "$student_name"},
			{Key: "grade_id", Value: "$grade_id"},
			{Key: "score", Value: "$score"},
			{Key: "comments", Value: "$comments"},
			{Key: "date_recorded", Value: "$date_recorded"},
		}}}},
		{Key: "total_grades", Value: bson.D{{Key: "$sum", Value: 1}}},
	}}}

	cursor, err := gradeCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookUpStudentStage,
		unwindStudentStage,
		lookUpClassStage,
		unwindClassStage,
		projectStage,
		groupStage,
	})
	if err != nil {
		return nil, err
	}

	// Collect results
	if err = cursor.All(ctx, &classGrades); err != nil {
		return nil, err
	}

	return classGrades, nil
}

func CreateGrade() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var grade models.Grade
		var validate = validator.New()

		// Bind incoming JSON to the grade model
		if err := c.BindJSON(&grade); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Example: Retrieve Class based on some identifier (e.g., class name)
		classID := grade.ClassID
		var class models.Class
		if err := classCollection.FindOne(ctx, bson.M{"_id": classID}).Decode(&class); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Class not found"})
			return
		}

		// Example: Retrieve Student based on student_id provided
		studentID := grade.StudentID
		var student models.Student
		if err := studentCollection.FindOne(ctx, bson.M{"_id": studentID}).Decode(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Student not found"})
			return
		}

		// Perform validation on the grade struct
		if validationErr := validate.Struct(grade); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Set additional fields before inserting the grade
		grade.ID = primitive.NewObjectID()
		grade.CreatedAt = time.Now()
		grade.UpdatedAt = time.Now()
		grade.GradeID = grade.ID.Hex()
		grade.DateRecorded = time.Now()

		// Insert grade into the database
		GradesToBeInserted := []interface{}{grade}
		insertedGradeItems, err := gradeCollection.InsertMany(ctx, GradesToBeInserted)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Respond with the inserted grade
		c.JSON(http.StatusOK, insertedGradeItems)
	}
}

func UpdateGrade() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		gradeID := c.Param("grade_id")
		var grade models.Grade
		var validate = validator.New()

		// Bind incoming JSON to the grade model to validate new data
		if err := c.BindJSON(&grade); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Perform validation on the grade struct
		if validationErr := validate.Struct(grade); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Convert gradeID from string to ObjectID type
		objectID, err := primitive.ObjectIDFromHex(gradeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade ID"})
			return
		}

		// Define the update object with fields to be updated
		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "score", Value: grade.Score},
				{Key: "comments", Value: grade.Comments},
				{Key: "updated_at", Value: time.Now()},
				// Add more fields here if necessary, e.g., "student_id" or "class_id"
			}},
		}

		// Perform the update operation
		result, err := gradeCollection.UpdateOne(
			ctx,
			bson.M{"_id": objectID},
			updateObj,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating grade"})
			return
		}

		// Check if any document was updated
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Grade not found"})
			return
		}

		// Respond with a success message
		c.JSON(http.StatusOK, gin.H{"message": "Grade updated successfully"})
	}
}
