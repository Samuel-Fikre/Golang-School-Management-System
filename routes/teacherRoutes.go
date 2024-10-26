package routes

import (
	controller "sms-system/controllers"

	"github.com/gin-gonic/gin"
)

func TeacherRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/students", controller.GetStudents())
	incomingRoutes.GET("/students/:student_id", controller.GetStudent())
	incomingRoutes.GET("timetabels/:class_id", controller.GetTimeTables())
	incomingRoutes.POST("/grades/:grade_id", controller.AddGrade())
	incomingRoutes.PATCH("/grades/:grade_id", controller.UpdateGrade())

}
