package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

// type Photo struct {
// 	gorm.Model
// 	ID       string
// 	Title    string
// 	Caption  string
// 	PhotoUrl string
// 	UserID   int
// }
