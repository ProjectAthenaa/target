package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"github.com/prometheus/common/log"
	"regexp"
	"strconv"
	"strings"
)

var (
	cartIdRe             = regexp.MustCompile(`"cart_id":"(\w+)"`)
	cartItemIdRe         = regexp.MustCompile(`"cart_item_id":"(\w+)"`)
	locationIdRe         = regexp.MustCompile(`"location_id":"(\d+)"`)
	guestIdRe            = regexp.MustCompile(`"targetGuid":"(\d+)"`)
	paymentInstructionRe = regexp.MustCompile(`"payment_instruction_id":"(\w+)"`)
)

func (tk *Task) ATC() {
	tk.SetStatus(module.STATUS_ADDING_TO_CART)
	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/cart_items?field_groups=CART%%2CCART_ITEMS%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(fmt.Sprintf(`{"cart_type":"REGULAR","channel_id":10,"shopping_context":"DIGITAL","cart_item":{"tcin":"%s","quantity":1,"item_channel_id":"10"}}`, tk.pid)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating atc request")
		tk.Stop()
		return
	}

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making atc request")
		tk.Stop()
		return
	}

	tk.cartid = cartIdRe.FindStringSubmatch(string(res.Body))[1]
	tk.cartitemid = cartItemIdRe.FindStringSubmatch(string(res.Body))[1]
	tk.SetStatus(module.STATUS_ADDED_TO_CART)
}

