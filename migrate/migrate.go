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
	initializers.DB.AutoMigrate(&models.Profile{})
}