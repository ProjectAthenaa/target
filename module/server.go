package module

import (
	"context"
	"fmt"
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
	}

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

	fmt.Println("Card Number: ", data.Profile.Billing.Number)
	fmt.Println("Expiry Month / Expiry Year: ", data.Profile.Billing.ExpirationMonth, "/", data.Profile.Billing.ExpirationYear)
	fmt.Println("CVV: ", data.Profile.Billing.CVV)

	return &module.StartResponse{Started: true}, nil
}
