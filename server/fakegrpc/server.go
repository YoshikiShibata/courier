// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package fakegrpc

import (
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
)

// GRPCServer provides a gRPC server.
type GRPCServer struct {
	listener net.Listener
	server   *grpc.Server
}

// NewGRPCServer() assigns a new port and creates a gRPC server.
func NewGRPCServer() *GRPCServer {
	l, err := net.Listen("tcp4", ":0")
	if err != nil {
		log.Fatalf("net.Listen failed: %v", err)
	}
	server := grpc.NewServer()
	return &GRPCServer{
		listener: l,
		server:   server,
	}
}

// Addr returns the address with a port number.
func (gs *GRPCServer) Addr() string {
	return gs.listener.Addr().String()
}

// Port returns the port number of the address.
func (gs *GRPCServer) Port() string {
	addr := gs.Addr()
	return strings.Split(addr, ":")[1]
}

// Server returns a grpcServer.
func (gs *GRPCServer) Server() *grpc.Server {
	return gs.server
}

// Serve starts the gRPC server.
func (gs *GRPCServer) Serve() {
	err := gs.server.Serve(gs.listener)
	if err != nil {
		log.Fatalf("gs.server.Serve failed: %v", err)
	}
}
