package controllers

import (
	"context"
	"log"
	"net/http"
	helper "sms-system/helpers"
	"sms-system/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Timetable Controllers
func GetTimeTables() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := timetableCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listening the timetable"})
		}
		var allTimeTables []bson.M
		if err = result.All(ctx, &allTimeTables); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allTimeTables)
	}
}

func GetTimeTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var timetable models.Timetable
		tableID := c.Param("table_id")

		err := classCollection.FindOne(ctx, bson.M{"table_id": tableID}).Decode(&timetable)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the timetable"})
		}

		c.JSON(http.StatusOK, timetable)
	}
}

func CreateTimeTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Timetable
		var validate = validator.New()

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.Timetable_ID = table.ID.Hex()
		result, insertErr := timetableCollection.InsertOne(ctx, table)
		if insertErr != nil {
			msg := "Time Table was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateTimeTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Timetable
		tableId := c.Param("timetable_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var updateObj primitive.D

		if table.DayOfWeek != "" {
			updateObj = append(updateObj, bson.E{Key: "day_of_week", Value: table.DayOfWeek})
		}

		//if table.StartTime != nil {
		//updateObj = append(updateObj, bson.E{Key: //"start_time", Value: table.StartTime})
		//}

		//if table.EndTime != nil {
		//updateObj = append(updateObj, bson.E{Key: "end_time", Value: table.EndTime})
		//}

		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"table_id": tableId}

		result, err := timetableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := "Failed to update the timetable"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}
