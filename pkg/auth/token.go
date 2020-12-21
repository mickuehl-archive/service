package auth

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateJWTToken creates a token that can be used for JWT authentication / authorization
func CreateJWTToken(secret, realm, clientID, userID, scope string, duration int64) (string, error) {

	a, err := GetSecureJWTMiddleware(realm, secret)
	if err != nil {
		return "", err
	}

	a.Timeout = (time.Duration)(duration*24) * time.Hour

	// the claim
	client := Client{
		ClientID: clientID,
		UserID:   userID,
		Scope:    scope,
	}

	token, _, err := a.TokenGenerator(&client)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetBearerToken extracts the bearer token
func GetBearerToken(c *gin.Context) string {

	auth := c.Request.Header["Authorization"]
	if len(auth) == 0 {
		return ""
	}

	parts := strings.Split(auth[0], " ")
	if len(parts) != 2 {
		return ""
	}

	if parts[0] == "Bearer" {
		return parts[1]
	}

	return ""
}
