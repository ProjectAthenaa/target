package module

import (
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
			if sonic.ErrorContains(err, "Debug modes need custom address") {
				goto cloudClient
			}
			panic(err)
		}

		_, err = shapeClient.GenHeaders(context.Background(), nil)
		if sonic.ErrorContains(err, "Error while dialing dial tcp [::1]:3000") {
			goto cloudClient
		}

		return
	}

cloudClient:
	shapeClient, err = sonic.NewShapeClient()
	if err != nil {
		panic(err)
	}

}

func (s Server) Task(_ context.Context, data *module.Data) (*module.StartResponse, error) {
	//v, _ := json.Marshal(data)
	//fmt.Println(string(v))
	task := NewTask(data)
	if err := task.Start(data); err != nil {
		return nil, err
	}

	return &module.StartResponse{Started: true}, nil
}
