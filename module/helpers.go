package module

import (
	"fmt"
	http "github.com/ProjectAthenaa/sonic-core/fasttls"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/json-iterator/go"
	"github.com/prometheus/common/log"
	"regexp"
)

var (
	apikeyRe = regexp.MustCompile(`"apiKey":"(\w+)"`)
	json     = jsoniter.ConfigFastest
)

type CheckoutResponse struct {
	Orders []struct {
		OrderID string `json:"order_id"`
	} `json:"orders"`
}

type PaymentInstructions struct {
	PaymentInstructionID     string  `json:"payment_instruction_id"`
	PaymentInstructionAmount float64 `json:"payment_instruction_amount"`
	RemainingBalance         float64 `json:"remaining_balance"`
}

type RefreshTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

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

func (tk *Task) RefreshToken() {
	req, err := tk.NewRequest("PUT", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/checkout?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CDELIVERY_WINDOWS%%2CPAYMENT_INSTRUCTIONS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(`{"grant_type":"refresh_token","client_credential":{"client_id":"ecom-web-1.0.0"},"device_info":{"user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","language":"en-US","color_depth":"24","device_memory":"8","pixel_ratio":"unknown","hardware_concurrency":"12","resolution":"[3148,886]","available_resolution":"[3098,886]","timezone_offset":"240","session_storage":"1","local_storage":"1","indexed_db":"1","add_behavior":"unknown","open_database":"1","cpu_class":"unknown","navigator_platform":"Win32","do_not_track":"unknown","regular_plugins":"[\"Chrome PDF Plugin::Portable Document Format::application/x-google-chrome-pdf~pdf\",\"Chrome PDF Viewer::::application/pdf~pdf\",\"Native Client::::application/x-nacl~,application/x-pnacl~\"]","adblock":"false","has_lied_languages":"false","has_lied_resolution":"false","has_lied_os":"false","has_lied_browser":"false","touch_support":"[0,false,false]","js_fonts":"[\"Arial\",\"Arial Black\",\"Arial Narrow\",\"Calibri\",\"Cambria\",\"Cambria Math\",\"Comic Sans MS\",\"Consolas\",\"Courier\",\"Courier New\",\"Georgia\",\"Helvetica\",\"Impact\",\"Lucida Console\",\"Lucida Sans Unicode\",\"Microsoft Sans Serif\",\"MS Gothic\",\"MS PGothic\",\"MS Sans Serif\",\"MS Serif\",\"Palatino Linotype\",\"Segoe Print\",\"Segoe Script\",\"Segoe UI\",\"Segoe UI Light\",\"Segoe UI Semibold\",\"Segoe UI Symbol\",\"Tahoma\",\"Times\",\"Times New Roman\",\"Trebuchet MS\",\"Verdana\",\"Wingdings\"]","navigator_vendor":"Google Inc.","navigator_app_name":"Netscape","navigator_app_code_name":"Mozilla","navigator_app_version":"5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","navigator_languages":"[\"en-US\"]","navigator_cookies_enabled":"true","navigator_java_enabled":"false","visitor_id":"017B6432D2940201872B1A2D05B771B8","tealeaf_id":"vxgLajdjLuk7vh_rMsnwW29e4rdbILHs","webgl_vendor":"Google Inc. (NVIDIA)~ANGLE (NVIDIA, NVIDIA GeForce GTX 1660 SUPER Direct3D11 vs_5_0 ps_5_0, D3D11-27.21.14.5671)","browser_name":"Chrome","browser_version":"92.0.4515.159","cpu_architecture":"amd64","device_vendor":"Unknown","device_model":"Unknown","device_type":"Unknown","engine_name":"Blink","engine_version":"92.0.4515.159","os_name":"Windows","os_version":"10"}}`))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating compare card request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://target.com")
	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making payment request")
		tk.Stop()
		return
	}

	var tokenresponse *RefreshTokenResp
	json.Unmarshal(res.Body, &tokenresponse)

	tk.FastClient.Jar.Set("accessToken", tokenresponse.AccessToken)
	tk.FastClient.Jar.Set("idToken", tokenresponse.IDToken)
	tk.FastClient.Jar.Set("refreshToken", tokenresponse.RefreshToken)
}

func (tk *Task) GenerateDefaultHeaders(referrer string) http.Headers {
	return http.Headers{
		`user-agent`:         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"},
		`accept`:             {`application/json`},
		`accept-encoding`:    {`gzip, deflate, br`},
		`accept-language`:    {`en-us`},
		`content-type`:       {`application/json`},
		`sec-ch-ua`:          {`"Chromium";v="91", " Not A;Brand";v="99", "Google Chrome";v="91"`},
		`x-application-name`: {`web`},
		`referer`:            {referrer},
	}
}

func (tk *Task) NearestStore() {
	req, err := tk.NewRequest("GET", fmt.Sprintf("https://api.target.com/shipt_deliveries/v1/stores?zip=%s&key=%s", tk.Data.Profile.Shipping.ShippingAddress.ZIP, tk.apikey), nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating find store request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com")

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making find store request")
		tk.Stop()
		return
	}
	if res.StatusCode == 401 {
		log.Info(string(res.Body))
		tk.RefreshToken()
		tk.NearestStore()
		return
	}

	tk.locationid = locationIdRe.FindStringSubmatch(string(res.Body))[1]
}
