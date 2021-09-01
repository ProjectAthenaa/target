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
	req, err := tk.FastClient.NewRequest("POST", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/cart_items?field_groups=CART%%2CCART_ITEMS%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(fmt.Sprintf(`{"cart_type":"REGULAR","channel_id":10,"shopping_context":"DIGITAL","cart_item":{"tcin":"%s","quantity":1,"item_channel_id":"10"}}`, tk.pid)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating atc request")
		tk.Stop()
		return
	}

	res, err := tk.FastClient.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making atc request")
		tk.Stop()
		return
	}

	tk.cartid = cartIdRe.FindStringSubmatch(string(res.Body))[1]
	tk.cartitemid = cartItemIdRe.FindStringSubmatch(string(res.Body))[1]
}

func (tk *Task) NearestStore() {
	storereq, err := tk.FastClient.NewRequest("GET", fmt.Sprintf("https://api.target.com/shipt_deliveries/v1/stores?zip=%s&key=%s", tk.Data.Profile.Shipping.ShippingAddress.ZIP, tk.apikey), nil)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating find store request")
		tk.Stop()
		return
	}
	res, err := tk.FastClient.Do(storereq)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making find store request")
		tk.Stop()
		return
	}
	if res.StatusCode == 401 {
		tk.RefreshToken()
		tk.NearestStore()
	}

	tk.locationid = locationIdRe.FindStringSubmatch(string(res.Body))[1]
}

func (tk *Task) Login() {
	req, err := tk.FastClient.NewRequest("POST", "https://gsp.target.com/gsp/authentications/v1/credential_validations?client_id=ecom-web-1.0.0", []byte(fmt.Sprintf(`{"username":"%s","password":"%s","keep_me_signed_in":true,"device_info":{"user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","language":"en-US","color_depth":"24","device_memory":"8","pixel_ratio":"unknown","hardware_concurrency":"12","resolution":"[3148,886]","available_resolution":"[3098,886]","timezone_offset":"240","session_storage":"1","local_storage":"1","indexed_db":"1","add_behavior":"unknown","open_database":"1","cpu_class":"unknown","navigator_platform":"Win32","do_not_track":"unknown","regular_plugins":"[\"Chrome PDF Plugin::Portable Document Format::application/x-google-chrome-pdf~pdf\",\"Chrome PDF Viewer::::application/pdf~pdf\",\"Native Client::::application/x-nacl~,application/x-pnacl~\"]","adblock":"false","has_lied_languages":"false","has_lied_resolution":"false","has_lied_os":"false","has_lied_browser":"false","touch_support":"[0,false,false]","js_fonts":"[\"Arial\",\"Arial Black\",\"Arial Narrow\",\"Calibri\",\"Cambria\",\"Cambria Math\",\"Comic Sans MS\",\"Consolas\",\"Courier\",\"Courier New\",\"Georgia\",\"Helvetica\",\"Impact\",\"Lucida Console\",\"Lucida Sans Unicode\",\"Microsoft Sans Serif\",\"MS Gothic\",\"MS PGothic\",\"MS Sans Serif\",\"MS Serif\",\"Palatino Linotype\",\"Segoe Print\",\"Segoe Script\",\"Segoe UI\",\"Segoe UI Light\",\"Segoe UI Semibold\",\"Segoe UI Symbol\",\"Tahoma\",\"Times\",\"Times New Roman\",\"Trebuchet MS\",\"Verdana\",\"Wingdings\"]","navigator_vendor":"Google Inc.","navigator_app_name":"Netscape","navigator_app_code_name":"Mozilla","navigator_app_version":"5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36","navigator_languages":"[\"en-US\"]","navigator_cookies_enabled":"true","navigator_java_enabled":"false","visitor_id":"%s","tealeaf_id":"%s","webgl_vendor":"Google Inc. (NVIDIA)~ANGLE (NVIDIA, NVIDIA GeForce GTX 1660 SUPER Direct3D11 vs_5_0 ps_5_0, D3D11-27.21.14.5671)","browser_name":"Chrome","browser_version":"92.0.4515.159","cpu_architecture":"amd64","device_vendor":"Unknown","device_model":"Unknown","device_type":"Unknown","engine_name":"Blink","engine_version":"92.0.4515.159","os_name":"Windows","os_version":"10"}}`, tk.Data.Metadata["username"], tk.Data.Metadata["password"], tk.visitorid, tk.tealid)))
	if err != nil {
		log.Error("create req: ", err)
		tk.SetStatus(module.STATUS_ERROR, "error creating login request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-login")

	headers, err := shapeClient.GenHeaders(tk.Ctx, &shape.Site{Value: shape.SITE_TARGET})
	if err != nil {
		log.Error("shape gen: ", err)
		tk.SetStatus(module.STATUS_ERROR, "error generating shape headers")
	}

	for k, v := range headers.Values {
		req.Headers[k] = []string{v}
	}

	res, err := tk.FastClient.Do(req)
	if err != nil {
		log.Error("make req: ", err)
		tk.SetStatus(module.STATUS_ERROR, "error making login request")
		tk.Stop()
		return
	}

	tk.guestid = guestIdRe.FindStringSubmatch(string(res.Body))[1]

}

func (tk *Task) RefreshCartId() {
	req, err := tk.FastClient.NewRequest("POST", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/pre_checkout?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CDELIVERY_WINDOWS%%2CPAYMENT_INSTRUCTIONS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(`{"cart_type":"REGULAR"}`))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating cartid refresh request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-review?precheckout=true")

	res, err := tk.FastClient.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making cartid refresh request")
		tk.Stop()
		return
	}

	tk.cartid = cartIdRe.FindStringSubmatch(string(res.Body))[1]
}

