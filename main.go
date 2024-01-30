package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/SE-TEAM-66/CPEvent-Backend/controllers"
	"github.com/SE-TEAM-66/CPEvent-Backend/initializers"
	"github.com/SE-TEAM-66/CPEvent-Backend/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	// r.GET("/auth", authHandler)
	// r.GET("/auth/callback", callbackHandler)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/getusers", middleware.RequireAuth, controllers.Getusers)
	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30        // 30 days
	isProd := false             // Set to true when serving over https

	store := cookie.NewStore([]byte(key))
	store.Options(sessions.Options{
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
	})
	r.Use(sessions.Sessions("gothic-session", store))

	gothic.Store = store

	goth.UseProviders(
		google.New("945403094249-qlcdv5r0ju6n3a17effe3osffaesub9k.apps.googleusercontent.com", "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", "http://localhost:4000/auth/callback?provider=google"),
	)

	r.GET("/auth/callback", func(c *gin.Context) {
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
			return
		}
		tmpl, err := template.ParseFiles("templates/success.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
			return
		}
		tmpl.Execute(c.Writer, user)
	})

	r.GET("/auth", func(c *gin.Context) {
		if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
			t, _ := template.ParseFiles("templates/index.html")
			t.Execute(c.Writer, gothUser)
		} else {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		}
	})

	r.GET("/", func(c *gin.Context) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err))
			return
		}
		tmpl.Execute(c.Writer, false)
	})

	r.Run()
}

// var googleOauthConfig = &oauth2.Config{
// 	RedirectURL:  "http://localhost:4000/auth/callback", // Set this to your callback URL
// 	ClientID:     "945403094249-qlcdv5r0ju6n3a17effe3osffaesub9k.apps.googleusercontent.com",
// 	ClientSecret: "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",
// 	Scopes:       []string{"profile", "email"},
// 	Endpoint:     google.Endpoint,
// }

// func authHandler(c *gin.Context) {
// 	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
// 	c.Redirect(http.StatusTemporaryRedirect, url)
// }

// func callbackHandler(c *gin.Context) {
// 	code := c.Query("code")
// 	token, err := googleOauthConfig.Exchange(context.TODO(), code)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
// 		return
// 	}

// 	// Use 'token' to make authenticated requests to Google APIs or store it for later use

// 	// Fetch user information from Google UserInfo API
// 	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
// 	userInfoResponse, err := http.Get(userInfoURL + "?access_token=" + token.AccessToken)
// 	defer userInfoResponse.Body.Close()

// 	if err != nil || userInfoResponse.StatusCode != http.StatusOK {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user information"})
// 		return
// 	}

// 	// Decode the user information from the response body
// 	var userInfo map[string]interface{}
// 	if err := json.NewDecoder(userInfoResponse.Body).Decode(&userInfo); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user information"})
// 		return
// 	}

// 	// Now 'userInfo' contains the user information obtained from Google
// 	c.JSON(http.StatusOK, gin.H{"user_info": userInfo})
// }
