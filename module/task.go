package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/base"
	"github.com/ProjectAthenaa/sonic-core/sonic/face"
	"github.com/ProjectAthenaa/threatmatrix"
	"github.com/prometheus/common/log"
)

var _ face.ICallback = (*Task)(nil)

type Task struct {
	*base.BTask
	pid                  string
	apikey               string
	cartid               string
	cartitemid           string
	storeid              string
	locationid           string
	guestid              string
	paymentinstructionid string
	imagelink            string
	username             string
	password             string
}

func NewTask(data *module.Data) *Task {
	task := &Task{BTask: &base.BTask{Data: data}}
	task.Callback = task
	task.Init()
	return task
}

func (tk *Task) OnInit() {
	return
}
func (tk *Task) OnPreStart() error {
	return nil
}
func (tk *Task) OnStarting() {
	tk.SetStatus(module.STATUS_STARTING, "starting task")
	tk.FastClient.CreateCookieJar()
	tk.FastClient.Jar.Set("UserLocation", fmt.Sprintf(`%s|||%s|%s`, tk.Data.Profile.Shipping.ShippingAddress.ZIP, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.Country))
	tk.FastClient.Jar.Set("hasApp", "false")
	tk.Flow()
}
func (tk *Task) OnPause() error {
	return nil
}
func (tk *Task) OnStopping() {
	tk.FastClient.Destroy()
	log.Info("stop called")
	panic("")
}

func (tk *Task) Flow() {
	tk.APIKey()
	tk.InitData()
	tk.NearestStore()
	//tk.Login()
	//tk.ATC()
	tk.NearestStore()
	//tk.RefreshCartId()
	threatmatrix.SendRequests(tk.cartid)
	//tk.SubmitShipping()
	//tk.SubmitPayment()
}
