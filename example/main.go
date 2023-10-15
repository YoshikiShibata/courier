package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YoshikiShibata/courier/example/internal/server"
)

func main() {
	grpcServer, err := server.NewGRPCServer(5000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "server.NewGRPCServer failed: %v", err)
		os.Exit(1)
	}

	go func() {
		log.Printf("Start Server")
		if err := grpcServer.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "grpcServer.Start failed: %v", err)
			os.Exit(1)
		}
	}()

	quitC := make(chan os.Signal, 1)
	signal.Notify(quitC, syscall.SIGINT, syscall.SIGTERM)
	<-quitC

	log.Printf("Stop Server")
	grpcServer.Stop()
}
