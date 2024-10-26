package routes

import (
	controller "sms-system/controllers"

	"github.com/gin-gonic/gin"
)

func ParentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/students/:student_id", controller.GetStudent())
}
