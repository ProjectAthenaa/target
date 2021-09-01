package module

import (
	"fmt"
	http "github.com/ProjectAthenaa/sonic-core/fasttls"
	"github.com/ProjectAthenaa/sonic-core/fasttls/tls"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"github.com/ProjectAthenaa/sonic-core/sonic/base"
	"github.com/ProjectAthenaa/sonic-core/sonic/core"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent"
	"github.com/ProjectAthenaa/sonic-core/sonic/database/ent/task"
	"github.com/ProjectAthenaa/sonic-core/sonic/face"
	"github.com/ProjectAthenaa/threatmatrix"
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
	tealid               string
	visitorid            string
	guestid              string
	paymentinstructionid string
	imagelink         	string
	username             string
	password             string

	product *ent.Product
}

func (tk *Task) OnInit() {
	dbtask, err := core.Base.GetPg("pg").
		Task.
		Query().
		WithProduct().
		Where(
			task.ID(
				sonic.UUIDParser(tk.Data.TaskID),
			),
		).
		First(tk.Ctx)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, sonic.EntErr(err))
		tk.Stop()
		return
	}
	tk.product = dbtask.Edges.Product[0]
	tk.FastClient = http.NewClient(tls.HelloChrome_91, tk.FormatProxy())
	tk.FastClient.Jar = http.NewJar()
}
func (tk *Task) OnPreStart() error {
	tk.FastClient.Jar.Set("UserLocation", fmt.Sprintf(`%s|||%s|%s`, tk.Data.Profile.Shipping.ShippingAddress.ZIP, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.Country))
	tk.FastClient.Jar.Set("hasApp", "false")
	tk.InitData()
	tk.NearestStore()
	return nil
}
func (tk *Task) OnStarting() {
	tk.Login()
	tk.SetStatus(module.STATUS_STARTING, "starting task")
	//targetinfo := <-tk.Monitor.Chan(tk.Ctx)
	//tk.imageguestid = targetinfo["guestid"].(string)
	//tk.pid = targetinfo["pid"].(string)
	tk.ATC()
	tk.NearestStore()
	tk.RefreshCartId()
	threatmatrix.SendRequests(tk.cartid)
	tk.SubmitShipping()
	tk.SubmitPayment()
}
func (tk *Task) OnPause() error {
	return nil
}
func (tk *Task) OnStopping() {

}
