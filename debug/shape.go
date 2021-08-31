package debug

import (
	"context"
	"github.com/ProjectAthenaa/shape"
	protos "github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"google.golang.org/grpc"
	"log"
	"net"
)

func StartShapeServer() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	protos.RegisterShapeServer(server, Server{})

	log.Println("Started Shape Server on port 5000")
	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

type Server struct {
	protos.UnimplementedShapeServer
}

func (s Server) GenHeaders(ctx context.Context, site *protos.Site) (*protos.Headers, error) {
	return &protos.Headers{Values: shape.GenerateHeaders(site.Value)}, nil
}
