package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/svc"
)

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
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
	// add basic routes
	svc.AddDefaultEndpoints()
	svc.ServeStaticAssets("/", "./example/public")
	// custom endpoints
	svc.GET("/api/public", TestAPIResponse)
	svc.GET("/api/private", TestAPIResponse)

	// add the router to a server on $PORT
	svc.Start()
}

// TestAPIResponse is the default way to respond to API requests
func TestAPIResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
