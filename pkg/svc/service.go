package svc

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/service/pkg/static"
)

type (
	// APIService abstracts an API endpoint
	APIService struct {
		router *gin.Engine
		addr   string
	}
)

var service *APIService

func init() {
	// basic http stack config
	gin.DisableConsoleColor()
	// make sure the service is not initialized
	service = nil
}

// New creates a new service instance on addr
func New(addr ...string) error {
	localAddr := "0.0.0.0:8080"
	if addr == nil && len(addr) > 0 {
		localAddr = addr[0]
	}
	if service != nil {
		return errors.New(fmt.Sprintf("A router for this address already exists: %s", localAddr))
	}

	// a new router
	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	s := &APIService{
		router: r,
		addr:   localAddr,
	}
	service = s

	return nil
}

// AddDefaultEndpoints adds a couple of simple handlers to the router
func AddDefaultEndpoints() {
	// default static endpoints
	service.router.GET("/robots.txt", RobotsEndpoint)
	service.router.GET("/ads.txt", NullEndpoint)    // FIXME change to the real handler
	service.router.GET("/humans.txt", NullEndpoint) // FIXME change to the real handler
}

// ServeStaticAssets adds handlers to serve static assets
func ServeStaticAssets(path, dir string) {
	// routes to load static assets and templates
	service.router.Use(static.Serve(path, static.LocalFile(dir, true)))
}

// Start attaches the router to a server. This function blocks!
func Start() {
	service.router.Run(service.addr)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func GET(relativePath string, handler gin.HandlerFunc) gin.IRoutes {
	return service.router.Handle(http.MethodGet, relativePath, handler)
	//return service.router.Handle(http.MethodGet, relativePath, handlers)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func POST(relativePath string, handler gin.HandlerFunc) gin.IRoutes {
	return service.router.Handle(http.MethodPost, relativePath, handler)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func DELETE(relativePath string, handler gin.HandlerFunc) gin.IRoutes {
	return service.router.Handle(http.MethodDelete, relativePath, handler)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func PATCH(relativePath string, handler gin.HandlerFunc) gin.IRoutes {
	return service.router.Handle(http.MethodPatch, relativePath, handler)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func PUT(relativePath string, handler gin.HandlerFunc) gin.IRoutes {
	return service.router.Handle(http.MethodPut, relativePath, handler)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func OPTIONS(relativePath string, handler gin.HandlerFunc) gin.IRoutes {
	return service.router.Handle(http.MethodOptions, relativePath, handler)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func HEAD(relativePath string, handler gin.HandlerFunc) gin.IRoutes {
	return service.router.Handle(http.MethodHead, relativePath, handler)
}
