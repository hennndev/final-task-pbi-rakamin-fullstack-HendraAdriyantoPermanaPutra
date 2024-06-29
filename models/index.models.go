package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(40);primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(255)" json:"username" validate:"required"`
	Email     string    `gorm:"type:varchar(255);unique" json:"email" validate:"required,email"`
	Password  string    `gorm:"type:varchar(255)" json:"password" validate:"required,min=6"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Photo     Photo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Photo struct {
	ID       uuid.UUID `gorm:"type:char(40);primaryKey" json:"id"`
	Title    string    `gorm:"type:varchar(255)" json:"title"`
	Caption  string    `gorm:"type:varchar(255)" json:"caption"`
	PhotoUrl string    `gorm:"type:varchar(255)" json:"photo_url"`
	UserID   uuid.UUID `gorm:"index"`
}
