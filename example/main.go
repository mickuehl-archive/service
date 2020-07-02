package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
	api, err := svc.New()
	if err != nil {
		os.Exit(1)
	}
	// add basic routes
	api.AddDefaultEndpoints()
	api.ServeStaticAssets("/", "./example/public")
	// add the router to a server on $PORT
	api.Start()
}
