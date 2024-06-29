package routes

import (
	"final-task-pbi-fullstackdev/routes/photo"
	"final-task-pbi-fullstackdev/routes/user"

	"github.com/gin-gonic/gin"
)

func InitRoute(app *gin.Engine) {
	route := app
	user.UserRoute(route)
	photo.PhotoRoute(route)
}