func (tk *Task) Login() {
	tk.SetStatus(module.STATUS_LOGGING_IN)
	req, err := tk.NewRequest("POST", "https://gsp.target.com/gsp/authentications/v1/credential_validations?client_id=ecom-web-1.0.0", []byte(fmt.Sprintf(`{"username":"%s","password":"%s","keep_me_signed_in":true,"device_info":{"user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","language":"en-US","color_depth":"24","device_memory":"8","pixel_ratio":"unknown","hardware_concurrency":"12","resolution":"[3148,886]","available_resolution":"[3098,886]","timezone_offset":"240","session_storage":"1","local_storage":"1","indexed_db":"1","add_behavior":"unknown","open_database":"1","cpu_class":"unknown","navigator_platform":"Win32","do_not_track":"unknown","regular_plugins":"[\"Chrome PDF Plugin::Portable Document Format::application/x-google-chrome-pdf~pdf\",\"Chrome PDF Viewer::::application/pdf~pdf\",\"Native Client::::application/x-nacl~,application/x-pnacl~\"]","adblock":"false","has_lied_languages":"false","has_lied_resolution":"false","has_lied_os":"false","has_lied_browser":"false","touch_support":"[0,false,false]","js_fonts":"[\"Arial\",\"Arial Black\",\"Arial Narrow\",\"Calibri\",\"Cambria\",\"Cambria Math\",\"Comic Sans MS\",\"Consolas\",\"Courier\",\"Courier New\",\"Georgia\",\"Helvetica\",\"Impact\",\"Lucida Console\",\"Lucida Sans Unicode\",\"Microsoft Sans Serif\",\"MS Gothic\",\"MS PGothic\",\"MS Sans Serif\",\"MS Serif\",\"Palatino Linotype\",\"Segoe Print\",\"Segoe Script\",\"Segoe UI\",\"Segoe UI Light\",\"Segoe UI Semibold\",\"Segoe UI Symbol\",\"Tahoma\",\"Times\",\"Times New Roman\",\"Trebuchet MS\",\"Verdana\",\"Wingdings\"]","navigator_vendor":"Google Inc.","navigator_app_name":"Netscape","navigator_app_code_name":"Mozilla","navigator_app_version":"5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","navigator_languages":"[\"en-US\"]","navigator_cookies_enabled":"true","navigator_java_enabled":"false","visitor_id":"%s","tealeaf_id":"%s","webgl_vendor":"Google Inc. (NVIDIA)~ANGLE (NVIDIA, NVIDIA GeForce GTX 1660 SUPER Direct3D11 vs_5_0 ps_5_0, D3D11-27.21.14.5671)","browser_name":"Chrome","browser_version":"92.0.4515.159","cpu_architecture":"amd64","device_vendor":"Unknown","device_model":"Unknown","device_type":"Unknown","engine_name":"Blink","engine_version":"92.0.4515.159","os_name":"Windows","os_version":"10"}}`, tk.Data.Metadata["username"], tk.Data.Metadata["password"], string(tk.FastClient.Jar.Peek("visitorId").Value()), string(tk.FastClient.Jar.Peek("TealeafAkaSid").Value()))))
	if err != nil {
		log.Error("create req: ", err)
		tk.SetStatus(module.STATUS_ERROR, "error creating login request")
		tk.Stop()
		return
	}

	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/login?client_id=ecom-web-1.0.0&ui_namespace=ui-default&back_button_action=browser&keep_me_signed_in=true&kmsi_default=false&actions=create_session_signin")

	tk.SetStatus(module.STATUS_GENERATING_COOKIES, "waiting for shape")
	headers, err := shapeClient.GenHeaders(tk.Ctx, &shape.Site{Value: shape.SITE_TARGET})
	if err != nil {
		log.Error("shape gen: ", err)
		tk.SetStatus(module.STATUS_ERROR, "error generating shape headers")
	}

	for k, v := range headers.Values {
		req.Headers[k] = []string{v}
	}

	req.Headers["Cookie"] = []string{`TealeafAkaSid=XcJfR4AGzMqu2Tis79OkjiIfyOMPZe4o; visitorId=017B9D0A2D0F020187D63FD82440F3D7; sapphire=1; UserLocation=00853|36.890|27.290|L|GR; criteo={}; ci_pixmgr=other; cd_user_id=17b9d0a47fc0-0979b60371ca31-c343365-1fa400-17b9d0a47fd1aa; adaptiveSessionId=A47545153; tlThirdPartyIds={"pt":"v2:29d8662f724798f235974866db75f3b9a5c0d6c1f196dd39b77c9253e63c9c46|3ae017d439d84697e45bf6ef263e53fd5a9cfed3915c9dd0dedb13ad7501f656"}; login-session=1AosAFM91egbloNszJgSuDMKuLsY8BjydetjpYA_G1MuCH75Bj7b_4WUXMZ2C24q; 3YCzT93n=A_WlvqV7AQAAcxp0ywxBMhyHaedCe2EA9p2NdgU_oZxSaKYWXRaMeJr58nv6AS6wMiSuctWowH8AAEB3AAAAAA|1|1|a6c3a9b353d0476f86de48fccf24708762a12e98; egsSessionId=243c1e00-313d-4f02-9895-0d1427553ade; fiatsCookie=DSI_|DSN_|DSZ_undefined; mid=10033639439; hasApp=false; accessToken=eyJraWQiOiJlYXMyIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiJjMDJjMTEzMC0xNGNkLTQ4ZDEtYTA2Ni01ZWMxOWRjMGMyNTQiLCJpc3MiOiJNSTYiLCJleHAiOjE2MzA2Njk0ODEsImlhdCI6MTYzMDU4MzA4MSwianRpIjoiVEdULjY2YmFmOTFjNzhiNjRjMzA4MmQ4MWRlMjJkNWY2ZjUzLWwiLCJza3kiOiJlYXMyIiwic3V0IjoiRyIsImRpZCI6IjEyYTI1MmRkZWU0ZjkzNmE0OWQ4MzZmZmU0ODgwNGJjMWYyOWEwNGM2OWRkMzllNTEyMWE4ZWE2ZTk3NjZhZjgiLCJzY28iOiJlY29tLm5vbmUsb3BlbmlkIiwiY2xpIjoiZWNvbS13ZWItMS4wLjAiLCJhc2wiOiJMIn0.UEGPKSY98f6HxH_Riq1C2t0a6tORZCPV_K67m1Ies_4ybzJJ_CciUqVlNnRXC9ovjswAuk6slah6BFxqY_OJActxUyhef6h8yIZjjUL1Yz2OX0aLo54uDsT0Ye9glG6CaCgKJrJfdgOyq4dmoelnFKpVOOlEQYeYKTsEqMinIESuLZ1Fjwh2o8pR9kx9wJa43KkFyMELKYaGzgsvB6HVhZdQrfMLpgsuViC_sBJdlQYUwGNn2LWfjUmC_Xybfy6edGlIIpwMLP7sraQdAwjF_OX7-vZxf640fUR7coE2C08a8kIjbJjN0fvmlfcFlf7ilTwUXlxcGqViqlEAP4WQsQ; idToken=eyJhbGciOiJub25lIn0.eyJzdWIiOiJjMDJjMTEzMC0xNGNkLTQ4ZDEtYTA2Ni01ZWMxOWRjMGMyNTQiLCJpc3MiOiJNSTYiLCJleHAiOjE2MzA2Njk0ODEsImlhdCI6MTYzMDU4MzA4MSwiYXNzIjoiTCIsInN1dCI6IkciLCJjbGkiOiJlY29tLXdlYi0xLjAuMCIsInBybyI6eyJmbiI6bnVsbCwiZW0iOm51bGwsInBoIjpmYWxzZSwibGVkIjpudWxsLCJsdHkiOmZhbHNlfX0.; refreshToken=IsVhle59tqDx9OxzxDuIV5DLM62cL4rVo4i0UYBOjF5vnUJQTKgDHqZEPzo0tzSpI1nfSj31FFPPKA0Doe4RjQ; guestType=G|1630583081000; targetMobileCookie=hasRC:false~guestLogonId:null~guestDisplayName:null~guestHasVerifiedPhone:false; mystate=1630583090909; ffsession={"sessionHash":"111c8b0351f01a1630573311063","sessionHit":35,"prevPageType":"/login","prevPageName":"Login: Login-Microsite","prevPageUrl":"https://www.target.com/login?client_id=ecom-web-1.0.0&ui_namespace=ui-default&back_button_action=browser&keep_me_signed_in=true&kmsi_default=false&actions=create_session_signin&username=agi***"}; JSESSIONID=7697E50FFEF96E1A3FAC3571ED1033B0`}

	res, err := tk.Do(req)
	if err != nil {
		log.Error("make req: ", err)
		tk.SetStatus(module.STATUS_ERROR, "error making login request")
		tk.Stop()
		return
	}

	log.Info(string(res.Body))
	tk.guestid = guestIdRe.FindStringSubmatch(string(res.Body))[1]
	tk.SetStatus(module.STATUS_LOGGED_IN)
}

