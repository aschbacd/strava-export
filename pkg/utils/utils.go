package utils

import (
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ReturnErrorPage redirects to the error page
func ReturnErrorPage(c *gin.Context) {
	c.Redirect(http.StatusFound, "/error")
	c.Abort()
}

// GetEnv returns an environment variable or a default value
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetRandomString returns a random string with a given length
func GetRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}
