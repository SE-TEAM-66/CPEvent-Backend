package main

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
}

func main() {
	groupDB()
}

func groupDB() {
	initializers.DB.AutoMigrate(&models.Group{})
	initializers.DB.AutoMigrate(&models.ReqPosition{})
}
