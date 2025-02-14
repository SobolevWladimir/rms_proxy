package parameters

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RMSConnectParameter struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Login           string `json:"login"`
	Password        string
	NeedPassEncrupt bool   `json:"needPassEncrupt"`
	Domain          string `json:"domain"`
}

func (rm *RMSConnectParameter) Handle(r *http.Request, log *LogItem) (*http.Response, error) {
	headers := rm.getHeaders(r)
	if r.URL.Path == "/resto/api/auth" {
		return rm.ProxyGetToken(r, headers, log)
	}
	if r.URL.Path == "/resto/api/logout" {
		return rm.ProxyRestapi(r, headers, log)
	}
			// return rm.ProxyRestapi(r, headers, log)
	val, ok := headers["X-Resto-Authtype"]
	if ok {
		if val[0] == "INTEGRATION" {
			return rm.ProxyIntegrations(r, headers, log)
		} else {
			return rm.ProxyRestapi(r, headers, log)
		}
	}
	fmt.Println("Не известный протокол")
	fmt.Println(r.URL.Path)
	return rm.ProxySimple(r, headers, log)
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

func (rm *RMSConnectParameter) getPassword() string {
	if rm.NeedPassEncrupt {
		hasher := sha1.New()
		hasher.Write([]byte(rm.Password))
		pass := hex.EncodeToString(hasher.Sum(nil))
		return pass
	}
	return rm.Password
}

func (rm *RMSConnectParameter) ProxyGetToken(r *http.Request, headers http.Header, log *LogItem) (*http.Response, error) {
	uri, err := url.Parse(rm.URL)
	if err != nil {
		return nil, err
	}
	uri.Path = r.URL.Path
	query := r.URL.Query()
	query.Set("login", rm.Login)
	query.Set("pass", rm.getPassword())
	uri.RawQuery = query.Encode()
	return rm.Proxy(r, uri, headers, log)
}

func (rm *RMSConnectParameter) GetToken() (*http.Response, error) {
	uri, err := url.Parse(rm.URL)
	if err != nil {
		return nil, err
	}
	uri.Path = "/resto/api/auth"
	query := uri.Query()
	query.Set("login", rm.Login)
	query.Set("pass", rm.getPassword())
	uri.RawQuery = query.Encode()
	resp, err := http.Get(uri.String())
	return resp, err
}

func (rm *RMSConnectParameter) ProxyRestapi(r *http.Request, headers http.Header, log *LogItem) (*http.Response, error) {
	respToken, err := rm.GetToken()
	if r.URL.Path == "/resto/api/auth" {
		return respToken, err
	}
	if err != nil {
		return nil, err
	}
	b, _ := io.ReadAll(respToken.Body)
	token := string(b)
	defer rm.Logout(token)
	uri, err := url.Parse(rm.URL)
	if err != nil {
		return nil, err
	}
	uri.Path = r.URL.Path
	query := r.URL.Query()
	query.Set("key", token)
	uri.RawQuery = query.Encode()
	return rm.Proxy(r, uri, headers, log)
}

func (rm *RMSConnectParameter) ProxyIntegrations(r *http.Request, headers http.Header, log *LogItem) (*http.Response, error) {
	uri, err := url.Parse(rm.URL)
	if err != nil {
		return nil, err
	}
	uri.Path = r.URL.Path
	query := r.URL.Query()
	uri.RawQuery = query.Encode()
	headers["X-Resto-Loginname"] = []string{rm.Login}
	headers["X-Resto-Passwordhash"] = []string{rm.getPassword()}
	return rm.Proxy(r, uri, headers, log)
}

func (rm *RMSConnectParameter) ProxySimple(r *http.Request, headers http.Header, log *LogItem) (*http.Response, error) {
	uri, err := url.Parse(rm.URL)
	if err != nil {
		return nil, err
	}
	uri.Path = r.URL.Path
	query := r.URL.Query()
	uri.RawQuery = query.Encode()
	return rm.Proxy(r, uri, headers, log)
}

func (rm *RMSConnectParameter) Proxy(r *http.Request, uri *url.URL, headers http.Header, log *LogItem) (*http.Response, error) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, uri.String(), bytes.NewBuffer(requestBody))
	req.Header = headers
	log.ClientProxyRequest = CreateHTTPRequest(req)

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
