package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
)

var build = "develop"

// setThreadsForService sets the correct number of threads for service based on what is available
// either by the machine or quotas
func setThreadsForService() error {
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	return nil
}

func main() {
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("maxprocs: %w", err)
		os.Exit(1)
	}
	err := setThreadsForService()
	if err != nil {
		fmt.Println("unable to set threads for service: %w", err)
		os.Exit(1)
	}
	g := runtime.GOMAXPROCS(0)
	log.Printf("starting service build[%s] CPU[%d]", build, g)
	defer log.Println("service ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	log.Println("stopping service")
}
