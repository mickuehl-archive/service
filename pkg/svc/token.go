package svc

import (
	"strings"

	"github.com/gin-gonic/gin"
)

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
