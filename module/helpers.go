package module

import (
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"io/ioutil"
	"regexp"
)

var (
	apikeyRe = regexp.MustCompile(`"apiKey":"(\w+)"`)
)

func (tk *Task) APIKey(){
	res, err := tk.Client.Get("https://www.target.com/")
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "could not fetch site")
		tk.Stop()
		return
	}

	defer res.Body.Close()

	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "could not parse homepage body")
		tk.Stop()
		return
	}

	tk.apikey = apikeyRe.FindStringSubmatch(string(resbody))[1]
	for _, cookie := range res.Cookies(){
		switch cookie.Name{
			case "visitorId":
				tk.visitorid = cookie.Value
			case "TealeafAkaSid":
				tk.tealid = cookie.Value
		}
	}
}
