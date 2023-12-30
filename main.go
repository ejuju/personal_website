package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ejuju/my-website/internal/service"
)

func main() {
	service, err := service.New()
	if err != nil {
		panic(err)
	}
	go service.Run()

	// Wait for termination.
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm
	err = service.Shutdown()
	if err != nil {
		panic(err)
	}
}
