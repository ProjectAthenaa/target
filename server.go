package main

import (
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	moduleServer "github.com/ProjectAthenaa/target/module"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"net"
	"os"
	"strings"
)

func init() {
	var name = "target"

	if podName := os.Getenv("POD_NAME"); podName != "" {
		name = strings.Split(podName, "-")[0]
	}

	target := &sonic.Module{
		Name: name,
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
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalln("start listener: ", err)
	}

	server := grpc.NewServer()

	module.RegisterModuleServer(server, moduleServer.Server{})

	log.Info("Target Module Initialized")
	if err = server.Serve(listener); err != nil {
		log.Fatalln("start server: ", err)
	}
}
