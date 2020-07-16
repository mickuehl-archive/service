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
		router       *gin.Engine
		addr         string
		scopeMapping map[string]string // format: METHOD+PATH -> SCOPE
	}

	// SecureRouterGroup wraps a gin.RouterGroup and adds metadata for request authorization
	SecureRouterGroup struct {
		router *gin.RouterGroup
		path   string
	}
)

// service is a singleton
var service *APIService

func init() {
	// basic http stack config
	gin.DisableConsoleColor()
	// make sure the service is not initialized
	service = nil
}

// New creates a new service instance on addr. Only one instance per process is supported!
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
		router:       r,
		addr:         localAddr,
		scopeMapping: make(map[string]string),
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
	service.router.NoRoute(StandardNoRouteResponse)
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

// Group creates a new router group
func Group(relativePath string) *gin.RouterGroup {
	return service.router.Group(relativePath)
}

// SecureGroup creates a new router group with a security handler
func SecureGroup(relativePath string, secureHandler gin.HandlerFunc) *SecureRouterGroup {
	g := service.router.Group(relativePath)
	g.Use(secureHandler)
	return &SecureRouterGroup{
		router: g,
		path:   relativePath,
	}

}

// GetRequiredScopes returns the required scopes/scopes for this request or an empty string if none are required
func GetRequiredScopes(c *gin.Context) string {
	// FIXME just a naive implementation, optimizations etc
	return service.getRequiredScopes(c.Request.Method, c.FullPath())
}

func (s *APIService) registerSecureRoute(method, path, scope string) {
	// FIXME just a naive implementation, no safety net!
	s.scopeMapping[method+path] = scope
}

func (s *APIService) getRequiredScopes(method, path string) string {
	// FIXME just a naive implementation, no safety net!
	return s.scopeMapping[method+path]
}

//
// helper methods for 'standard' routes
//

// GET is a shortcut for router.Handle("GET", path, handle).
func GET(relativePath string, handler gin.HandlerFunc) {
	service.router.Handle(http.MethodGet, relativePath, handler)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func POST(relativePath string, handler gin.HandlerFunc) {
	service.router.Handle(http.MethodPost, relativePath, handler)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func DELETE(relativePath string, handler gin.HandlerFunc) {
	service.router.Handle(http.MethodDelete, relativePath, handler)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func PATCH(relativePath string, handler gin.HandlerFunc) {
	service.router.Handle(http.MethodPatch, relativePath, handler)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func PUT(relativePath string, handler gin.HandlerFunc) {
	service.router.Handle(http.MethodPut, relativePath, handler)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func OPTIONS(relativePath string, handler gin.HandlerFunc) {
	service.router.Handle(http.MethodOptions, relativePath, handler)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func HEAD(relativePath string, handler gin.HandlerFunc) {
	service.router.Handle(http.MethodHead, relativePath, handler)
}

//
// helper methods for secure routes
//

// GET is a shortcut for router.Handle("GET", path, handle).
func (srg *SecureRouterGroup) GET(relativePath, scope string, handler gin.HandlerFunc) {
	service.registerSecureRoute("GET", srg.path+relativePath, scope)
	srg.router.Handle(http.MethodGet, relativePath, handler)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (srg *SecureRouterGroup) POST(relativePath, scope string, handler gin.HandlerFunc) {
	service.registerSecureRoute("POST", srg.path+relativePath, scope)
	srg.router.Handle(http.MethodPost, relativePath, handler)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (srg *SecureRouterGroup) DELETE(relativePath, scope string, handler gin.HandlerFunc) {
	service.registerSecureRoute("DELETE", srg.path+relativePath, scope)
	srg.router.Handle(http.MethodDelete, relativePath, handler)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (srg *SecureRouterGroup) PATCH(relativePath, scope string, handler gin.HandlerFunc) {
	service.registerSecureRoute("PATCH", srg.path+relativePath, scope)
	srg.router.Handle(http.MethodPatch, relativePath, handler)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (srg *SecureRouterGroup) PUT(relativePath, scope string, handler gin.HandlerFunc) {
	service.registerSecureRoute("PUT", srg.path+relativePath, scope)
	srg.router.Handle(http.MethodPut, relativePath, handler)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (srg *SecureRouterGroup) OPTIONS(relativePath, scope string, handler gin.HandlerFunc) {
	service.registerSecureRoute("OPTIONS", srg.path+relativePath, scope)
	srg.router.Handle(http.MethodOptions, relativePath, handler)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (srg *SecureRouterGroup) HEAD(relativePath, scope string, handler gin.HandlerFunc) {
	service.registerSecureRoute("HEAD", srg.path+relativePath, scope)
	srg.router.Handle(http.MethodHead, relativePath, handler)
}
