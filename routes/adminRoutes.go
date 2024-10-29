package routes

import (
	controller "sms-system/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(incomingRoutes *gin.Engine) {
	// GET REQUESTS
	incomingRoutes.GET("/students", controller.GetStudents())
	incomingRoutes.GET("/students/:student_id", controller.GetStudent())
	incomingRoutes.GET("/teachers", controller.GetTeachers())
	incomingRoutes.GET("/teachers/:teacher_id", controller.GetTeacher())
	incomingRoutes.GET("/classes", controller.GetClasses())
	incomingRoutes.GET("/classes/:class_id", controller.GetClass())
	incomingRoutes.GET("/timetables", controller.GetTimeTables())
	incomingRoutes.GET("/timetables/:class_id", controller.GetTimeTable())
	// POST REQUESTS
	incomingRoutes.POST("/students", controller.CreateStudent())
	incomingRoutes.PATCH("/students/:student_id", controller.UpdateStudent())
	incomingRoutes.POST("/teachers", controller.CreateTeacher())
	incomingRoutes.PATCH("/teachers/:teacher_id", controller.UpdateTeacher())
	incomingRoutes.POST("/classes", controller.CreateClass())
	incomingRoutes.PATCH("/classes/:class_id", controller.UpdateClass())
	incomingRoutes.POST("/timetables", controller.CreateTimeTable())
	incomingRoutes.PATCH("/timetables/:timetable_id", controller.UpdateTimeTable())

}
