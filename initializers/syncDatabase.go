package initializers

import "github.com/SE-TEAM-66/CPEvent-Backend/models"

func SyncDatabase() {
	DB.AutoMigrate(models.User{})
}