func (tk *Task) SubmitShipping() {
	var form string
	if *tk.Data.Profile.Shipping.ShippingAddress.AddressLine2 != "" {
		form = fmt.Sprintf(`{"cart_type":"REGULAR","address":{"address_line1":"%s","address_line2":"%s","address_type":"SHIPPING","city":"%s","country":"%s","first_name":"%s","last_name":"%s","mobile":"%s","save_as_default":false,"state":"%s","zip_code":"%s"},"selected":true,"save_to_profile":true,"skip_verification":true}`, tk.Data.Profile.Shipping.ShippingAddress.AddressLine, tk.Data.Profile.Shipping.ShippingAddress.AddressLine2, tk.Data.Profile.Shipping.ShippingAddress.City, tk.Data.Profile.Shipping.ShippingAddress.Country, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.ZIP)
	} else {
		form = fmt.Sprintf(`{"cart_type":"REGULAR","address":{"address_line1":"%s","address_type":"SHIPPING","city":"%s","country":"%s","first_name":"%s","last_name":"%s","mobile":"%s","save_as_default":false,"state":"%s","zip_code":"%s"},"selected":true,"save_to_profile":true,"skip_verification":true}`, tk.Data.Profile.Shipping.ShippingAddress.AddressLine, tk.Data.Profile.Shipping.ShippingAddress.City, tk.Data.Profile.Shipping.ShippingAddress.Country, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.ZIP)
	}

	req, err := tk.FastClient.NewRequest("POST", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/cart_shipping_addresses?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(form))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating shipping request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-delivery")

	_, err = tk.FastClient.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making shipping request")
		tk.Stop()
		return
	}

}

