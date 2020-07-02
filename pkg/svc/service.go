package svc

import (
	"fmt"

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

var routes map[string]*gin.Engine

func init() {
	// initialize the registry
	routes = make(map[string]*gin.Engine)
	// basic http stack config
	gin.DisableConsoleColor()
}

// New creates a new service instance on addr
func New(addr ...string) (*APIService, error) {
	localAddr := "0.0.0.0:8080"
	if addr == nil && len(addr) > 0 {
		localAddr = addr[0]
	}
	if _, ok := routes[localAddr]; ok {
		return nil, errors.New(fmt.Sprintf("A router for this address already exists: %s", localAddr))
	}

	// a new router
	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	// register the router
	routes[localAddr] = r

	s := &APIService{
		router: r,
		addr:   localAddr,
	}

	return s, nil
}

// AddDefaultEndpoints adds a couple of simple handlers to the router
func (api *APIService) AddDefaultEndpoints() {
	// default static endpoints
	api.router.GET("/robots.txt", RobotsEndpoint)
	api.router.GET("/ads.txt", NullEndpoint)    // FIXME change to the real handler
	api.router.GET("/humans.txt", NullEndpoint) // FIXME change to the real handler
}

// ServeStaticAssets adds handlers to serve static assets
func (api *APIService) ServeStaticAssets(path, dir string) {
	// routes to load static assets and templates
	api.router.Use(static.Serve(path, static.LocalFile(dir, true)))
}

// Start attaches the router to a server
func (api *APIService) Start() {
	api.router.Run(api.addr)
}

// Router gives access to the router instance
func (api *APIService) Router() *gin.Engine {
	return api.router
}
