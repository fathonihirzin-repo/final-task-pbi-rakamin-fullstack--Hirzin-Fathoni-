package initializers

import "crud/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
