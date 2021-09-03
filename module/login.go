package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"strings"
	"time"
)

func (tk *Task) APIKey() {
	req, err := tk.NewRequest("GET", "https://www.target.com/", nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not make homepage req")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://target.com")

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not fetch site")
		tk.Stop()
		return
	}

	tk.apikey = apikeyRe.FindStringSubmatch(string(res.Body))[1]
}

func (tk *Task) OauthPost() {
	req, err := tk.NewRequest("POST", "https://gsp.target.com/gsp/oauth_tokens/v2/client_tokens", []byte(fmt.Sprintf(`{"grant_type":"refresh_token","client_credential":{"client_id":"ecom-web-1.0.0"},"device_info":{"user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","language":"en-US","color_depth":"24","device_memory":"8","pixel_ratio":"unknown","hardware_concurrency":"12","resolution":"[3148,886]","available_resolution":"[3098,886]","timezone_offset":"240","session_storage":"1","local_storage":"1","indexed_db":"1","add_behavior":"unknown","open_database":"1","cpu_class":"unknown","navigator_platform":"Win32","do_not_track":"unknown","regular_plugins":"[\"Chrome PDF Plugin::Portable Document Format::application/x-google-chrome-pdf~pdf\",\"Chrome PDF Viewer::::application/pdf~pdf\",\"Native Client::::application/x-nacl~,application/x-pnacl~\"]","adblock":"false","has_lied_languages":"false","has_lied_resolution":"false","has_lied_os":"false","has_lied_browser":"false","touch_support":"[0,false,false]","js_fonts":"[\"Arial\",\"Arial Black\",\"Arial Narrow\",\"Calibri\",\"Cambria\",\"Cambria Math\",\"Comic Sans MS\",\"Consolas\",\"Courier\",\"Courier New\",\"Georgia\",\"Helvetica\",\"Impact\",\"Lucida Console\",\"Lucida Sans Unicode\",\"Microsoft Sans Serif\",\"MS Gothic\",\"MS PGothic\",\"MS Sans Serif\",\"MS Serif\",\"Palatino Linotype\",\"Segoe Print\",\"Segoe Script\",\"Segoe UI\",\"Segoe UI Light\",\"Segoe UI Semibold\",\"Segoe UI Symbol\",\"Tahoma\",\"Times\",\"Times New Roman\",\"Trebuchet MS\",\"Verdana\",\"Wingdings\"]","navigator_vendor":"Google Inc.","navigator_app_name":"Netscape","navigator_app_code_name":"Mozilla","navigator_app_version":"5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","navigator_languages":"[\"en-US\"]","navigator_cookies_enabled":"true","navigator_java_enabled":"false","visitor_id":"%s","tealeaf_id":"%s","webgl_vendor":"Google Inc. (NVIDIA)~ANGLE (NVIDIA, NVIDIA GeForce GTX 1660 SUPER Direct3D11 vs_5_0 ps_5_0, D3D11-27.21.14.5671)","browser_name":"Chrome","browser_version":"92.0.4515.159","cpu_architecture":"amd64","device_vendor":"Unknown","device_model":"Unknown","device_type":"Unknown","engine_name":"Blink","engine_version":"92.0.4515.159","os_name":"Windows","os_version":"10"}}`, tk.FastClient.Jar.Peek("visitorId"), tk.FastClient.Jar.Peek("TealeafAkaSid"))))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not create oauth request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com")

	_, err = tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not fetch oauth request")
		tk.Stop()
		return
	}
}

func (tk *Task) OauthSession() {
	req, err := tk.NewRequest("GET", fmt.Sprintf(`https://gsp.target.com/gsp/authentications/v1/spa_auth_codes?client_id=ecom-web-1.0.0&state=%d&keep_me_signed_in=false`, time.Now().Unix()), nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not create oauth2 second request")
		tk.Stop()
		return
	}

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not get oauth2 second request")
		tk.Stop()
		return
	}

	if v := authCodeRe.FindSubmatch(res.Body); len(v) == 2 {
		tk.authcode = string(v[1])
	}

}

