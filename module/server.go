package module

import (
	"context"
	"github.com/ProjectAthenaa/sonic-core/protos/module"

)

type Server struct {
	module.UnimplementedModuleServer
}

func (s Server) Task(ctx context.Context, data *module.Data) (*module.StartResponse, error) {
	task := Task{}
	task.Init()
	if err := task.Start(data); err != nil {
		return nil, err
	}

	return &module.StartResponse{Started: true}, nil
}