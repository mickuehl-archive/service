package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/auth"
	"github.com/txsvc/service/pkg/svc"
)

func init() {
	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdown()
		os.Exit(1)
	}()
}

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
}

// testAPIResponse is the default way to respond to API requests
func testAPIResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	// used to secure cookies and sign the JWT token
	secret := env.GetString("SECRET", "supersecretsecret")

	// create the JWT middleware
	a, err := auth.GetSecureJWTMiddleware("svcexample", secret)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// add basic routes
	svc.AddDefaultEndpoints()
	svc.ServeStaticAssets("/", "./examples/app/public")

	// add custom endpoints with authentication
	api := svc.SecureGroup("/api", a.MiddlewareFunc())
	api.GET("/public", "chat.read", testAPIResponse)
	api.POST("/private", "chat.write", testAPIResponse)

	// add CORS handler, allowing all. See https://github.com/gin-contrib/cors
	svc.Use(cors.Default())

	// add session handler
	store := cookie.NewStore([]byte(secret))
	svc.Use(sessions.Sessions("svcexample", store))

	// add the service/router to a server on $PORT and launch it. This call BLOCKS !
	svc.Start()
}
