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
	// initializers.DB.AutoMigrate(&models.Group{})
	// initializers.DB.AutoMigrate(&models.User{})
	// initializers.DB.AutoMigrate(&models.Profile{})
	// initializers.DB.AutoMigrate(&models.Exp{})


	// initializers.DB.AutoMigrate(&models.Skill{})
	// initializers.DB.AutoMigrate(&models.Soft_skill{})
	// initializers.DB.AutoMigrate(&models.Lang_skill{})

	initializers.DB.AutoMigrate(&models.Tec_skills{})
	// initializers.DB.AutoMigrate(&models.DataAna{})
	initializers.DB.AutoMigrate(&models.DBmanage{})
	// initializers.DB.AutoMigrate(&models.GraphicDesign{})
	// initializers.DB.AutoMigrate(&models.Programming{})
	// initializers.DB.AutoMigrate(&models.WebDev{})


}