package swagger

import (
	"final-task-pbi-fullstackdev/models"
)

type LoginValue struct {
	Message  string
	ID       string
	Username string
	Photo    models.Photo
}

type ReturnValue struct {
	Message string
}

type GetPhotosValue struct {
	Message string
	Photos  []models.Photo
}
