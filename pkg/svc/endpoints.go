package svc

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NullEndpoint just respondes with an empty 200
func NullEndpoint(c *gin.Context) {
	c.Status(http.StatusAccepted)
}

// RobotsEndpoint maps to GET /robots.txt
func RobotsEndpoint(c *gin.Context) {
	// simply write text back ...
	c.Header("Content-Type", "text/plain")

	// a simple robots.txt file, disallow the API
	c.Writer.Write([]byte("User-agent: *\n\n"))
	c.Writer.Write([]byte("Disallow: /a/\n")) // FIXME make this configurable, e.g. using ROBOTS_TXT env
}

// StandardAPIResponse is the default way to respond to API requests
func StandardAPIResponse(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
	}
}

// StandardNoRouteResponse is the default response to unknown routes
func StandardNoRouteResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"status": "error", "msg": "page not found"})
}

// StandardJSONResponse is the default way to respond to API requests
func StandardJSONResponse(c *gin.Context, res interface{}, err error) {
	if err == nil {
		if res == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusOK, res)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
	}
}