func (tk *Task) SubmitPayment() {
	var form string

	if *tk.Data.Profile.Shipping.BillingAddress.AddressLine2 != "" {
		form = fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"card_name":"%s","card_number":"%s","cvv":"%s","expiry_month":"%s","expiry_year":"%s"},"billing_address":{"address_line1":"%s","address_line2":"%s","city":"%s","first_name":"%s","last_name":"%s","phone":"%s","state":"%s","zip_code":"%s","country":"%s"}}`, tk.cartid, tk.Data.Profile.Shipping.FirstName+" "+tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Billing.Number, tk.Data.Profile.Billing.CVV, tk.Data.Profile.Billing.ExpirationMonth, "20"+tk.Data.Profile.Billing.ExpirationYear, tk.Data.Profile.Shipping.BillingAddress.AddressLine, *tk.Data.Profile.Shipping.BillingAddress.AddressLine2, tk.Data.Profile.Shipping.BillingAddress.City, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.BillingAddress.StateCode, tk.Data.Profile.Shipping.BillingAddress.ZIP, tk.Data.Profile.Shipping.BillingAddress.Country)
	} else {
		form = fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"card_name":"%s","card_number":"%s","cvv":"%s","expiry_month":"%s","expiry_year":"%s"},"billing_address":{"address_line1":"%s","city":"%s","first_name":"%s","last_name":"%s","phone":"%s","state":"%s","zip_code":"%s","country":"%s"}}`, tk.cartid, tk.Data.Profile.Shipping.FirstName+" "+tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Billing.Number, tk.Data.Profile.Billing.CVV, tk.Data.Profile.Billing.ExpirationMonth, "20"+tk.Data.Profile.Billing.ExpirationYear, tk.Data.Profile.Shipping.BillingAddress.AddressLine, tk.Data.Profile.Shipping.BillingAddress.City, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.BillingAddress.StateCode, tk.Data.Profile.Shipping.BillingAddress.ZIP, tk.Data.Profile.Shipping.BillingAddress.Country)
	}

	req, err := tk.FastClient.NewRequest("POST", fmt.Sprintf("https://carts.target.com/checkout_payments/v1/payment_instructions?key=%s", tk.apikey), []byte(form))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating payment request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-payment")

	res, err := tk.FastClient.Do(req)
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
	req, err := tk.FastClient.NewRequest("POST", fmt.Sprintf("https://carts.target.com/checkout_payments/v1/payment_instructions/%s?key=%s", tk.paymentinstructionid, tk.apikey), []byte(fmt.Sprintf(`{"cart_id":"%s","card_number":"%s"}`, tk.cartid, tk.Data.Profile.Billing.CVV)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating compare card request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-payment")

	res, err := tk.FastClient.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making compare card request")
		tk.Stop()
		return
	}
	if res.StatusCode == 401 {
		tk.RefreshToken()
		tk.CompareCard()
	}

	if strings.Contains(string(res.Body), "SUCCESS") {
		tk.SetStatus(module.STATUS_CONTINUING, "card valid")
	} else {
		tk.SetStatus(module.STATUS_ERROR, "card not valid")
		tk.Stop()
		return
	}
}

func (tk *Task) SubmitCVV() {
	req, err := tk.FastClient.NewRequest("PUT", fmt.Sprintf("https://carts.target.com/checkout_payments/v1/credit_card_compare?key=%s", tk.apikey), []byte(fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"cvv":"%s"}}`, tk.cartid, tk.Data.Profile.Billing.Number)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating compare card request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-payment")

	res, err := tk.FastClient.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making cvv request")
		tk.Stop()
		return
	}

	if res.StatusCode == 401 {
		tk.RefreshToken()
		tk.SubmitCVV()
	} else {
		tk.SetStatus(module.STATUS_CONTINUING, "payment submitted")
	}
}

func (tk *Task) SubmitCheckout() {
	req, err := tk.FastClient.NewRequest("PUT", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/checkout?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CDELIVERY_WINDOWS%%2CPAYMENT_INSTRUCTIONS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.apikey), []byte(`{"cart_type":"REGULAR","channel_id":10}`))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating compare card request")
		tk.Stop()
		return
	}
	res, err := tk.FastClient.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making checkout request")
		tk.Stop()
		return
	}

	var orderdata *CheckoutResponse
	json.Unmarshal(res.Body, &orderdata)

	tk.ReturningFields.Size = "os"
	tk.ReturningFields.ProductImage = fmt.Sprintf(`https://target.scene7.com/is/image/Target/GUEST_%s?wid=175&hei=175&qlt=80&fmt=webp`, tk.imageguestid)
	tk.ReturningFields.Color = "na"
	tk.ReturningFields.OrderNumber = orderdata.Orders[0].OrderID

	if res.StatusCode == 401 {
		tk.RefreshToken()
		tk.SubmitCheckout()
	} else {
		tk.SetStatus(module.STATUS_CHECKED_OUT, "checked out")
	}
}
