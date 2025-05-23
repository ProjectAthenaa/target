package module

import (
	"fmt"
	http "github.com/ProjectAthenaa/sonic-core/fasttls"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/target/config"
	"github.com/json-iterator/go"
	"regexp"
	"time"
)

var (
	referenceIdRe        = regexp.MustCompile(`"reference_id":"([\w-]+)"`)
	shapeSeedRe          = regexp.MustCompile(`init\("(.*?)"`)
	orderTotalRe         = regexp.MustCompile(`"total_authorization_amount":(\d+\.\d+)`)
	paymentInstructionRe = regexp.MustCompile(`"payment_instruction_id":"([\w-]+)"`)
	authCodeRe           = regexp.MustCompile(`"auth_code":"([\w-]+)"`)
	cartIdRe             = regexp.MustCompile(`"cart_id":"([\w-]+)"`)
	cartItemIdRe         = regexp.MustCompile(`"cart_item_id":"([\w-]+)"`)
	locationIdRe         = regexp.MustCompile(`"location_id":"(\d+)"`)
	guestIdRe            = regexp.MustCompile(`"targetGuid":"(\d+)"`)
	apikeyRe             = regexp.MustCompile(`"apiKey":"(\w+)"`)
	cartApiKeyRe         = regexp.MustCompile(`carts\.target\.com","apiKey":"(\w+)"`)
	loginErrRe           = regexp.MustCompile(`"errorKey":"(\w+)"`)
	totalCountRe         = regexp.MustCompile(`"total_count":(\d+)`)
	checkoutErrRe        = regexp.MustCompile(`"code":\s*"([\w-]+)"`)
	redirectCodeRe       = regexp.MustCompile(`code=([\w-]+)&`)
	json                 = jsoniter.ConfigFastest
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

func (tk *Task) GenerateDefaultHeaders(referrer string) http.Headers {
	return http.Headers{
		`user-agent`:         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"},
		`accept`:             {`application/json`},
		`accept-encoding`:    {`gzip, deflate, br`},
		`accept-language`:    {`en-us`},
		`content-type`:       {`application/json`},
		`sec-ch-ua`:          {`"Chromium";v="91", " Not A;Brand";v="99", "Google Chrome";v="91"`},
		`sec-ch-ua-mobile`:   {`?0`},
		`Sec-Fetch-Site`:     {`same-site`},
		`Sec-Fetch-Dest`:     {`empty`},
		`Sec-Fetch-Mode`:     {`cors`},
		`x-application-name`: {`web`},
		`referer`:            {referrer},
		`origin`:             {`https://www.target.com`},
		`Pragma`:             {`no-cache`},
		`Cache-Control`:      {`no-cache`},
		`Connection`:         {`keep-alive`},
	}
}

func (tk *Task) RefreshToken() {
	req, err := tk.NewRequest("PUT", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/checkout?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CDELIVERY_WINDOWS%%2CPAYMENT_INSTRUCTIONS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(`{"grant_type":"refresh_token","client_credential":{"client_id":"ecom-web-1.0.0"},"device_info":{"user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","language":"en-US","color_depth":"24","device_memory":"8","pixel_ratio":"unknown","hardware_concurrency":"12","resolution":"[3148,886]","available_resolution":"[3098,886]","timezone_offset":"240","session_storage":"1","local_storage":"1","indexed_db":"1","add_behavior":"unknown","open_database":"1","cpu_class":"unknown","navigator_platform":"Win32","do_not_track":"unknown","regular_plugins":"[\"Chrome PDF Plugin::Portable Document Format::application/x-google-chrome-pdf~pdf\",\"Chrome PDF Viewer::::application/pdf~pdf\",\"Native Client::::application/x-nacl~,application/x-pnacl~\"]","adblock":"false","has_lied_languages":"false","has_lied_resolution":"false","has_lied_os":"false","has_lied_browser":"false","touch_support":"[0,false,false]","js_fonts":"[\"Arial\",\"Arial Black\",\"Arial Narrow\",\"Calibri\",\"Cambria\",\"Cambria Math\",\"Comic Sans MS\",\"Consolas\",\"Courier\",\"Courier New\",\"Georgia\",\"Helvetica\",\"Impact\",\"Lucida Console\",\"Lucida Sans Unicode\",\"Microsoft Sans Serif\",\"MS Gothic\",\"MS PGothic\",\"MS Sans Serif\",\"MS Serif\",\"Palatino Linotype\",\"Segoe Print\",\"Segoe Script\",\"Segoe UI\",\"Segoe UI Light\",\"Segoe UI Semibold\",\"Segoe UI Symbol\",\"Tahoma\",\"Times\",\"Times New Roman\",\"Trebuchet MS\",\"Verdana\",\"Wingdings\"]","navigator_vendor":"Google Inc.","navigator_app_name":"Netscape","navigator_app_code_name":"Mozilla","navigator_app_version":"5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","navigator_languages":"[\"en-US\"]","navigator_cookies_enabled":"true","navigator_java_enabled":"false","visitor_id":"017B6432D2940201872B1A2D05B771B8","tealeaf_id":"vxgLajdjLuk7vh_rMsnwW29e4rdbILHs","webgl_vendor":"Google Inc. (NVIDIA)~ANGLE (NVIDIA, NVIDIA GeForce GTX 1660 SUPER Direct3D11 vs_5_0 ps_5_0, D3D11-27.21.14.5671)","browser_name":"Chrome","browser_version":"92.0.4515.159","cpu_architecture":"amd64","device_vendor":"Unknown","device_model":"Unknown","device_type":"Unknown","engine_name":"Blink","engine_version":"92.0.4515.159","os_name":"Windows","os_version":"10"}}`))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating refresh token request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://target.com")
	_, err = tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making refresh token request")
		tk.Stop()
		return
	}
}

