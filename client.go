package robinhood

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"astuart.co/clyde"
)

const (
	epBase        = "https://api.robinhood.com/"
	epLogin       = epBase + "oauth2/token/"
	epAccounts    = epBase + "accounts/"
	epQuotes      = epBase + "quotes/"
	epPortfolios  = epBase + "portfolios/"
	epWatchlists  = epBase + "watchlists/"
	epInstruments = epBase + "instruments/"
	epOrders      = epBase + "orders/"
)

type Client struct {
	Token string
	*http.Client
}

func Dial(t TokenGetter) (*Client, error) {
	tkn, err := t.GetToken()
	if err != nil {
		return nil, err
	}

	return &Client{
		Token:  tkn,
		Client: &http.Client{Transport: clyde.HeaderRoundTripper{"Authorization": "Bearer " + tkn}},
	}, nil
}

func (c *Client) GetAndDecode(url string, dest Detailable) error {
	res, err := c.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return unmarshalJSON(res.Body, dest)
}

type Detailable interface {
	Details() string
}

func (c *Client) PostAndDecode(url string, data interface{}, dest Detailable) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if DebugMode {
		fmt.Println(string(bytes))
	}
	res, err := c.Post(url, "application/json", strings.NewReader(string(bytes)))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return unmarshalJSON(res.Body, dest)
}

func unauthenticatedPostAndDecode(url string, data interface{}, dest Detailable) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if DebugMode {
		fmt.Println(string(bytes))
	}
	res, err := http.Post(url, "application/json", strings.NewReader(string(bytes)))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		fmt.Printf("Got status code %v from URL %v\n", res.StatusCode, url)
	}

	defer res.Body.Close()

	return unmarshalJSON(res.Body, dest)
}

// unmarshalJSON wraps json.Unmarshal
// It will log the response if DebugMode is true
func unmarshalJSON(r io.Reader, dest Detailable) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if DebugMode {
		fmt.Println(string(body))
	}
	err = json.Unmarshal(body, &dest)
	if err != nil {
		return err
	}
	deets := dest.Details()
	if deets != "" {
		return fmt.Errorf("response contained error %v", deets)
	}
	return nil
}

type Meta struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
}
