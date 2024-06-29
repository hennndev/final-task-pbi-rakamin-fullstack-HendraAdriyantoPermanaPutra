package user

import (
	"final-task-pbi-fullstackdev/controllers/user"

	"github.com/gin-gonic/gin"
)

func UserRoute(route *gin.Engine) {
	route.GET("/users/logout", user.Logout)
	route.GET("/users/login", user.Login)
	route.POST("/users/register", user.Register)
	route.PUT("/users/:userId", user.UpdateUser)
	route.DELETE("/users/:userId", user.DeleteUser)
}
