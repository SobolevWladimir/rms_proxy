package parameters

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RMSConnectParameter struct {
	URL             string `json:"url"`
	Login           string `json:"login"`
	Password        string
	MeedPassEncrupt bool `json:"needPassEncrupt"`
}

func (rm *RMSConnectParameter) Handle(r *http.Request) (*http.Response, error) {
	respToken, err := rm.GetToken()
	if r.URL.Path == "/resto/api/auth" {
		return respToken, err
	}
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(respToken.Body)
	token := string(b)
	defer rm.Logout(token)

	return rm.Proxy(r, token)
}

func (rm *RMSConnectParameter) Logout(token string) error {
	uri, err := url.Parse(rm.URL)
	if err != nil {
		return nil
	}

	uri.Path = "/resto/api/logout"
	query := uri.Query()
	query.Set("key", rm.Login)
	query.Set("client-type", "iikoweb")
	uri.RawQuery = query.Encode()
	_, err = http.Get(uri.String())
	return err
}

func (rm *RMSConnectParameter) GetToken() (*http.Response, error) {
	uri, err := url.Parse(rm.URL)
	if err != nil {
		return nil, err
	}
	uri.Path = "/resto/api/auth"
	query := uri.Query()
	query.Set("login", rm.Login)
	if rm.MeedPassEncrupt {
		hasher := sha1.New()
		hasher.Write([]byte(rm.Password))
		pass := hex.EncodeToString(hasher.Sum(nil))
		query.Set("pass", pass)
	} else {
		query.Set("pass", rm.Password)
	}
	uri.RawQuery = query.Encode()
	resp, err := http.Get(uri.String())
	return resp, err
}

func (rm *RMSConnectParameter) Proxy(r *http.Request, token string) (*http.Response, error) {
	uri, err := url.Parse(rm.URL)
	if err != nil {
		return nil, err
	}
	uri.Path = r.URL.Path
	query := r.URL.Query()
	query.Set("key", token)
	uri.RawQuery = query.Encode()
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, uri.String(), bytes.NewBuffer(requestBody))
	req.Header = rm.getHeaders(r)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	return resp, err
}

func (rm *RMSConnectParameter) getHeaders(r *http.Request) http.Header {
	result := http.Header{}
	for key, value := range r.Header {
		if strings.Contains(key, "X-Resto") {
			result.Set(key, value[0])
		}
		if strings.Contains(key, "Content-Type") {
			result.Set(key, value[0])
		}
	}

	return result
}
