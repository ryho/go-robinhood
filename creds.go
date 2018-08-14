package robinhood

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var DebugMode bool

type TokenGetter interface {
	GetToken() (string, error)
}

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
	MFA      string `json:"mfa_code"`
	// Expiration time in seconds
	ExpiresIn int    `json:"expires_in"`
	Scope     string `json:"scope"`
	ClientId  string `json:"client_id"`
	GrantType string `json:"grant_type"`
}

type LoginResponse struct {
	Token       string `json:"access_token"`
	MFAType     string `json:"mfa_type"`
	MFARequired bool   `json:"mfa_required"`
	Detail      string `json:"detail"`
}

func (resp *LoginResponse) Details() string {
	return resp.Detail
}

func NewCreds(username, password string) *Creds {
	return NewCredsWithMFA(username, password, "")
}

func NewCredsWithMFA(username, password, mfa string) *Creds {
	return &Creds{
		Username: username,
		Password: password,
		MFA:      mfa,
		// 72 hours, in seconds
		ExpiresIn: int((72 * time.Hour) / time.Second),
		// These are the values that the Robinhood Website uses:
		Scope:     "internal",
		ClientId:  "c82SH0WZOsabOXGP2sxqcj34FxkvfnWRZBKlBjFS",
		GrantType: "password",
	}
}

func (c *Creds) GetToken() (string, error) {
	var resp LoginResponse
	err := unauthenticatedPostAndDecode(epLogin, c, &resp)
	if err != nil {
		return "", err
	}
	if resp.MFARequired {
		return "", fmt.Errorf("this account requires two factor. Two factor type: %v", resp.MFAType)
	}
	return resp.Token, nil
}

// A CredsCacher takes user credentials and a file path. The token obtained
// from the RobinHood API will be cached at the file path, and a new token will
// not be obtained.
type CredsCacher struct {
	Creds TokenGetter
	Path  string
}

// GetToken implements TokenGetter. It may fail if an error is encountered
// checking the file path provided, or if the underlying creds return an error
// when retrieving their token.
func (c *CredsCacher) GetToken() (string, error) {
	mustLogin := false

	err := os.MkdirAll(path.Dir(c.Path), 0750)
	if err != nil {
		return "", fmt.Errorf("error creating path for token: %s", err)
	}

	_, err = os.Stat(c.Path)
	if err != nil {
		if os.IsNotExist(err) {
			mustLogin = true
		} else {
			return "", err
		}
	}

	if !mustLogin {
		bs, err := ioutil.ReadFile(c.Path)
		bsstr := string(bs)
		if err != nil || bsstr != "" {
			return bsstr, err
		}
	}

	f, err := os.OpenFile(c.Path, os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		return "", err
	}
	defer f.Close()

	tok, err := c.Creds.GetToken()
	if err != nil {
		return "", err
	}

	if tok == "" {
		return "", fmt.Errorf("Empty token")
	}

	_, err = f.Write([]byte(tok))
	return tok, err
}

type Token string

func (t *Token) GetToken() (string, error) {
	return string(*t), nil
}