func (tk *Task) Login() {
	tk.SetStatus(module.STATUS_LOGGING_IN)
	req, err := tk.NewRequest("POST", "https://gsp.target.com/gsp/authentications/v1/credential_validations?client_id=ecom-web-1.0.0", []byte(fmt.Sprintf(`{"username":"%s","password":"%s","keep_me_signed_in":true,"device_info":{"user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","language":"en-US","color_depth":"24","device_memory":"8","pixel_ratio":"unknown","hardware_concurrency":"12","resolution":"[3148,886]","available_resolution":"[3098,886]","timezone_offset":"240","session_storage":"1","local_storage":"1","indexed_db":"1","add_behavior":"unknown","open_database":"1","cpu_class":"unknown","navigator_platform":"Win32","do_not_track":"unknown","regular_plugins":"[\"Chrome PDF Plugin::Portable Document Format::application/x-google-chrome-pdf~pdf\",\"Chrome PDF Viewer::::application/pdf~pdf\",\"Native Client::::application/x-nacl~,application/x-pnacl~\"]","adblock":"false","has_lied_languages":"false","has_lied_resolution":"false","has_lied_os":"false","has_lied_browser":"false","touch_support":"[0,false,false]","js_fonts":"[\"Arial\",\"Arial Black\",\"Arial Narrow\",\"Calibri\",\"Cambria\",\"Cambria Math\",\"Comic Sans MS\",\"Consolas\",\"Courier\",\"Courier New\",\"Georgia\",\"Helvetica\",\"Impact\",\"Lucida Console\",\"Lucida Sans Unicode\",\"Microsoft Sans Serif\",\"MS Gothic\",\"MS PGothic\",\"MS Sans Serif\",\"MS Serif\",\"Palatino Linotype\",\"Segoe Print\",\"Segoe Script\",\"Segoe UI\",\"Segoe UI Light\",\"Segoe UI Semibold\",\"Segoe UI Symbol\",\"Tahoma\",\"Times\",\"Times New Roman\",\"Trebuchet MS\",\"Verdana\",\"Wingdings\"]","navigator_vendor":"Google Inc.","navigator_app_name":"Netscape","navigator_app_code_name":"Mozilla","navigator_app_version":"5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","navigator_languages":"[\"en-US\"]","navigator_cookies_enabled":"true","navigator_java_enabled":"false","visitor_id":"%s","tealeaf_id":"%s","webgl_vendor":"Google Inc. (NVIDIA)~ANGLE (NVIDIA, NVIDIA GeForce GTX 1660 SUPER Direct3D11 vs_5_0 ps_5_0, D3D11-27.21.14.5671)","browser_name":"Chrome","browser_version":"92.0.4515.159","cpu_architecture":"amd64","device_vendor":"Unknown","device_model":"Unknown","device_type":"Unknown","engine_name":"Blink","engine_version":"92.0.4515.159","os_name":"Windows","os_version":"10"}}`, tk.Data.Metadata["username"], tk.Data.Metadata["password"], string(tk.FastClient.Jar.Peek("visitorId").Value()), string(tk.FastClient.Jar.Peek("TealeafAkaSid").Value()))))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating login request")
		tk.Stop()
		return
	}

	tk.FastClient.Jar.Set("sapphire", "1")
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/login?client_id=ecom-web-1.0.0&ui_namespace=ui-default&back_button_action=browser&keep_me_signed_in=true&kmsi_default=false&actions=create_session_signin")

	tk.SetStatus(module.STATUS_GENERATING_COOKIES, "waiting for shape")
	headers, err := shapeClient.GenHeaders(tk.Ctx, &shape.Site{Value: shape.SITE_TARGET})
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error generating shape headers")
	}

	for k, v := range headers.Values {
		req.Headers[k] = []string{v}
	}
	req.Headers["Accept"] = []string{"application/json"}

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making login request")
		tk.Stop()
		return
	}

	if strings.Contains(string(res.Body), "_ERR_AUTH_DENIED") {
		data, _ := json.Marshal(headers.Values)
		fmt.Println(string(data))
		fmt.Println(string(res.Body))
		tk.SetStatus(module.STATUS_ERROR)
		tk.Stop()
		return
	}

	//if strings.Contains(string(res.Body), "_ERR_AUTH_DENIED") && tk.logincount < 3 {
	//	tk.logincount++
	//	tk.Login()
	//	return
	//} else if tk.logincount >= 3 {
	//	tk.SetStatus(module.STATUS_ERROR, "login error")
	//	tk.Stop()
	//	return
	//}

	tk.guestid = guestIdRe.FindStringSubmatch(string(res.Body))[1]
	tk.SetStatus(module.STATUS_LOGGED_IN)
}