func (tk *Task) RefreshCartId() {
	tk.SetStatus(module.STATUS_CHECKING_OUT, "refreshing cart")
	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/pre_checkout?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CDELIVERY_WINDOWS%%2CPAYMENT_INSTRUCTIONS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(`{"cart_type":"REGULAR"}`))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating cartid refresh request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-review?precheckout=true")

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making cartid refresh request")
		tk.Stop()
		return
	}

	tk.cartid = cartIdRe.FindStringSubmatch(string(res.Body))[1]
}

func (tk *Task) SubmitShipping() {
	tk.SetStatus(module.STATUS_SUBMITTING_SHIPPING)
	var form string
	if *tk.Data.Profile.Shipping.ShippingAddress.AddressLine2 != "" {
		form = fmt.Sprintf(`{"cart_type":"REGULAR","address":{"address_line1":"%s","address_line2":"%s","address_type":"SHIPPING","city":"%s","country":"%s","first_name":"%s","last_name":"%s","mobile":"%s","save_as_default":false,"state":"%s","zip_code":"%s"},"selected":true,"save_to_profile":true,"skip_verification":true}`, tk.Data.Profile.Shipping.ShippingAddress.AddressLine, tk.Data.Profile.Shipping.ShippingAddress.AddressLine2, tk.Data.Profile.Shipping.ShippingAddress.City, tk.Data.Profile.Shipping.ShippingAddress.Country, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.ZIP)
	} else {
		form = fmt.Sprintf(`{"cart_type":"REGULAR","address":{"address_line1":"%s","address_type":"SHIPPING","city":"%s","country":"%s","first_name":"%s","last_name":"%s","mobile":"%s","save_as_default":false,"state":"%s","zip_code":"%s"},"selected":true,"save_to_profile":true,"skip_verification":true}`, tk.Data.Profile.Shipping.ShippingAddress.AddressLine, tk.Data.Profile.Shipping.ShippingAddress.City, tk.Data.Profile.Shipping.ShippingAddress.Country, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.ZIP)
	}

	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/cart_shipping_addresses?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(form))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating shipping request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-delivery")

	_, err = tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making shipping request")
		tk.Stop()
		return
	}
}