func (tk *Task) OauthAuthCode() {
	req, err := tk.NewRequest("POST", "https://gsp.target.com/gsp/oauth_tokens/v2/client_tokens", []byte(fmt.Sprintf(`{"grant_type":"authorization_code","client_credential":{"client_id":"ecom-web-1.0.0"},"merge":"save","code":"%s","device_info":{"user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","language":"en-US","color_depth":"24","device_memory":"8","pixel_ratio":"unknown","hardware_concurrency":"12","resolution":"[3148,886]","available_resolution":"[3098,886]","timezone_offset":"240","session_storage":"1","local_storage":"1","indexed_db":"1","add_behavior":"unknown","open_database":"1","cpu_class":"unknown","navigator_platform":"Win32","do_not_track":"unknown","regular_plugins":"[\"Chrome PDF Plugin::Portable Document Format::application/x-google-chrome-pdf~pdf\",\"Chrome PDF Viewer::::application/pdf~pdf\",\"Native Client::::application/x-nacl~,application/x-pnacl~\"]","adblock":"false","has_lied_languages":"false","has_lied_resolution":"false","has_lied_os":"false","has_lied_browser":"false","touch_support":"[0,false,false]","js_fonts":"[\"Arial\",\"Arial Black\",\"Arial Narrow\",\"Calibri\",\"Cambria\",\"Cambria Math\",\"Comic Sans MS\",\"Consolas\",\"Courier\",\"Courier New\",\"Georgia\",\"Helvetica\",\"Impact\",\"Lucida Console\",\"Lucida Sans Unicode\",\"Microsoft Sans Serif\",\"MS Gothic\",\"MS PGothic\",\"MS Sans Serif\",\"MS Serif\",\"Palatino Linotype\",\"Segoe Print\",\"Segoe Script\",\"Segoe UI\",\"Segoe UI Light\",\"Segoe UI Semibold\",\"Segoe UI Symbol\",\"Tahoma\",\"Times\",\"Times New Roman\",\"Trebuchet MS\",\"Verdana\",\"Wingdings\"]","navigator_vendor":"Google Inc.","navigator_app_name":"Netscape","navigator_app_code_name":"Mozilla","navigator_app_version":"5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","navigator_languages":"[\"en-US\"]","navigator_cookies_enabled":"true","navigator_java_enabled":"false","visitor_id":"%s","tealeaf_id":"%s","webgl_vendor":"Google Inc. (NVIDIA)~ANGLE (NVIDIA, NVIDIA GeForce GTX 1660 SUPER Direct3D11 vs_5_0 ps_5_0, D3D11-27.21.14.5671)","browser_name":"Chrome","browser_version":"92.0.4515.159","cpu_architecture":"amd64","device_vendor":"Unknown","device_model":"Unknown","device_type":"Unknown","engine_name":"Blink","engine_version":"92.0.4515.159","os_name":"Windows","os_version":"10"}}`, tk.redirectcode, tk.FastClient.Jar.PeekValue("TealeafAkaSid"), tk.FastClient.Jar.PeekValue("visitorId"))))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not create auth code request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-login?shouldMergeCart=false&redirectToStep=PRECHECKOUT")

	_, err = tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not post auth code request")
		tk.Stop()
		return
	}
}

