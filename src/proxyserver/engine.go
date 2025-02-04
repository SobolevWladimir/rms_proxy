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
	Engine       *parameters.ProxyEngine
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ChanLog      chan parameters.LogItem
	server       *http.Server
}

func (sv *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.String())
	result, log := sv.Engine.Handle(r)
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
}

func (sv *ProxyServer) setHeader(w http.ResponseWriter, resp *http.Response) {
	for key, value := range resp.Header {
		w.Header().Add(key, value[0])
	}
}

func (sv *ProxyServer) Start() {
	fmt.Println("----------------- --------------------------")
	fmt.Println("Start server for: ", sv.Engine.MainRms.Name, "port: ", sv.Engine.Port )
	fmt.Println("--------------------------------------------")
	sv.server = &http.Server{
		Addr:           sv.Engine.Port,
		Handler:        sv,
		ReadTimeout:    sv.ReadTimeout,
		WriteTimeout:   sv.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(sv.server.ListenAndServe())
}

func (sv *ProxyServer) Stop() error {
	fmt.Println("----------------- --------------------------")
	fmt.Println("Stop server : ", sv.Engine.MainRms.Name)
	err := sv.server.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("--------------------------------------------")
	return err
}
