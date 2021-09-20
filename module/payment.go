package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"strconv"
	"strings"
)

func (tk *Task) RefreshCartId() {
	tk.SetStatus(module.STATUS_CHECKING_OUT, "refreshing cart")
	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/pre_checkout?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CDELIVERY_WINDOWS%%2CPAYMENT_INSTRUCTIONS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.cartApiKey), []byte(`{"cart_type":"REGULAR"}`))
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

	if v := paymentInstructionRe.FindSubmatch(res.Body); len(v) == 2 {
		tk.paymentinstructionid = string(v[1])
		tk.ReturningFields.Price = string(orderTotalRe.FindSubmatch(res.Body)[1])
	} else {
		fmt.Println(orderTotalRe.FindStringSubmatch(string(res.Body)))
		tk.ReturningFields.Price = string(orderTotalRe.FindSubmatch(res.Body)[1])
	}
}

func (tk *Task) SubmitPayment() {
	tk.SetStatus(module.STATUS_SUBMITTING_PAYMENT)
	var form string

	if tk.Data.Profile.Shipping.BillingIsShipping {
		if tk.Data.Profile.Shipping.ShippingAddress.AddressLine2 != nil {
			form = fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"card_name":"%s","card_number":"%s","cvv":"%s","expiry_month":"%s","expiry_year":"%s"},"billing_address":{"address_line1":"%s","address_line2":"%s","city":"%s","first_name":"%s","last_name":"%s","phone":"%s","state":"%s","zip_code":"%s","country":"%s"}}`, tk.cartid, tk.Data.Profile.Shipping.FirstName+" "+tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Billing.Number, tk.Data.Profile.Billing.CVV, tk.Data.Profile.Billing.ExpirationMonth, "20"+tk.Data.Profile.Billing.ExpirationYear, tk.Data.Profile.Shipping.ShippingAddress.AddressLine, *tk.Data.Profile.Shipping.ShippingAddress.AddressLine2, tk.Data.Profile.Shipping.ShippingAddress.City, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.ZIP, tk.Data.Profile.Shipping.ShippingAddress.Country)
		} else {
			form = fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"card_name":"%s","card_number":"%s","cvv":"%s","expiry_month":"%s","expiry_year":"%s"},"billing_address":{"address_line1":"%s","city":"%s","first_name":"%s","last_name":"%s","phone":"%s","state":"%s","zip_code":"%s","country":"%s"}}`, tk.cartid, tk.Data.Profile.Shipping.FirstName+" "+tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Billing.Number, tk.Data.Profile.Billing.CVV, tk.Data.Profile.Billing.ExpirationMonth, "20"+tk.Data.Profile.Billing.ExpirationYear, tk.Data.Profile.Shipping.ShippingAddress.AddressLine, tk.Data.Profile.Shipping.ShippingAddress.City, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.ShippingAddress.StateCode, tk.Data.Profile.Shipping.ShippingAddress.ZIP, tk.Data.Profile.Shipping.ShippingAddress.Country)
		}
	} else {
		if tk.Data.Profile.Shipping.BillingAddress.AddressLine2 != nil {
			form = fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"card_name":"%s","card_number":"%s","cvv":"%s","expiry_month":"%s","expiry_year":"%s"},"billing_address":{"address_line1":"%s","address_line2":"%s","city":"%s","first_name":"%s","last_name":"%s","phone":"%s","state":"%s","zip_code":"%s","country":"%s"}}`, tk.cartid, tk.Data.Profile.Shipping.FirstName+" "+tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Billing.Number, tk.Data.Profile.Billing.CVV, tk.Data.Profile.Billing.ExpirationMonth, "20"+tk.Data.Profile.Billing.ExpirationYear, tk.Data.Profile.Shipping.BillingAddress.AddressLine, *tk.Data.Profile.Shipping.BillingAddress.AddressLine2, tk.Data.Profile.Shipping.BillingAddress.City, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.BillingAddress.StateCode, tk.Data.Profile.Shipping.BillingAddress.ZIP, tk.Data.Profile.Shipping.BillingAddress.Country)
		} else {
			form = fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"card_name":"%s","card_number":"%s","cvv":"%s","expiry_month":"%s","expiry_year":"%s"},"billing_address":{"address_line1":"%s","city":"%s","first_name":"%s","last_name":"%s","phone":"%s","state":"%s","zip_code":"%s","country":"%s"}}`, tk.cartid, tk.Data.Profile.Shipping.FirstName+" "+tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Billing.Number, tk.Data.Profile.Billing.CVV, tk.Data.Profile.Billing.ExpirationMonth, "20"+tk.Data.Profile.Billing.ExpirationYear, tk.Data.Profile.Shipping.BillingAddress.AddressLine, tk.Data.Profile.Shipping.BillingAddress.City, tk.Data.Profile.Shipping.FirstName, tk.Data.Profile.Shipping.LastName, tk.Data.Profile.Shipping.PhoneNumber, tk.Data.Profile.Shipping.BillingAddress.StateCode, tk.Data.Profile.Shipping.BillingAddress.ZIP, tk.Data.Profile.Shipping.BillingAddress.Country)
		}
	}

	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/checkout_payments/v1/payment_instructions?key=%s", tk.cartApiKey), []byte(form))
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

	if strings.Contains(string(res.Body), "CARD_PAYMENT_EXISTS") {
		return
	} else {
		var instructionresponse *PaymentInstructions
		json.Unmarshal(res.Body, &instructionresponse)

		tk.paymentinstructionid = instructionresponse.PaymentInstructionID
		tk.ReturningFields.Price = strconv.FormatFloat(instructionresponse.PaymentInstructionAmount, 'f', -1, 64)
	}
}

