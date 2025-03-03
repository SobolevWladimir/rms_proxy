package parameters

import (
	"io"
	"net/http"
	"strings"
)

// Сам лог для записи
type LogItem struct {
	ClientRequest      *HTTPRequst          `json:"clientRequest"`
	ClientProxyRequest *HTTPRequst          `json:"clientProxyRequest"`
	IsProxy            bool                 `json:"isProxy"` // Если проксируется
	MainRMS            *RMSConnectParameter `json:"mainRms"`
	ProxyTo            *ReplacedItem        `json:"proxyTo"`         // Куда проксируется
	IsErrorResponse    bool                 `json:"isErrorResponse"` //  Если клиент ответил с ошибокй
	ClientResponse     *HTTPResponse        `json:"clientResponse"`
	ErrorResponse      string               `json:"error"`
}

type HTTPRequst struct {
	Method string
	URL    string
	Header http.Header
	Body   string
}

func CreateHTTPRequest(r *http.Request) *HTTPRequst {
	b, _ := io.ReadAll(r.Body)
	// io.ReadAll(r io.Reader)
	body := string(b)
	r.Body = io.NopCloser(strings.NewReader(body))

	return &HTTPRequst{
		Header: r.Header,
		Body:   body,
		URL:    r.URL.String(),
		Method: r.Method,
	}
}


type HTTPResponse struct {
	// Data *http.Response
	Body   string
	Header http.Header
	Status int
}

func CreateHTTPResponse(r *http.Response) *HTTPResponse {
	b, _ := io.ReadAll(r.Body)
	body := string(b)
	r.Body = io.NopCloser(strings.NewReader(body))

	return &HTTPResponse{
		Header: r.Header,
		Body:   body,
		Status: r.StatusCode,
		// Body: "",
	}
}
