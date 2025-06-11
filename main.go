//go:build !android || cli

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/freefly-systems/mavp2p-library/router"
)

func main() {
	// Create the router program with command line arguments
	p, err := router.NewProgram(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either the program to finish or a signal
	go func() {
		<-sigChan
		log.Println("Received shutdown signal, closing...")
		p.Close()
	}()

	// Wait for the program to finish
	p.Wait()
}
