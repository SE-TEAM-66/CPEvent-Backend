package controllers

import (
	"fmt"
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
	user := models.User{Fname: body.Fname, Lname: body.Lname, Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password1"})
		return
	}
	//Compare
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password2"})
		return
	}
	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	//sign n get completeed encoded token
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create token"})
		return
	}
	//send back
	c.SetSameSite((http.SameSiteLaxMode))
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
	fmt.Println(user.(models.User).Fname)
}
