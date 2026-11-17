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
	"sync"
	"time"
)

const tokenTTL = 14 * time.Minute

type cachedToken struct {
	token     string
	expiresAt time.Time
}

var (
	tokenCacheMu sync.Mutex
	tokenCache   = map[string]cachedToken{}
)

type RMSConnectParameter struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Login           string `json:"login"`
	Password        string `json:"-"`
	NeedPassEncrupt bool   `json:"needPassEncrupt"`
	Domain          string `json:"domain"`
}

func (rm *RMSConnectParameter) tokenCacheKey() string {
	return rm.URL + "|" + rm.Login
}

func (rm *RMSConnectParameter) getCachedToken() (string, bool) {
	tokenCacheMu.Lock()
	defer tokenCacheMu.Unlock()
	c, ok := tokenCache[rm.tokenCacheKey()]
	if !ok || time.Now().After(c.expiresAt) {
		return "", false
	}
	return c.token, true
}

func (rm *RMSConnectParameter) setCachedToken(token string) {
	tokenCacheMu.Lock()
	defer tokenCacheMu.Unlock()
	tokenCache[rm.tokenCacheKey()] = cachedToken{token: token, expiresAt: time.Now().Add(tokenTTL)}
}

func (rm *RMSConnectParameter) clearCachedToken() {
	tokenCacheMu.Lock()
	defer tokenCacheMu.Unlock()
	delete(tokenCache, rm.tokenCacheKey())
}

func (rm *RMSConnectParameter) Handle(r *http.Request, log *LogItem) (*http.Response, error) {
	headers := rm.getHeaders(r)
	if r.URL.Path == "/resto/api/auth" {
		return rm.ProxyGetToken(r, headers, log)
	}
	if r.URL.Path == "/resto/api/logout" {
		rm.clearCachedToken()
		return rm.ProxySimple(r, headers, log)
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
	if r.URL.Query().Has("key") {
			return rm.ProxyRestapi(r, headers, log)

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
	query.Set("key", token)
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

func (rm *RMSConnectParameter) getToken() (string, error) {
	if token, ok := rm.getCachedToken(); ok {
		return token, nil
	}
	resp, err := rm.GetToken()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	token := string(b)
	rm.setCachedToken(token)
	return token, nil
}

func (rm *RMSConnectParameter) ProxyRestapi(r *http.Request, headers http.Header, log *LogItem) (*http.Response, error) {
	token, err := rm.getToken()
	if err != nil {
		return nil, err
	}
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
