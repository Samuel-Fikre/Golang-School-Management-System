package controllers

import (
	"context"
	"log"
	"net/http"
	helper "sms-system/helpers"
	"sms-system/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin" // Add this import
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Student Controllers
func GetStudents() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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

		if err := c.BindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(student)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := classCollection.FindOne(ctx, bson.M{"class_id": student.ClassID}).Decode(&class)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu was not found"})
			return
		}

		student.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		student.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		student.ID = primitive.NewObjectID()
		student.StudentID = student.ID.Hex()

		result, insertErr := studentCollection.InsertOne(ctx, student)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var class models.Class
		var student models.Student

		studentID := c.Param("student_id")
		if err := c.BindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if student.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: student.Name})
		}

		if student.Age != 0 {
			updateObj = append(updateObj, bson.E{Key: "age", Value: student.Age})
		}
		if student.Email != nil {
			updateObj = append(updateObj, bson.E{Key: "email", Value: *student.Email})
		}
		if student.Phone != nil {
			updateObj = append(updateObj, bson.E{Key: "phone", Value: *student.Phone})
		}

		if student.ClassID != "" {
			err := classCollection.FindOne(ctx, bson.M{"class_id": student.ClassID}).Decode(&class)
			defer cancel()

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "class was not found"})
				return
			}

			updateObj = append(updateObj, bson.E{Key: "class", Value: student.ClassID})
		}

		updateObj = append(updateObj, bson.E{Key: "enrolled", Value: student.Enrolled})
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: time.Now()})

		upsert := true

		filter := bson.M{"student_id": studentID}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := studentCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Student Individual update failed"})
		}

		c.JSON(http.StatusOK, result)

	}

}
