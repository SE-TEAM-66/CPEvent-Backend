package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("In middleware")
	//Get cookie of req
	tokenString, err := c.Cookie("Authorization")
	//Decode n validate
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//Check exp
		unixTime := time.Now().Unix()
		unixTimeFloat := float64(unixTime)
		if unixTimeFloat > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//Find user w/ token sub
		var user models.User
		initializers.DB.First(&user, "email = ?", claims["sub"])

		if user.Email == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//attach to req
		c.Set("user", user)
		//cont
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