func (tk *Task) ClearCart() {
	req, err := tk.NewRequest("PUT", fmt.Sprintf(`https://carts.target.com/web_checkouts/v1/cart?field_groups=ADDRESSES%%2CCART_ITEMS%%2CCART%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s`, tk.cartApiKey),[]byte(fmt.Sprintf(`{"cart_type":"REGULAR","channel_id":"10","shopping_context":"DIGITAL","guest_location":{"state":"%s","latitude":"","zip_code":"%s","longitude":"","country":"US"},"shopping_location_id":"%s"}`, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.ZIP, tk.storeid)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating tax request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf("https://www.target.com/p/-/A-%s", tk.Data.Metadata[*config.Module.Fields[0].FieldKey]))

	//	cookiejar.ReleaseCookieJar(tk.FastClient.Jar)
	//	tk.FastClient.Jar = nil
	//	req.Headers["Cookie"] = []string{`TealeafAkaSid=4yhXI0qbi4RDOUJwHFSlW1RQHfdkGYGs;`+
	//`visitorId=017BFD3904FF02019F3B4429BEEBBD59;`+
	//`login-session=7pWxn9BICBjEOQPEpSrZyoAec6qop4lvEJOGcXu39Wd_kdRH_gTgh5hYDk35cNas;`+
	//`accessToken=eyJraWQiOiJlYXMyIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIyMDA3Njg0NTM1NyIsImlzcyI6Ik1JNiIsImV4cCI6MTYzMjA1NTM5NCwiaWF0IjoxNjMyMDQwOTk0LCJqdGkiOiJUR1QuYWVjNGJjYWUzMDMxNDg2YThiMjIzMWIzMjI5M2UzYTMtbSIsInNreSI6ImVhczIiLCJzdXQiOiJSIiwiZGlkIjoiNzQwOTU1Mjc0ZTQyMjQyODBiYTYyM2IyODQzNGM3ZTBmN2UzOTEwZGUzNWI2YThkMGY1NjMzM2Y5ZDIzOWYzNSIsImVpZCI6InRlcnJ5ZGF2aXM5MDNAZ21haWwuY29tIiwiZ3NzIjoxLjAsInNjbyI6ImVjb20ubWVkLG9wZW5pZCIsImNsaSI6ImVjb20td2ViLTEuMC4wIiwiYXNsIjoiTSJ9.pioTTg9Ret_vMb8vnmt2SwlX03i6_4KY0XUL5n408Zvf3PSmS7teHk14tGN0tFbA9IjOqJk1uwE2XXyEzh_N471cCEKQ8m91wZ7VpRMjUIhyrzqWKU4zFgxHeoSE8kr5pQ0TCoMuImMWuVJHwvAhk0YkGGU0ZNSpnzzNIROIXf0GiJntOTq2ASD8Jg2tGaT6ra9iPoo_THzYeJKkr7m3hCwf0VrOnv5kjb504BQmx0MysejH3pIrTdwFFkB6gOW5oHHL2deE9bHoBMFmjC7dtQhnY24XPniQJ9z2Y-gFO4W00FiF31rwjnDuWJtvYibJ69Bglu0rIkS3OC6JQU66fQ;`+
	//`idToken=eyJhbGciOiJub25lIn0.eyJzdWIiOiIyMDA3Njg0NTM1NyIsImlzcyI6Ik1JNiIsImV4cCI6MTYzMjA1NTM5NCwiaWF0IjoxNjMyMDQwOTk0LCJhc3MiOiJNIiwic3V0IjoiUiIsImNsaSI6ImVjb20td2ViLTEuMC4wIiwicHJvIjp7ImZuIjoidGVycnkiLCJlbSI6InRlcnJ5ZGF2aXM5MDNAZ21haWwuY29tIiwicGgiOmZhbHNlLCJsZWQiOm51bGwsImx0eSI6ZmFsc2V9fQ.;`,
	//}

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not post clear cart request")
		tk.Stop()
		return
	}

	if res.StatusCode == 403 {
		tk.SetStatus(module.STATUS_ERROR, "could not clear cart, likely proxy error")
		return
	}

	for _, sm := range cartItemIdRe.FindAllSubmatch(res.Body, -1) {
		sm := sm
		go func() {
			delreq, delerr := tk.NewRequest("DELETE", fmt.Sprintf(`https://carts.target.com/web_checkouts/v1/cart_items/%s?cart_type=REGULAR&field_groups=CART%%2CCART_ITEMS%%2CSUMMARY%%2CPROMOTION_CODES%%2CADDRESSES&key=%s`, string(sm[1]), tk.cartApiKey), nil)
			if delerr != nil {
				tk.SetStatus(module.STATUS_ERROR, "could not create clear cart request")
				tk.Stop()
				return
			}
			delreq.Headers = tk.GenerateDefaultHeaders("https://www.target.com")

			delres, delerr := tk.Do(delreq)
			if delerr != nil {
				tk.SetStatus(module.STATUS_ERROR, "could not post clear cart request")
				tk.Stop()
				return
			}

			if delres.StatusCode != 200 {
				tk.SetStatus(module.STATUS_ERROR, "could not delete item "+string(sm[1]))
				tk.Stop()
				return
			}
		}()
	}
}

