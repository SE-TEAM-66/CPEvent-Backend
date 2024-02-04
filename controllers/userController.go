package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//Get
	var body struct {
		Fname    string
		Lname    string
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}
	//Hash
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to hash password"})
		return
	}
	//Create
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}
	initializers.DB.First(&user, "email = ?", body.Email)
	user_profile := models.Profile{
		Fname:  body.Fname,
		Lname:  body.Lname,
		Email:  body.Email,
		UserID: user.ID}
	result_profile := initializers.DB.Create(&user_profile)
	if result_profile.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}
	//Res
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {

	//Get
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}
	//Lookup
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}
	//Compare
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}
	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Second * 3600 * 2).Unix(),
	})
	//sign n get completed encoded token
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create token"})
		return
	}
	//send back
	// c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "localhost", false, true)
	c.String(http.StatusOK, tokenString)
}
func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)
	// googleLogoutURL := "https://accounts.google.com/logout"
	// c.Redirect(http.StatusSeeOther, googleLogoutURL)
	c.String(http.StatusOK, "Cookie has been deleted!")
}
func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func Getusers(c *gin.Context) {
	var users []models.User
	initializers.DB.Find(&users)
	//Res
	c.JSON(200, gin.H{
		"users": users,
	})
	return
}
