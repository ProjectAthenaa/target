package main

import (
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	moduleServer "github.com/ProjectAthenaa/target/module"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"net"
)

func init() {
	target := &sonic.Module{
		Name: "Target",
		Fields: []sonic.InputField{
			{
				Validation: "https://www.target.*?",
				Label:      "URL",
			},
		},
	}

	if err := sonic.RegisterModule(target); err != nil {
		panic(err)
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
