package main

import (
	"net/http"
	"os"

	"github.com/aschbacd/strava-export/controllers"
	"github.com/aschbacd/strava-export/pkg/utils"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	// Load .env file (if exists)
	godotenv.Load()

	// Debug logs
	if os.Getenv("HTTP_DEBUG") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Gin
	r := gin.Default()

	// User interface
	r.HTMLRender = ginview.Default()
	r.Static("/assets", "./assets")

	// Session storage
	store := cookie.NewStore([]byte(utils.GetRandomString(64)))
	r.Use(sessions.Sessions("session", store))

	// OAuth config
	config := &oauth2.Config{
		ClientID:     os.Getenv("STRAVA_CLIENT_ID"),
		ClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		Scopes:       []string{"activity:read_all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/oauth/authorize",
			TokenURL: "https://www.strava.com/oauth/token",
		},
		RedirectURL: os.Getenv("BASE_URL") + "/authenticate",
	}

	authController := controllers.AuthController{OAuthConfig: *config}

	// Unauthenticated routes
	r.GET("/login", authController.GetLoginPage)
	r.GET("/authenticate", authController.AuthenticateUser)
	r.GET("/rate-limit", func(c *gin.Context) {
		c.HTML(http.StatusTooManyRequests, "rate-limit", nil)
	})
	r.GET("/error", func(c *gin.Context) {
		c.HTML(http.StatusInternalServerError, "error", nil)
	})

	// Authenticated routes
	auth := r.Group("")
	auth.Use(authController.AuthMiddleware())
	auth.GET("/", controllers.GetActivitiesPage)
	auth.GET("/export", controllers.ExportData)
	auth.POST("/logout", controllers.Logout)

	r.Run(utils.GetEnv("HTTP_ADDRESS", "localhost") + ":" + utils.GetEnv("HTTP_PORT", "8080"))
}
