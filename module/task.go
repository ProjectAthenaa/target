package module

import (
	"fmt"
	log "github.com/ProjectAthenaa/sonic-core/logs"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/base"
	"github.com/ProjectAthenaa/sonic-core/sonic/face"
	"github.com/ProjectAthenaa/target/config"
	"sync"
)

var _ face.ICallback = (*Task)(nil)

type Task struct {
	*base.BTask
	logincount           int
	pid                  string
	apikey               string
	cartApiKey           string
	cartid               string
	cartitemid           string
	storeid              string
	guestid              string
	paymentinstructionid string
	authcode             string
	imagelink            string
	submitCVV            bool
	submitAddress        bool
	redirectcode         string
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
	tk.ReturningFields.Size = "OneSize"
	tk.ReturningFields.ProductName = "A-" + tk.Data.Metadata[*config.Module.Fields[0].FieldKey]
	tk.submitCVV = true
	tk.submitAddress = true
	tk.Flow()
}
func (tk *Task) OnPause() error {
	return nil
}
func (tk *Task) OnStopping() {
	tk.FastClient.Destroy()
	return
}

func (tk *Task) Flow() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("recovered: ", err)
			tk.SetStatus(module.STATUS_ERROR, "internal error")
			tk.Stop()
		}
	}()

	defer func() {
		tk.Stop()
	}()

	funcArr := []func(){
		tk.InitData,     //InitData and NearestStore have to be done before monitoring as they fill in critical variables like apikey and storeid
		tk.NearestStore, //add cache for nearest store?
		tk.OauthPost,
		tk.OauthSession,
		tk.AuthRedirect,
		tk.Login,
		tk.CartAPIKey,
		tk.AuthCode,
		tk.OauthAuthCode,
		tk.ClearCart,
		tk.CheckDetails,
		tk.OauthSession,
		tk.WaitForInstock, //monitoring
		//tk.sessionLock.Lock,
		tk.ATC,
		tk.RefreshCartId, //do we really need it?   //optimise get session
		tk.SubmitPayment,
		tk.SubmitShipping, //remove once better implementation is done, kiwi you did good job :)
		tk.SubmitCVV,      //remove once better implementation is done, this seems to be mandatory regardless if theres a payment or not
		tk.CompareCard,
		tk.SubmitCheckout,
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
