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

type AuthController struct {
	OAuthConfig oauth2.Config
}

// GetLoginPage returns the login page
func (ac *AuthController) GetLoginPage(c *gin.Context) {
	// Get session
	session := sessions.Default(c)

	// Redirect if already logged in
	if session.Get("token") != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	// Set oauth state string
	state := utils.GetRandomString(10)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Return login page
	c.HTML(http.StatusOK, "login", gin.H{
		"authURL": ac.OAuthConfig.AuthCodeURL(state),
	})
}

// AuthenticateUser requests an oauth token and stores it in the session
func (ac *AuthController) AuthenticateUser(c *gin.Context) {
	// Check if state is correct
	if c.Request.FormValue("state") != sessions.Default(c).Get("state").(string) {
		logger.Info("invalid state string passed by user")
		utils.ReturnErrorPage(c)
		return
	}

	// Get token from code
	token, err := ac.OAuthConfig.Exchange(context.Background(), c.Request.FormValue("code"))
	if err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Encode token to store in session
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Store token in session
	session := sessions.Default(c)
	session.Set("token", tokenJSON)
	if err := session.Save(); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Redirect to activities page
	c.Redirect(http.StatusFound, "/")
}

// Logout invalidates the token and clears the session
func Logout(c *gin.Context) {
	// Get client from authentication middleware
	client, exists := c.Get("client")
	if !exists {
		logger.Error("client not passed by authentication middleware")
		utils.ReturnErrorPage(c)
		return
	}

	// Invalidate token
	if _, err := client.(*http.Client).Post("https://www.strava.com/oauth/deauthorize", "", nil); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Clear session
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Redirect to login page
	c.Redirect(http.StatusFound, "/login")
}