func (tk *Task) CompareCard() {
	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/checkout_payments/v1/credit_card_compare?key=%s", tk.cartApiKey), []byte(fmt.Sprintf(`{"cart_id":"%s","card_number":"%s"}`, tk.cartid, tk.Data.Profile.Billing.Number)))
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
		//tk.RefreshToken()
		tk.CompareCard()
		return
	}

	if strings.Contains(string(res.Body), "SUCCESS") {
		tk.SetStatus(module.STATUS_CHECKING_OUT, "card valid")
	} else {
		tk.SetStatus(module.STATUS_CHECKOUT_DECLINE, "card not valid")
		tk.Stop()
		return
	}
}

func (tk *Task) SubmitCVV() {
	if !tk.submitCVV {
		return
	}
	req, err := tk.NewRequest("PUT", fmt.Sprintf("https://carts.target.com/checkout_payments/v1/payment_instructions/%s?key=%s", tk.paymentinstructionid, tk.cartApiKey), []byte(fmt.Sprintf(`{"cart_id":"%s","wallet_mode":"NONE","payment_type":"CARD","card_details":{"cvv":"%s"}}`, tk.cartid, tk.Data.Profile.Billing.CVV)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating cvv request")
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
		//tk.RefreshToken()
		tk.SubmitCVV()
		return
	} else {
		tk.SetStatus(module.STATUS_CHECKING_OUT, "payment submitted")
	}
}

func (tk *Task) PaymentOauth() {
	req, err := tk.NewRequest("POST", "https://gsp.target.com/gsp/oauth_validations/v3/token_validations", []byte(`{}`))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating payment oauth request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-payment")

	_, err = tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making payment oauth request")
		tk.Stop()
		return
	}

}