func (tk *Task) SubmitPayment() {
	tk.SetStatus(module.STATUS_SUBMITTING_PAYMENT)
	var form string

	if *tk.Data.Profile.Shipping.BillingAddress.AddressLine2 != "" {
		form = fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"card_name":"%s","card_number":"%s","cvv":"%s","expiry_month":"%s","expiry_year":"%s"},"billing_address":{"address_line1":"%s","address_line2":"%s","city":"%s","first_name":"%s","last_name":"%s","phone":"%s","state":"%s","zip_code":"%s","country":"%s"}}`, tk.cartid, tk.Data.Profile.Shipping.FirstName+" "+tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Billing.Number, tk.Data.Profile.Billing.CVV, tk.Data.Profile.Billing.ExpirationMonth, "20"+tk.Data.Profile.Billing.ExpirationYear, tk.Data.Profile.Shipping.BillingAddress.AddressLine, *tk.Data.Profile.Shipping.BillingAddress.AddressLine2, tk.Data.Profile.Shipping.BillingAddress.City, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.BillingAddress.StateCode, tk.Data.Profile.Shipping.BillingAddress.ZIP, tk.Data.Profile.Shipping.BillingAddress.Country)
	} else {
		form = fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"card_name":"%s","card_number":"%s","cvv":"%s","expiry_month":"%s","expiry_year":"%s"},"billing_address":{"address_line1":"%s","city":"%s","first_name":"%s","last_name":"%s","phone":"%s","state":"%s","zip_code":"%s","country":"%s"}}`, tk.cartid, tk.Data.Profile.Shipping.FirstName+" "+tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Billing.Number, tk.Data.Profile.Billing.CVV, tk.Data.Profile.Billing.ExpirationMonth, "20"+tk.Data.Profile.Billing.ExpirationYear, tk.Data.Profile.Shipping.BillingAddress.AddressLine, tk.Data.Profile.Shipping.BillingAddress.City, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.BillingAddress.StateCode, tk.Data.Profile.Shipping.BillingAddress.ZIP, tk.Data.Profile.Shipping.BillingAddress.Country)
	}

	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/checkout_payments/v1/payment_instructions?key=%s", tk.apikey), []byte(form))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating payment request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-payment")

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making payment request")
		tk.Stop()
		return
	}

	var instructionresponse *PaymentInstructions
	json.Unmarshal(res.Body, &instructionresponse)

	tk.paymentinstructionid = instructionresponse.PaymentInstructionID
	tk.ReturningFields.Price = strconv.FormatFloat(instructionresponse.PaymentInstructionAmount, 'f', -1, 64)

}

func (tk *Task) CompareCard() {
	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/checkout_payments/v1/payment_instructions/%s?key=%s", tk.paymentinstructionid, tk.apikey), []byte(fmt.Sprintf(`{"cart_id":"%s","card_number":"%s"}`, tk.cartid, tk.Data.Profile.Billing.CVV)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating compare card request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-payment")

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making compare card request")
		tk.Stop()
		return
	}
	if res.StatusCode == 401 {
		tk.RefreshToken()
		tk.CompareCard()
		return
	}

	if strings.Contains(string(res.Body), "SUCCESS") {
		tk.SetStatus(module.STATUS_CHECKING_OUT, "card valid")
	} else {
		tk.SetStatus(module.STATUS_ERROR, "card not valid")
		tk.Stop()
		return
	}
}

func (tk *Task) SubmitCVV() {
	req, err := tk.NewRequest("PUT", fmt.Sprintf("https://carts.target.com/checkout_payments/v1/credit_card_compare?key=%s", tk.apikey), []byte(fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"cvv":"%s"}}`, tk.cartid, tk.Data.Profile.Billing.Number)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating compare card request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-payment")

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making cvv request")
		tk.Stop()
		return
	}

	if res.StatusCode == 401 {
		tk.RefreshToken()
		tk.SubmitCVV()
		return
	} else {
		tk.SetStatus(module.STATUS_CHECKING_OUT, "payment submitted")
	}
}

func (tk *Task) SubmitCheckout() {
	tk.SetStatus(module.STATUS_SUBMITTING_CHECKOUT)
	req, err := tk.NewRequest("PUT", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/checkout?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CDELIVERY_WINDOWS%%2CPAYMENT_INSTRUCTIONS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(`{"cart_type":"REGULAR","channel_id":10}`))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating compare card request")
		tk.Stop()
		return
	}
	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making checkout request")
		tk.Stop()
		return
	}

	var orderdata *CheckoutResponse
	json.Unmarshal(res.Body, &orderdata)

	tk.ReturningFields.Size = "os"
	tk.ReturningFields.ProductImage = tk.imagelink
	tk.ReturningFields.Color = "na"
	tk.ReturningFields.OrderNumber = orderdata.Orders[0].OrderID

	if res.StatusCode == 401 {
		tk.RefreshToken()
		tk.SubmitCheckout()
		return
	} else {
		tk.SetStatus(module.STATUS_CHECKED_OUT, "checked out")
	}
}
