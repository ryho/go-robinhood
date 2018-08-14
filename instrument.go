package robinhood

type Instrument struct {
	BloombergUnique    string      `json:"bloomberg_unique"`
	Country            string      `json:"country"`
	DayTradeRatio      string      `json:"day_trade_ratio"`
	Fundamentals       string      `json:"fundamentals"`
	ID                 string      `json:"id"`
	ListDate           string      `json:"list_date"`
	MaintenanceRatio   string      `json:"maintenance_ratio"`
	MarginInitialRatio string      `json:"margin_initial_ratio"`
	Market             string      `json:"market"`
	MinTickSize        interface{} `json:"min_tick_size"`
	Name               string      `json:"name"`
	Quote              string      `json:"quote"`
	SimpleName         interface{} `json:"simple_name"`
	Splits             string      `json:"splits"`
	State              string      `json:"state"`
	Symbol             string      `json:"symbol"`
	Tradeable          bool        `json:"tradeable"`
	URL                string      `json:"url"`
	Detail             string      `json:"detail"`
}

func (resp *Instrument) Details() string {
	return resp.Detail
}

func (c Client) GetInstrument(instURL string) (*Instrument, error) {
	var i Instrument
	err := c.GetAndDecode(instURL, &i)
	return &i, err
}

type GetInstrumentsResponse struct {
	Results []Instrument
	Detail  string `json:"detail"`
}

func (resp *GetInstrumentsResponse) Details() string {
	return resp.Detail
}

func (c Client) GetInstrumentForSymbol(sym string) (*Instrument, error) {
	var i GetInstrumentsResponse
	err := c.GetAndDecode(epInstruments+"?symbol="+sym, &i)
	return &i.Results[0], err
}