func (tk *Task) SubmitCheckout() {
	tk.SetStatus(module.STATUS_SUBMITTING_CHECKOUT)
	//req, err := tk.NewRequest("POST", `https://carts.target.com/web_checkouts/v1/checkout?field_groups=ADDRESSES%2CCART%2CCART_ITEMS%2CDELIVERY_WINDOWS%2CPAYMENT_INSTRUCTIONS%2CPICKUP_INSTRUCTIONS%2CPROMOTION_CODES%2CSUMMARY%2CFINANCE_PROVIDERS&key=feaf228eb2777fd3eee0fd5192ae7107d6224b39`, []byte(`{"cart_type":"REGULAR","channel_id":10}`))
	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/checkout?field_groups=ADDRESSES%%2CCART%%2CCART_ITEMS%%2CDELIVERY_WINDOWS%%2CPAYMENT_INSTRUCTIONS%%2CPICKUP_INSTRUCTIONS%%2CPROMOTION_CODES%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.cartApiKey), []byte(`{"cart_type":"REGULAR","channel_id":"10"}`))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating compare card request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders("https://www.target.com/co-review")

	//cookiejar.ReleaseCookieJar(tk.FastClient.Jar)
	//tk.FastClient.Jar = nil
	//req.Headers["Cookie"] = []string{
	//	`TealeafAkaSid=1qO9JRb606YJBVl-mD_ODKToDqTy4hqu;visitorId=017C01BAB1410201B3864A91B7E33229;sapphire=1;UserLocation=07093|40.790|-74.020|NJ|US;adaptiveSessionId=A5843947546;egsSessionId=dabe2cbb-7454-4866-9a89-0f7535fb5d27;fiatsCookie=DSI_1865|DSN_North%20Bergen%20Commons|DSZ_07047;criteo={};tlThirdPartyIds={%22pt%22:%22v2:c4b3f7c735ea7f46fa912375296760ab69f5abadb7c505fd578a1b7b7bf048d1|8e934a7648e1c54b8d7ff86fdad39a22040c8a95d771edba1104d850c79011a4%22};ci_pixmgr=other;mystate=1632116586956;_gcl_au=1.1.1183904347.1632116587;login-session=Tv0AjgcJPPZY9iq8vNaxRtimuxZ0dMpTuVyf_-6TZ-do3cMk9c5WcIWbutd6AieY;3YCzT93n=A6HKugF8AQAAmMzN7mtZotZfuK1BxRnrf1gaO4uP_1fk7wjkFKGl7hyp-M6YAUha_m6ucvlowH8AAEB3AAAAAA|1|1|80163e867ee4de0ae908befcfa891c2d072cbf3b;accessToken=eyJraWQiOiJlYXMyIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiIxMDAzNTg3NDE3MiIsImlzcyI6Ik1JNiIsImV4cCI6MTYzMjEzMDk5NCwiaWF0IjoxNjMyMTE2NTk0LCJqdGkiOiJUR1QuZWU1N2NiYThhZTU0NDhkZTgxYTkwYTc5YTNjNmU3ZjMtbSIsInNreSI6ImVhczIiLCJzdXQiOiJSIiwiZGlkIjoiYmFmMjA4NjFhNzFiZDVlMzg2OTVjYTE4MjU5NjBhZGMzOTU1MGViMDc2MmZhYTczMmU4OWRjZDdhMTgwYTA0YSIsImVpZCI6Im5pZGFpdWR3YWl1QGdtYWlsLmNvbSIsImdzcyI6MS4wLCJzY28iOiJlY29tLm1lZCxvcGVuaWQiLCJjbGkiOiJlY29tLXdlYi0xLjAuMCIsImFzbCI6Ik0ifQ.fb1IfK_iwqtNvvFNRZ_rmKLSinaiuFCqBab5BD11qxilDLaCqwGN8iu_KI_8-X47SbLzKu3r3jhaUQVQ0g4Vmlj_LZX4HCRk_4yWw0Eb35BfMkaejoZL6C47HPO7dVb-BGndRAadHiafO4hJ_6QLXXtSK9LKyd5V89xGUjQRcyP9BiW5mdOAfuw4PoU3SQxrsP4Bsl9mwSx3zyyykDubp6Omhmz8Np0mceAQ1D3aVmMK2i8KSfrD88PilqSV1GGwMXOQ7IPcN2doVbsHeO_YPbAq8YOu2ljSyv3aX1NZs3TQ9IyJeLzLT2Q1rxlLM9zGymI1vXbFQP9StLpN-yciYg;idToken=eyJhbGciOiJub25lIn0.eyJzdWIiOiIxMDAzNTg3NDE3MiIsImlzcyI6Ik1JNiIsImV4cCI6MTYzMjEzMDk5NCwiaWF0IjoxNjMyMTE2NTk0LCJhc3MiOiJNIiwic3V0IjoiUiIsImNsaSI6ImVjb20td2ViLTEuMC4wIiwicHJvIjp7ImZuIjoib21hciIsImVtIjoibmlkYWl1ZHdhaXVAZ21haWwuY29tIiwicGgiOmZhbHNlLCJsZWQiOm51bGwsImx0eSI6ZmFsc2V9fQ.;refreshToken=TGT.ee57cba8ae5448de81a90a79a3c6e7f3-m;guestType=R|1632116594000;__gads=ID=2063e0fb6318582c-22fc39bc16bb0070:T=1632116595:S=ALNI_MY_MBGppGjeIyA4UdDQduPLNqeteg;mid=10035874172;cd_user_id=17c01bade8b814-08c2afe1e49c74-a7d173c-2a8f08-17c01bade8c92e;crl8.fpcuid=3884d1ca-f318-45ef-97ff-fae2efdaf2ec;_uetsid=a18cee9019d511ec9c071beaeb3bab4b;_uetvid=a18cf82019d511ec9763ed6e6b433469;ffsession={%22sessionHash%22:%221061873f7943d71632116585762%22%2C%22sessionHit%22:43%2C%22prevPageType%22:%22checkout%22%2C%22prevPageName%22:%22checkout:%20order%20review%22%2C%22prevPageUrl%22:%22https://www.target.com/co-review%22};targetMobileCookie=hasRC:false~cartQty:1~guestLogonId:nidaiudwaiu@gmail.com~guestDisplayName:omar~guestHasVerifiedPhone:false`,
	//}
	
	
	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making checkout request")
		tk.Stop()
		return
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		//tk.RefreshToken()
		tk.SubmitCheckout()
		return
	}

	if strings.Contains(string(res.Body), "order_id") {
		var orderdata *CheckoutResponse
		json.Unmarshal(res.Body, &orderdata)

		tk.ReturningFields.Size = "os"
		tk.ReturningFields.ProductImage = tk.imagelink
		tk.ReturningFields.Color = "na"
		tk.ReturningFields.OrderNumber = orderdata.Orders[0].OrderID
		tk.SetStatus(module.STATUS_CHECKED_OUT, "checked out")
	} else if strings.Contains(string(res.Body), "PAYMENT_DECLINED_EXCEPTION") {
		tk.SetStatus(module.STATUS_CHECKOUT_DECLINE, "declined")
	} else {
		errMessage := checkoutErrRe.FindStringSubmatch(string(res.Body))[1]
		tk.SetStatus(module.STATUS_CHECKOUT_ERROR, errMessage)
	}

}
