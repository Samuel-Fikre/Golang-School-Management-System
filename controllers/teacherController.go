package controllers

import (
	"context"
	"log"
	"net/http"
	helper "sms-system/helpers"
	"sms-system/models"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Teacher Controllers
func GetTeachers() gin.HandlerFunc {
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
				{Key: "teacher_list", Value: bson.D{{Key: "$slice", Value: bson.A{
					"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, err := teacherCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing teachers"})
		}

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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var validate = validator.New()

		var teacher models.Teacher

		if err := c.BindJSON(&teacher); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(teacher)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		defer cancel()
		teacher.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		teacher.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		teacher.ID = primitive.NewObjectID()
		teacher.TeacherID = teacher.ID.Hex()

		result, insertErr := teacherCollection.InsertOne(ctx, teacher)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Teacher item was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var teacher models.Teacher
		defer cancel()

		if err := c.BindJSON(&teacher); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		teacherID := c.Param("teacher_id")

		//filter := bson.M{"teacher_id": teacherID}
		var updateObj primitive.D

		if teacher.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: teacher.Name})
		}

		if teacher.Age != 0 {
			updateObj = append(updateObj, bson.E{Key: "age", Value: teacher.Age})
		}
		if teacher.Email != nil {
			updateObj = append(updateObj, bson.E{Key: "email", Value: teacher.Email})
		}
		if teacher.Phone != "" {
			updateObj = append(updateObj, bson.E{Key: "phone", Value: teacher.Phone})
		}

		if teacher.ClassID != "" {
			err := classCollection.FindOne(ctx, bson.M{"class_id": teacher.ClassID}).Decode(&teacher)
			defer cancel()

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "class was not found"})
				return
			}

			updateObj = append(updateObj, bson.E{Key: "class", Value: teacher.ClassID})
		}

		teacher.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: time.Now()})

		upsert := true

		filter := bson.M{"teacher_id": teacherID}

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
