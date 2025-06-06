package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/antibots/shape"
	"github.com/ProjectAthenaa/target/config"
)

func (tk *Task) ATC() {

	tk.SetStatus(module.STATUS_ADDING_TO_CART)

	req, err := tk.NewRequest("POST", fmt.Sprintf("https://carts.target.com/web_checkouts/v1/cart_items?field_groups=CART%%2CCART_ITEMS%%2CSUMMARY%%2CFINANCE_PROVIDERS&key=%s", tk.cartApiKey), []byte(fmt.Sprintf(`{"cart_type":"REGULAR","channel_id":"10","shopping_context":"DIGITAL","cart_item":{"tcin":"%s","quantity":1,"item_channel_id":"10"}}`, tk.pid)))
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error creating atc request")
		tk.Stop()
		return
	}
	req.Headers = tk.GenerateDefaultHeaders(fmt.Sprintf("https://www.target.com/p/-/A-%s", tk.Data.Metadata[*config.Module.Fields[0].FieldKey]))
	headers, err := shapeClient.GenHeaders(tk.Ctx, &shape.Site{Value: shape.SITE_TARGET, ResString: tk.GetResString()})
	if err != nil {
		tk.SetStatus(module.STATUS_ERROR, "error generating shape headers")
		tk.Stop()
		return
	}

	for k, v := range headers.Values {
		req.Headers[k] = []string{v}
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
