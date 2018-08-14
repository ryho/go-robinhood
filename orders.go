package robinhood

import "time"

type OrderRequest struct {
	Account    string `json:"account"`
	Instrument string `json:"instrument"`
	Symbol     string `json:"symbol"`
	//Order type: market or limit
	Type OrderType `json:"type"`
	//gfd, gtc, ioc, or opg
	TimeInForce TimeInForce `json:"time_in_force"`
	// immediate or stop
	Trigger Trigger `json:"trigger"`
	// for use with limit
	Price float64 `json:"price"`
	// required when trigger equals stop
	StopPrice float64 `json:"stop_price,omitempty"`
	Quantity  int     `json:"quantity"`
	// buy or sell
	Side Side `json:"side"`
	//Would/Should order execute when exchanges are closed
	ExtendedHours          bool `json:"extended_hours"`
	OverrideDayTradeChecks bool `json:"override_day_trade_checks"`
	OverrideDtbpChecks     bool `json:"override_dtbp_checks"`
}
type OrderType string

const (
	OrderType_Market OrderType = "market"
	OrderType_Limit  OrderType = "limit"
)

type TimeInForce string

const (
	TimeInForce_GoodForDay        TimeInForce = "gfd"
	TimeInForce_GoodTillCanceled  TimeInForce = "gtc"
	TimeInForce_ImmediateOrCancel TimeInForce = "ioc"
	TimeInForce_Opening           TimeInForce = "opg"
)

type Trigger string

const (
	Trigger_Imediate Trigger = "immediate"
	Trigger_Stop     Trigger = "stop"
)

type Side string

const (
	Side_Buy  Side = "buy"
	Side_Sell Side = "sell"
)

type Order struct {
	Meta
	Id                 string      `json:"id"`
	Executions         []Execution `json:"executions"`
	Fees               float64     `json:"fees,string"`
	Cancel             string      `json:"cancel"`
	CumulativeQuantity float64     `json:"cumulative_quantity,string"`
	RejectReason       string      `json:"reject_reason"`
	//queued, unconfirmed, confirmed, partially_filled, filled, rejected, canceled, or failed
	State OrderState `json:"state"`
	// required when trigger equals stop
	LastTransactionAt      string  `json:"last_transaction_at"`
	ClientId               string  `json:"client_id"`
	URL                    string  `json:"url"`
	Position               string  `json:"position"`
	AveragePrice           float64 `json:"average_price,string"`
	ExtendedHours          bool    `json:"extended_hours"`
	OverrideDayTradeChecks bool    `json:"override_day_trade_checks"`
	OverrideDtbpChecks     bool    `json:"override_dtbp_checks"`

	Detail string `json:"detail"`
}

type Execution struct {
	Id             string
	Price          float64   `json:"price,string"`
	Quantity       float64   `json:"quantity,string"`
	Timestamp      time.Time `json:"timestamp,string"`
	SettlementDate string    `json:"settlement_date"`
}

type OrderState string

const (
	OrderState_Queued          OrderState = "queued"
	OrderState_Unconfirmed     OrderState = "unconfirmed"
	OrderState_Confirmed       OrderState = "confirmed"
	OrderState_PartiallyFilled OrderState = "partially_filled"
	OrderState_Filled          OrderState = "filled"
	OrderState_Rejected        OrderState = "rejected"
	OrderState_Canceled        OrderState = "canceled"
	OrderState_Failed          OrderState = "failed"
)

func (resp *Order) Details() string {
	return resp.Detail
}

// SendOrder will send an order to buy or sell
func (c *Client) SendOrder(request *OrderRequest) (Order, error) {
	var response Order
	err := c.PostAndDecode(epOrders, request, &response)
	return response, err
}

// GetOrder returns the order with the given id
func (c *Client) GetOrder(id string) (Order, error) {
	var response Order
	err := c.GetAndDecode(epOrders+id, &response)
	return response, err
}

type GetOrderRequest struct {
	Cursor     string    `json:"cursor,omitempty"`
	Instrument string    `json:"instrument,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type GetOrderResponse struct {
	Previous string `json:"prevoius"`
	Next     string `json:"next"`
	Results  []Order
	Detail   string `json:"detail"`
}

func (resp *GetOrderResponse) Details() string {
	return resp.Detail
}

// GetRecentOrders returns all recent orders for the instrument
func (c *Client) GetRecentOrders(id *Instrument) ([]Order, error) {
	var orders []Order
	url := epOrders + "?instrument=" + id.URL
	for {
		var response GetOrderResponse
		err := c.GetAndDecode(url, &response)
		if err != nil {
			return nil, err
		}
		orders = append(orders, response.Results...)
		if response.Next == "" {
			break
		}
		url = response.Next
		time.Sleep(time.Second)
	}

	return orders, nil
}

type CancelOrderResponse struct {
	Detail string `json:"detail"`
}

func (resp *CancelOrderResponse) Details() string {
	return resp.Detail
}

// CancelOrder will cancel the order with the given id
func (c *Client) CancelOrder(id string) error {
	var r CancelOrderResponse
	return c.PostAndDecode(epOrders+id+"/cancel/", struct{}{}, &r)
}
