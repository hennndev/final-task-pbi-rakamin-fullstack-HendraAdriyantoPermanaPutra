package user

import (
	"final-task-pbi-fullstackdev/controllers/user"
	"final-task-pbi-fullstackdev/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoute(route *gin.Engine) {
	route.GET("/users/login", user.Login)
	route.POST("/users/register", user.Register)
	//harus memakai middleware, untuk memastikan bahwa yang dapat melakukan akses route hanya user yang sedang login
	route.GET("/users/logout", middlewares.AuthMiddleware(), user.Logout)
	route.PUT("/users/:userId", middlewares.AuthMiddleware(), user.UpdateUser)
	route.DELETE("/users/:userId", middlewares.AuthMiddleware(), user.DeleteUser)
}
