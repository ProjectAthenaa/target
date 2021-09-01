package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"regexp"
	"sync"
	"sync/atomic"
)

var (
	imageRe = regexp.MustCompile(`"og:image" content=("https://target.scene7.com/is/image/Target/GUEST_[^\"]+?")`)
	tcinRe = regexp.MustCompile(`"tcin":"\d+?"`)
)

type StockInfo struct {
	Data struct {
		Product struct {
			Fulfillment struct {
				ShippingOptions struct {
					AvailabilityStatus  string  `json:"availability_status"`
				} `json:"shipping_options"`
			} `json:"fulfillment"`
		} `json:"product"`
	} `json:"data"`
}

func (tk *Task) InitData(){
	req, err := tk.NewRequest("GET", *tk.Data.TaskData.Link, nil)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "could not fetch product page")
		tk.Stop()
		return
	}
	res, err := tk.Do(req)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "could not read product page")
		tk.Stop()
		return
	}

	tk.imagelink = string(imageRe.FindSubmatch(res.Body)[1])
	tk.pid = string(tcinRe.FindSubmatch(res.Body)[1])
}

func (tk *Task) WaitForInstock() {
	req, err := tk.NewRequest("GET", fmt.Sprintf(`https://redsky.target.com/redsky_aggregations/v1/web/pdp_fulfillment_v1?key=%s&tcin=%s&store_id=%s&store_positions_store_id=%s&has_store_positions_store_id=true&zip=%s&state=NJ&scheduled_delivery_store_id=%s&pricing_store_id=%s&has_pricing_store_id=true&is_bot=false`, tk.apikey, tk.pid, tk.storeid, tk.storeid, tk.Data.Profile.Shipping.ShippingAddress.ZIP,tk.storeid, tk.storeid ), nil)
	if err != nil{
		tk.SetStatus(module.STATUS_ERROR, "could not fetch product availability")
		tk.Stop()
		return
	}

	var found int32
	var wg sync.WaitGroup

	tk.SetStatus(module.STATUS_MONITORING, "waiting for instock")

	for i := 0 ; i < 5; i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			for found == 0{
				res, err := tk.Do(req)
				if err != nil{
					tk.SetStatus(module.STATUS_ERROR, "could not read product availability")
					tk.Stop()
					return
				}

				var instock *StockInfo
				json.Unmarshal(res.Body, &instock)
				if instock.Data.Product.Fulfillment.ShippingOptions.AvailabilityStatus == "IN_STOCK"{
					atomic.AddInt32(&found, 1)
				}
			}
		}()
	}

	wg.Wait()
	tk.SetStatus(module.STATUS_PRODUCT_FOUND, "")
}