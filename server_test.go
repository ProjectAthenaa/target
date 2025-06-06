package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/core"
	"github.com/ProjectAthenaa/target/config"
	module2 "github.com/ProjectAthenaa/target/module"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
	"time"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func init() {
	//go debug.StartShapeServer()
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	module.RegisterModuleServer(server, module2.Server{})
	go func() {
		server.Serve(lis)
	}()
}

func TestModule(t *testing.T) {
	subToken, controlToken := uuid.NewString(), uuid.NewString()

	productID := "78808342"

	// 1moewci2:4k7cvljz:178.159.147.248:65112
	//username := "1moewci2"
	//password := "4k7cvljz"
	//ip := "178.159.147.248"
	//port := "65112"

	ip := "localhost"
	port := "8866"

	tk := &module.Data{
		TaskID: uuid.NewString(),
		Profile: &module.Profile{
			Email: "poprer656sad@gmail.com",
			Shipping: &module.Shipping{
				FirstName:   "Omar",
				LastName:    "Hu",
				PhoneNumber: "6463222013",
				ShippingAddress: &module.Address{
					AddressLine: "7004 JFK BLVD E",
					Country:     "US",
					State:       "NEW JERSEY",
					City:        "WEST NEW YORK",
					ZIP:         "07093",
					StateCode:   "NJ",
				},
				BillingAddress: &module.Address{
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
			Billing: &module.Billing{
				Number:          "4207670259298100",
				ExpirationMonth: "06",
				ExpirationYear:  "26",
				CVV:             "109",
			},
		},
		Proxy: &module.Proxy{
			//Username: &username,
			//Password: &password,
			IP:       ip,
			Port:     port,
		},
		TaskData: &module.TaskData{
			RandomSize:  false,
			RandomColor: false,
			Color:       []string{"1"},
			Size:        []string{"1"},
		},
		Metadata: map[string]string{
			"username":                        "dnwuiadiuwan@gmail.com",
			"password":                        "0o0p0o0P.",
			"UserID":                          "e99fa929-f1f2-4aad-b782-bfe6772fb2fc",
			*config.Module.Fields[0].FieldKey: productID,
		},
		Channels: &module.Channels{
			UpdatesChannel:  subToken,
			CommandsChannel: controlToken,
		},
	}

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	client := module.NewModuleClient(conn)
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))

	t.Log("connecting to redis")
	pubsub := core.Base.GetRedis("cache").Subscribe(ctx, fmt.Sprintf("tasks:updates:%s", subToken))
	t.Log("connected to redis")

	_, err = client.Task(context.Background(), tk)
	if err != nil {
		t.Fatal(err)
	}

	var start time.Time

	for msg := range pubsub.Channel() {
		var data module.Status
		_ = json.Unmarshal([]byte(msg.Payload), &data)
		fmt.Println(data.Status, data.Information["message"])

		if data.Status == module.STATUS_PRODUCT_FOUND {
			start = time.Now()
		}

		if data.Status == module.STATUS_CHECKED_OUT {
			fmt.Println(msg.Payload)
		}

		if data.Status == module.STATUS_STOPPED {
			fmt.Println(time.Since(start))
			return
		}
	}
}
