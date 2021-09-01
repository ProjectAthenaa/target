package module

import (
	"context"
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"github.com/prometheus/common/log"
	"os"
)

type Server struct {
	module.UnimplementedModuleServer
}

var shapeClient shape.ShapeClient

func init() {
	var err error
	if os.Getenv("DEBUG") == "1" {
		shapeClient, err = sonic.NewShapeClient("localhost:5000")
		if err != nil {
			panic(err)
		}

		fmt.Println(shapeClient.GenHeaders(context.Background(), &shape.Site{Value: shape.SITE_TARGET}))
		return
	}

	shapeClient, err = sonic.NewShapeClient()
	if err != nil {
		panic(err)
	}

}

func (s Server) Task(ctx context.Context, data *module.Data) (*module.StartResponse, error) {
	task := Task{}
	log.Info(data.TaskID)
	task.Init()
	if err := task.Start(data); err != nil {
		return nil, err
	}

	return &module.StartResponse{Started: true}, nil
}
