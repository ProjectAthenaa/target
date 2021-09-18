package main

import (
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/ProjectAthenaa/sonic-core/sonic/core"
	"github.com/ProjectAthenaa/target/config"
	moduleServer "github.com/ProjectAthenaa/target/module"
)

func init() {
	if err := sonic.RegisterModule(config.Module); err != nil {
		panic(err)
	}
}

func main() {
	core.RegisterModuleServer(config.Module.Name, &moduleServer.Server{})
}
