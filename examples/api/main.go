package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/auth"
	"github.com/txsvc/service/pkg/svc"
)

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
}

// testAPIResponse is the default way to respond to API requests
func testAPIResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdown()
		os.Exit(1)
	}()

	// create the service endpoint
	err := svc.New()
	if err != nil {
		os.Exit(1)
	}

	// create the JWT middleware
	a, err := auth.GetSecureJWTMiddleware("svcexample", "supersecretsecret")
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// add basic routes
	svc.AddDefaultEndpoints()
	svc.ServeStaticAssets("/", "./examples/api/public")

	// add custom endpoints with authentication
	api := svc.SecureGroup("/api", a.MiddlewareFunc())
	api.GET("/public", "chat.read", testAPIResponse)
	api.POST("/private", "chat.write", testAPIResponse)

	// add the service/router to a server on $PORT and launch it. This call BLOCKS !
	svc.Start()
}
