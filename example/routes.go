package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/txsvc/commons/pkg/services"
)

func setupRoutes() *gin.Engine {
	// a new router
	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// routes to load static assets and templates
	r.Use(static.Serve("/assets/css", static.LocalFile("./public/assets/css", true)))
	r.Use(static.Serve("/assets/javascript", static.LocalFile("./public/assets/javascript", true)))
	r.LoadHTMLGlob("public/templates/*")

	// default static endpoints
	r.GET("/robots.txt", services.RobotsEndpoint)
	r.GET("/ads.txt", services.NullEndpoint)    // FIXME change to the real handler
	r.GET("/humans.txt", services.NullEndpoint) // FIXME change to the real handler

	return r
}
