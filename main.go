package main

import (
	"os"
	"sms-system/database"
	"sms-system/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var studentCollection *mongo.Collection = database.OpenCollection(database.Client, "Student")

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	// Assuming you have a specific database, like "mydb"
	collection := client.Database("mydb").Collection(collectionName)
	return collection
}

func main() {
	port := os.Getenv(("PORT"))

	if port == "" {
		port = "8000"
	}

	router := gin.New()

	router.Use(gin.Logger())

	router.Use(middleware.Authentication)

	routes.AdminRoutes(router)
	routes.TeacherRoutes(router)
	routes.StudentRoutes(router)
	routes.ParentRoutes(router)
	routes.AuthRoutes(router)
	routes.NewsRoutes(router)

	router.Run(":" + port)
}
