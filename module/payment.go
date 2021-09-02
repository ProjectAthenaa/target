package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"strconv"
	"strings"
)

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
