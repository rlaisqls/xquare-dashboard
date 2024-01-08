package main

import (
	"github.com/xquare-dashboard/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	service, _ := server.NewService()
	service.Start()
	service.Running()

	// Wait for the process to be shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
