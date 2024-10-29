package routes

import (
	controller "sms-system/controllers"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/st/student/:student_id", controller.GetStudent())
	incomingRoutes.GET("/st/timetable/:class_id", controller.GetTimeTable())
	incomingRoutes.GET("/st/grades/:grade_id", controller.GetGrade())
}
