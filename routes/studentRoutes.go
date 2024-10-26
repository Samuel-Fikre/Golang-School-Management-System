package routes

import (
	controller "sms-system/controllers"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/student/:student_id", controller.GetStudent())
	incomingRoutes.GET("/timetable/:class_id", controller.GetTimeTable())
	incomingRoutes.GET("/grades/:grade_id", controller.GetGrade())
}
