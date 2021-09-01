package module

import (
	"fmt"
	http "github.com/ProjectAthenaa/sonic-core/fasttls"
	"github.com/ProjectAthenaa/sonic-core/fasttls/tls"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/base"
	"github.com/ProjectAthenaa/sonic-core/sonic/face"
	"github.com/ProjectAthenaa/sonic-core/sonic/frame"
	"github.com/ProjectAthenaa/threatmatrix"
	"strings"
)

var _ face.ICallback = (*Task)(nil)

type Task struct {
	*base.BTask
	Monitor              *frame.PubSub
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
	imageguestid		 string
	//temporary holders for metadata
	username string
	password string
}

func (tk *Task) OnInit() {
	proxy := tk.FormatProxy()
	tk.FastClient = http.NewClient(tls.HelloChrome_91, &proxy)
	tk.FastClient.Jar = http.NewJar()

	pubsub, err := frame.SubscribeToChannel(tk.Data.MonitorChannel)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, err.Error())
		tk.Stop()
		return
	}

	tk.Monitor = pubsub
}
func (tk *Task) OnPreStart() error {
	tk.FastClient.Jar.Set("UserLocation",fmt.Sprintf(`%s|||%s|%s`, tk.Data.Profile.Shipping.ShippingAddress.ZIP, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.Country),)
	tk.FastClient.Jar.Set("hasApp","false")
	return nil
}
func (tk *Task) OnStarting() {
	tk.Login()
	tk.SetStatus(module.STATUS_STARTING, "starting task")
	rawinfo := <-tk.Monitor.Channel
	monitorinfo := strings.Split(rawinfo.Payload, ":")
	tk.imageguestid = monitorinfo[1]
	tk.pid = monitorinfo[0]
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
