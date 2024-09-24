package parameters

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RMSConnectParameter struct {
	url      string
	login    string
	password string
}

func (rm *RMSConnectParameter) Handle(r *http.Request) (string, error) {
	fmt.Println(r.URL.String())
	fmt.Println(r.Header)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("BODY:")
	fmt.Printf("%s", b)
	fmt.Println("--------")
	fmt.Println(r.URL.Path)

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
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, uri.String(), r.Body)
	// req.Header = r.Header
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
