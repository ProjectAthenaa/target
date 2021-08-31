package main

import (
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/target/debug"
	moduleServer "github.com/ProjectAthenaa/target/module"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

func init() {
	if os.Getenv("DEBUG") == "1" {
		go debug.StartShapeServer()
	}
}

func main() {
	listener, err := net.Listen("tcp", "3000")
	if err != nil {
		log.Fatalln("start listener: ", err)
	}

	server := grpc.NewServer()

	module.RegisterModuleServer(server, moduleServer.Server{})

	if err = server.Serve(listener); err != nil {
		log.Fatalln("start server: ", err)
	}
}
