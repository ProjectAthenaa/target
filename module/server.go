package module

import (
	"context"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"os"
)

type Server struct {
	module.UnimplementedModuleServer
}

var shapeClient shape.ShapeClient

func init() {
	var err error
	if os.Getenv("DEBUG") == "1" {
		shapeClient, err = sonic.NewShapeClient("localhost:3000")
		if err != nil {
			panic(err)
		}
		return
	}

	shapeClient, err = sonic.NewShapeClient()
	if err != nil {
		panic(err)
	}
}

func (s Server) Task(_ context.Context, data *module.Data) (*module.StartResponse, error) {

	task := NewTask(data)
	if err := task.Start(data); err != nil {
		return nil, err
	}

	return &module.StartResponse{Started: true}, nil
}
