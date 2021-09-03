package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"strings"
)

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
		tk.RefreshToken()
		tk.NearestStore()
		return
	}

	tk.storeid = locationIdRe.FindStringSubmatch(string(res.Body))[1]
}

func (tk *Task) SubmitShipping() {
	tk.SetStatus(module.STATUS_SUBMITTING_SHIPPING)

	//threatmatrix.SendRequests(tk.cartid)

	var form string
	if tk.Data.Profile.Shipping.ShippingAddress.AddressLine2 != nil {
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

	res, err := tk.Do(req)
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error making shipping request")
		tk.Stop()
		return
	}

	if strings.Contains(string(res.Body), "ADDRESS_ALREADY_PRESENT"){
		return
	}
}