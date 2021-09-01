package main

import (
	"context"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/frame"
	"github.com/ProjectAthenaa/target/debug"
	module2 "github.com/ProjectAthenaa/target/module"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func init(){
	go debug.StartShapeServer()
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	module.RegisterModuleServer(server, module2.Server{})
	go func(){
		server.Serve(lis)
	}()
}

func TestModule(t *testing.T){
	subToken, controlToken := uuid.NewString(), uuid.NewString()

	productlink := "https://www.target.com/p/coleman-3pk-propane/-/A-81968997"

	tk := &module.Data{
		TaskID:         uuid.NewString(),
		Profile:        &module.Profile{
			Email:    "poprer656sad@gmail.com",
			Shipping: &module.Shipping{
				FirstName:         "Omar",
				LastName:          "Hu",
				PhoneNumber:       "6463222013",
				ShippingAddress:   &module.Address{
					AddressLine:  "7004 JFK BLVD E",
					AddressLine2: nil,
					Country:      "US",
					State:        "NEW JERSEY",
					City:         "WEST NEW YORK",
					ZIP:          "07093",
					StateCode:    "NJ",
				},
				BillingAddress:    &module.Address{
					AddressLine:  "7004 JFK BLVD E",
					AddressLine2: nil,
					Country:      "US",
					State:        "NEW JERSEY",
					City:         "WEST NEW YORK",
					ZIP:          "07093",
					StateCode:    "NJ",
				},
				BillingIsShipping: true,
			},
			Billing:  &module.Billing{
				Number:          "4207670236068972",
				ExpirationMonth: "05",
				ExpirationYear:  "25",
				CVV:             "997",
			},
		},
		Proxy:          &module.Proxy{
			IP:       "127.0.0.1",
			Port:     "8866",
		},
		TaskData:       &module.TaskData{
			RandomSize:  false,
			RandomColor: false,
			Color:       nil,
			Size:        nil,
			Link:        &productlink,
		},
		Metadata: map[string]string{
			"username":"chouerzi@gmail.com",
			"password":"Poprer656sad.",
		},
		Channels:       &module.Channels{
			UpdatesChannel:  subToken,
			CommandsChannel: controlToken,
		},
	}

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil{
		t.Fatal(err)
	}
	client := module.NewModuleClient(conn)

	pubsub, err := frame.SubscribeToChannel(subToken)
	if err != nil{
		t.Fatal(err)
	}

	_, err = client.Task(context.Background(), tk)
	if err != nil{
		t.Fatal(err)
	}

	for msg := range pubsub.Channel{
		log.Println(msg.Payload)
	}
}