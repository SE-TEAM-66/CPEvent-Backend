package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SE-TEAM-66/CPEvent-Backend/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	tokeng, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Code-Token Exchange Failed")
	}

	//fetch user info
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tokeng.AccessToken)
	if err != nil {
		fmt.Println("User data fetch failed")
	}
	//parse user info
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("JSON parsing fail")
	}
	email := strings.Split(string(userData), ",")[1]
	emailss := strings.Replace(email, " ", "", -1)
	emails := strings.Replace(strings.Split(emailss, ":")[1], "\"", "", -1)
	fmt.Println(emails)
	fmt.Print("Local Cookie Created!")
	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": emails,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	//sign n get completed encoded token
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create token"})
		return
	}
	//send back
	c.SetSameSite((http.SameSiteLaxMode))
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
	fmt.Print("Local Cookie Created!")
}
