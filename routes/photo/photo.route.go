package photo

import (
	"final-task-pbi-fullstackdev/controllers/photo"
	"final-task-pbi-fullstackdev/middlewares"

	"github.com/gin-gonic/gin"
)

func PhotoRoute(route *gin.Engine) {
	//photos route harus memakai middleware, untuk memastikan bahwa yang dapat melakukan akses route hanya user yang sedang login
	route.GET("/photos", middlewares.AuthMiddleware(), photo.GetPhotos)
	route.POST("/photos", middlewares.AuthMiddleware(), photo.PostPhoto)
	route.PUT("/photos/:photoId", middlewares.AuthMiddleware(), photo.UpdatePhoto)
	route.DELETE("/photos/:photoId", middlewares.AuthMiddleware(), photo.DeletePhoto)
}
