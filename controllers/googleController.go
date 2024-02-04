package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/SE-TEAM-66/CPEvent-Backend/config"
	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {
	googleConfig := config.SetupConfig()
	url := googleConfig.AuthCodeURL("randomstate")

	http.Redirect(c.Writer, c.Request, url, http.StatusSeeOther)
}

func Googlecallback(c *gin.Context) {
	//state
	state := c.Request.URL.Query()["state"][0]
	if state != "randomstate" {
		fmt.Println("state don't match")
		return
	}
	//code
	code := c.Request.URL.Query()["code"][0]

	//conf
	googleConfig := config.SetupConfig()

	//exchange code for token
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Code-Token Exchange Failed")
	}

	//fetch user info
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("User data fetch failed")
	}
	//parse user info
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("JSON parsing fail")
	}
	c.JSON(http.StatusOK, gin.H{
		"message": string(userData),
	})
	fmt.Println(string(userData))
}
