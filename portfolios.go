package robinhood

type Portfolio struct {
	Account                                string  `json:"account"`
	AdjustedEquityPreviousClose            float64 `json:"adjusted_equity_previous_close,string"`
	Equity                                 float64 `json:"equity,string"`
	EquityPreviousClose                    float64 `json:"equity_previous_close,string"`
	ExcessMaintenance                      float64 `json:"excess_maintenance,string"`
	ExcessMaintenanceWithUnclearedDeposits float64 `json:"excess_maintenance_with_uncleared_deposits,string"`
	ExcessMargin                           float64 `json:"excess_margin,string"`
	ExcessMarginWithUnclearedDeposits      float64 `json:"excess_margin_with_uncleared_deposits,string"`
	ExtendedHoursEquity                    float64 `json:"extended_hours_equity,string"`
	ExtendedHoursMarketValue               float64 `json:"extended_hours_market_value,string"`
	LastCoreEquity                         float64 `json:"last_core_equity,string"`
	LastCoreMarketValue                    float64 `json:"last_core_market_value,string"`
	MarketValue                            float64 `json:"market_value,string"`
	StartDate                              string  `json:"start_date"`
	UnwithdrawableDeposits                 float64 `json:"unwithdrawable_deposits,string"`
	UnwithdrawableGrants                   float64 `json:"unwithdrawable_grants,string"`
	URL                                    string  `json:"url"`
	WithdrawableAmount                     float64 `json:"withdrawable_amount,string"`
}

type GetPortfolioResponse struct {
	Results []Portfolio
	Detail  string `json:"detail"`
}

func (resp *GetPortfolioResponse) Details() string {
	return resp.Detail
}

// GetPortfolios returns all the portfolios associated with a client's
// credentials and accounts
func (c *Client) GetPortfolios() ([]Portfolio, error) {
	var p GetPortfolioResponse
	err := c.GetAndDecode(epPortfolios, &p)
	return p.Results, err
}
