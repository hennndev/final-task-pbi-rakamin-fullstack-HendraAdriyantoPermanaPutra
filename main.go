package main

import (
	"final-task-pbi-fullstackdev/database/dbconfig"
	_ "final-task-pbi-fullstackdev/docs"
	"final-task-pbi-fullstackdev/routes"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	dbconfig.ConnectDB() // connect database mysql
	app := gin.Default() // initialize gin
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	routes.InitRoute(app) // invoke routes users dan photos
	app.Run(":8080")      //run app di port 8080
}
