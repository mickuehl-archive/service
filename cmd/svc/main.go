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

	api, err := svc.New()
	if err != nil {
		os.Exit(1)
	}
	api.AddDefaultEndpoints()
	api.ServeStaticAssets("/", "public")
	api.Start()
}
