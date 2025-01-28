package proxyserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"rms_proxy/v2/src/parameters"
)

type ProxyServer struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ChanLog      chan parameters.LogItem
}

func (sv *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n----------------")
	fmt.Println(r.Method, r.URL.String())
	repository := parameters.SettingsRepositoryMemory{}
	engine := repository.GetActiveProxySettings()
	fmt.Println("result:")
	result, log := engine.Handle(r)
	sv.ChanLog <- log
	if log.IsErrorResponse {
		fmt.Println("error response", log.ErrorResponse)
		fmt.Fprintln(w, log.ErrorResponse)
	}
	sv.setHeader(w, result)
	w.WriteHeader(result.StatusCode)

	fmt.Println("status code", result.StatusCode)
	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println("error  read body", err.Error())
		fmt.Fprintln(w, err.Error())
	}
	fmt.Fprintln(w, string(body))
	fmt.Println("---------------- end ---------------------")
}

func (sv *ProxyServer) setHeader(w http.ResponseWriter, resp *http.Response) {
	for key, value := range resp.Header {
		w.Header().Add(key, value[0])
	}
}

func (sv *ProxyServer) Start() {
	fmt.Println("start server")
	s := &http.Server{
		Addr:           sv.Port,
		Handler:        sv,
		ReadTimeout:    sv.ReadTimeout,
		WriteTimeout:   sv.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
