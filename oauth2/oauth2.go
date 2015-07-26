package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type OAGithub struct {
	sync.Mutex
	client_id        string
	client_secret    string
	scope            string
	auth_uri         string
	redirect_uri     string
	access_token_uri string
	user_info_uri    string
	state            map[string]string
}

func NewOAGithub(client_id, client_secret, scope string) *OAGithub {
	return &OAGithub{
		client_id:        client_id,
		client_secret:    client_secret,
		scope:            scope,
		auth_uri:         "https://github.com/login/oauth/authorize",
		redirect_uri:     "http://shaalx-oauths.daoapp.io/callback",
		access_token_uri: "https://github.com/login/oauth/access_token",
		user_info_uri:    "https://api.github.com/user?access_token",
		state:            make(map[string]string),
	}
}

func (oa *OAGithub) AuthURL() string {
	oa.Lock()
	state := randomStr()
	oa.state[state] = state
	oa.Unlock()
	vals := url.Values{
		"client_id":    {oa.client_id},
		"redirect_uri": {oa.redirect_uri},
		"scope":        {oa.scope},
		"state":        {state},
	}
	url_, _ := url.Parse(oa.auth_uri)
	url_.RawQuery = vals.Encode()
	return url_.String()
}

func randomStr() string {
	return "29s-sdfwuefs"
}

func (oa *OAGithub) AuthCode(req *http.Request) (code string, err error) {
	url_ := req.URL
	q := url_.Query()
	state := q.Get("state")
	if _, ok := oa.state[state]; !ok {
		return "", errors.New("SCRF Attack!")
	}
	delete(oa.state, state)
	code = q.Get("code")
	return code, nil
}

func (oa *OAGithub) AccessTokenURL(code string) string {
	vals := url.Values{
		"client_id":     {oa.client_id},
		"client_secret": {oa.client_secret},
		"code":          {code},
		"redirect_uri":  {oa.redirect_uri},
	}
	url_, _ := url.Parse(oa.access_token_uri)
	url_.RawQuery = vals.Encode()
	return url_.String()
}

func (oa *OAGithub) AccessToken(code string) (token string, err error) {
	access_token_url := oa.AccessTokenURL(code)
	_req, err := http.NewRequest("POST", access_token_url, nil)
	_req.Header.Set("Accept", "application/json")
	if nil != err {
		return "nil", err
	}
	c := http.Client{}
	resp, err := c.Do(_req)
	if nil != err {
		return "nil", err
	}
	var ret map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if nil != err {
		return "nil", err
	}
	return fmt.Sprintf("%v", ret["access_token"]), nil
}

// GET https://api.github.com/user?access_token
func (oa *OAGithub) UserInfoURL(access_token string) string {
	vals := url.Values{
		"access_token": {access_token},
	}
	url_, _ := url.Parse(oa.user_info_uri)
	url_.RawQuery = vals.Encode()
	return url_.String()
}

func (oa *OAGithub) UserInfo(access_token string) ([]byte, error) {
	user_info_url := oa.UserInfoURL(access_token)
	resp, err := http.Get(user_info_url)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (oa *OAGithub) NextStep(req *http.Request) ([]byte, error) {
	code, err := oa.AuthCode(req)
	if nil != err {
		return nil, err
	}
	fmt.Printf("code :  %s\n", code)
	token, err := oa.AccessToken(code)
	if nil != err {
		return nil, err
	}
	fmt.Printf("token :  %s\n", token)
	b, err := oa.UserInfo(token)
	if nil != err {
		return nil, err
	}
	return b, nil
}
