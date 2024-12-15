package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"varchar:191"`
	Email    string `gorm:"type:varchar(191);unique"`
	Password string `gorm:"varchar:191"`
	Avatar   string `gorm:"varchar:191"`
}
