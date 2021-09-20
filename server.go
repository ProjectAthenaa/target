package main

import (
	log "github.com/ProjectAthenaa/sonic-core/logs"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/ProjectAthenaa/sonic-core/sonic/core"
	"github.com/ProjectAthenaa/target/config"
	moduleServer "github.com/ProjectAthenaa/target/module"
)

func init() {
	if err := log.Base().SetFormat("logger:stdout?json=true"); err != nil {
		log.Fatalln(err)
	}
	if err := sonic.RegisterModule(config.Module); err != nil {
		panic(err)
	}
}

func main() {
	core.ListenAndServe(config.Module.Name, &moduleServer.Server{})
}
