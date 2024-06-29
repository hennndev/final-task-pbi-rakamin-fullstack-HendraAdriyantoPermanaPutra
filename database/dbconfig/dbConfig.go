package dbconfig

import (
	"log"

	"final-task-pbi-fullstackdev/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/finaltask_pbi_fullstackdev?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Database no connected")
	}
	log.Println("Database connected")

	db.AutoMigrate(&models.User{}, &models.Photo{})
	DB = db
}
