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
	url      string
	login    string
	password string
}

func (rm *RMSConnectParameter) Handle(r *http.Request) (string, error) {
	token, err := rm.GetToken()
	if r.URL.Path == "/resto/api/auth" {
		return token, err
	}
	if err != nil {
		return "", err
	}
	return rm.Proxy(r, token)
}

func (rm *RMSConnectParameter) GetToken() (string, error) {
	uri, err := url.Parse(rm.url)
	if err != nil {
		return "", err
	}
	uri.Path = "/resto/api/auth"
	hasher := sha1.New()
	hasher.Write([]byte(rm.password))
	query := uri.Query()
	query.Set("login", rm.login)
	pass := hex.EncodeToString(hasher.Sum(nil))
	query.Set("pass", pass)
	uri.RawQuery = query.Encode()
	resp, err := http.Get(uri.String())
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)

	return string(b), err
}

func (rm *RMSConnectParameter) Proxy(r *http.Request, token string) (string, error) {
	uri, err := url.Parse(rm.url)
	if err != nil {
		return "", err
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
	fmt.Println("URL:", uri.String())
	fmt.Println("METHOD:", r.Method)
	fmt.Println("BODY:", string(requestBody))
	req, err := http.NewRequest(r.Method, uri.String(),bytes.NewBuffer(requestBody))
	req.Header = rm.getHeaders(r)
	fmt.Println("Headers:", req.Header)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	b, err := io.ReadAll(resp.Body)


	if err != nil {
		fmt.Println("server response error:")
		fmt.Println(err.Error())

	} else {
		fmt.Println("server response:")
		fmt.Println(string(b))

	}

	return string(b), err

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
