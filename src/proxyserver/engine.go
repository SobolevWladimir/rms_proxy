package proxyserver

import (
	"context"
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

	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println("error  read body", err.Error())
		fmt.Fprintln(w, err.Error())
	}
	w.Write(body)
}

func (sv *ProxyServer) setHeader(w http.ResponseWriter, resp *http.Response) {
	for key, value := range resp.Header {
		w.Header().Add(key, value[0])
	}
}

func (sv *ProxyServer) Start() {
	server := new(http.Server)
	server.Addr = sv.Engine.Port
	server.Handler = sv

	server.ReadTimeout = sv.ReadTimeout
	server.WriteTimeout = sv.WriteTimeout
	server.MaxHeaderBytes = 1 << 20

	sv.server = server

	log.Fatal(sv.server.ListenAndServe())
}

func (sv *ProxyServer) Stop() error {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	err := sv.server.Shutdown(shutdownCtx)
	// err := sv.server.Close()
	if err != nil {
		fmt.Println(" ------ server error")
		fmt.Println(err.Error())
	}
	return err
}
