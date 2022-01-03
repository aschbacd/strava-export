package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aschbacd/strava-export/pkg/logger"
	"github.com/aschbacd/strava-export/pkg/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// AuthMiddleware checks if a user is authenticated
func (a *AuthController) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from session storage
		tokenJSON := sessions.Default(c).Get("token")
		if tokenJSON == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Parse json token
		var token oauth2.Token
		if err := json.Unmarshal(tokenJSON.([]byte), &token); err != nil {
			logger.Error(err.Error())
			utils.ReturnErrorPage(c)
			return
		}

		// Set token source and client
		c.Set("tokenSource", a.OAuthConfig.TokenSource(context.Background(), &token))
		c.Set("client", a.OAuthConfig.Client(context.Background(), &token))

		c.Next()
	}
}