func (tk *Task) CheckDetails() {
	go func() {
		req, err := tk.NewRequest("GET", fmt.Sprintf("https://profile.target.com/WalletWEB/wallet/v5/tenders?key=%s&savings=true", tk.apikey), nil)
		if err != nil {
			tk.SetStatus(module.STATUS_ERROR, "could not check payment methods")
			tk.Stop()
			return
		}

		req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/account/payments")

		res, err := tk.Do(req)
		if err != nil {
			tk.SetStatus(module.STATUS_ERROR, "could not check payment methods")
			tk.Stop()
			return
		}

		if match := totalCountRe.FindStringSubmatch(string(res.Body)); len(match) == 2 && match[1] >= "1" {
			tk.submitCVV = false
		}
	}()

	req, err := tk.NewRequest("GET", fmt.Sprintf("https://api.target.com/guest_addresses/v1/addresses?key=%s", tk.apikey), nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not check addresses")
		tk.Stop()
		return
	}

	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/account/addresses")

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not check addresses")
		tk.Stop()
	}

	if match := totalCountRe.FindStringSubmatch(string(res.Body)); len(match) == 2 && match[1] >= "1" {
		tk.submitAddress = false
	}
}

func (tk *Task) AuthRedirect() {
	req, err := tk.NewRequest("GET", fmt.Sprintf(`https://gsp.target.com/gsp/authentications/v1/auth_codes?client_id=ecom-web-1.0.0&state=%d&redirect_uri=https%%3A%%2F%%2Fwww.target.com%%2F&assurance_level=M`, time.Now().Unix()), nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not get auth code")
		tk.Stop()
		return
	}

	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/")

	_, err = tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not get auth code")
		tk.Stop()
		return
	}

}

func (tk *Task) AuthCode() {
	req, err := tk.NewRequest("GET", `https://gsp.target.com/gsp/authentications/v1/auth_codes?client_id=ecom-web-1.0.0`, nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not get auth code")
		tk.Stop()
		return
	}

	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/")

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "could not get auth code")
		tk.Stop()
	}

	if res.StatusCode == 302 {
		tk.redirectcode = redirectCodeRe.FindStringSubmatch(res.Headers["Location"][0])[1]
	}
}
