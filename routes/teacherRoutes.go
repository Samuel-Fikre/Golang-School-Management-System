package routes

import (
	controller "sms-system/controllers"

	"github.com/gin-gonic/gin"
)

func TeacherRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/teacher/students", controller.GetStudents())
	incomingRoutes.GET("/students/:student_id", controller.GetStudent())
	incomingRoutes.GET("/timetables/:class_id", controller.GetTimeTables())
	incomingRoutes.POST("/grades/:grade_id", controller.CreateGrade())
	incomingRoutes.PATCH("/grades/:grade_id", controller.UpdateGrade())

	// New routes for grades
	incomingRoutes.GET("/grades", controller.GetGrades())          // Get all grades
	incomingRoutes.GET("/grades/:grade_id", controller.GetGrade()) // Get a specific grade by ID
}
