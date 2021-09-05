package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/base"
	"github.com/ProjectAthenaa/sonic-core/sonic/face"
	"sync"
)

var _ face.ICallback = (*Task)(nil)

type Task struct {
	*base.BTask
	logincount           int
	pid                  string
	apikey               string
	cartid               string
	cartitemid           string
	storeid              string
	guestid              string
	paymentinstructionid string
	authcode             string
	imagelink            string
	submitCVV            bool
	submitAddress        bool
	sessionLock          *sync.Mutex
}

func NewTask(data *module.Data) *Task {
	task := &Task{BTask: &base.BTask{Data: data}}
	task.Callback = task
	task.Init()
	return task
}

func (tk *Task) OnInit() {
	tk.logincount = 0
	tk.sessionLock = &sync.Mutex{}
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
	//panic("")
	return
}

func (tk *Task) GetSession() {
	tk.sessionLock.Lock()
	defer tk.sessionLock.Unlock()
	funcArr := []func(){
		tk.APIKey,
		tk.InitData,
		tk.NearestStore,
		tk.OauthPost,
		tk.OauthSession,
		tk.Login,
		tk.CheckDetails,
		tk.OauthSession,
		tk.ClearCart,
		tk.OauthAuthCode,
	}

	for _, f := range funcArr {
		select {
		case <-tk.Ctx.Done():
			return
		default:
			f()
		}
	}

}

func (tk *Task) Flow() {
	go tk.GetSession()
	funcArr := []func(){
		tk.WaitForInstock,
		tk.sessionLock.Lock,
		tk.ATC,
		tk.SubmitShipping,
		tk.RefreshCartId,
		tk.SubmitCVV,
		tk.PaymentOauth,
		tk.SubmitCheckout,
	}

	for _, f := range funcArr {
		select {
		case <-tk.Ctx.Done():
			tk.sessionLock.Unlock()
			return
		default:
			f()
		}
	}
}
