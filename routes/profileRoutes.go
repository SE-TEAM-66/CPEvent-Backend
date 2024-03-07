package routes

import (
	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/middleware"
	"github.com/gin-gonic/gin"
)

func ProfileRoutes(r *gin.Engine) {

	// GET
	r.GET("/profile", controllers.ProfileIndex)
	r.GET("/user_profile", middleware.RequireAuth, controllers.ProfileImage)
	r.GET("/profile/:profileID", controllers.ProfileShow)

	// POST
	r.POST("/profile", controllers.ProfileCreate)
	r.POST("/profiles/:profileID/exp", controllers.CreateExperience)

	// PUT
	r.PUT("/profile/:profileID", controllers.ProfileUpdate)

	// DELETE
	r.DELETE("/profile/:id", controllers.ProfileDelete)


	//Edit skills routes
	r.PUT("/profiles/:profileID/exp", controllers.EditExperience)
	r.PUT("/profiles/:profileID/soft-skills", controllers.EditSoftSkill)
	r.PUT("/profiles/:profileID/lang-skills", controllers.EditLangSkill)
	r.PUT("/profiles/:profileID/dataAna", controllers.EditDataAnaSkill)
	
	
	//Get skills routes
	r.GET("/profiles/:profileID/exp", controllers.GetExperience)
	r.GET("/profiles/:profileID/soft-skills", controllers.GetAllSoftSkills)
	r.GET("/profiles/:profileID/lang-skills", controllers.GetAllLanguageSkills)
	r.GET("/profiles/:profileID/dataAna", controllers.GetAllDataAnalysisSkills)
	

}
