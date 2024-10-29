package controllers

import (
	"context"
	"log"
	"net/http"
	helper "sms-system/helpers"
	"sms-system/models"
	"time"

	"github.com/gin-gonic/gin" // Add this import
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetClasses() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := classCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listening the class it"})
		}

		var allClasses []bson.M

		if err = result.All(ctx, &allClasses); err != nil {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the class"})
		}

		c.JSON(http.StatusOK, class)

	}
}

func CreateClass() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var class models.Class
		var validate = validator.New()

		if err := c.BindJSON(&class); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(class)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		defer cancel()

		class.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		class.ID = primitive.NewObjectID()
		// this line converts the newly generated ObjectID (which is used as the primary key for the food item in MongoDB) into a hexadecimal string representation.
		class.Class_ID = class.ID.Hex()
		//class.Category = class.Category

		result, err := classCollection.InsertOne(ctx, class)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "class was not created"})

			return
		}

		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}

func UpdateClass() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var class models.Class
		defer cancel()

		if err := c.BindJSON(&class); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		classID := c.Param("class_id")

		filter := bson.M{"class_id": classID}

		var updateObj primitive.D

		if class.Code != "" {
			updateObj = append(updateObj, bson.E{Key: "code", Value: class.Code})
		}

		if class.Description != nil {
			updateObj = append(updateObj, bson.E{Key: "description", Value: class.Description})
		}

		class.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: class.UpdatedAt})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := classCollection.UpdateOne(
			ctx, filter, bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := "class item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, result)

	}

}
