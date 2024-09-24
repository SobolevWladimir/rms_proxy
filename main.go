package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"rms_proxy/v2/src/parameters"
	"time"
)

type CounterHandler struct {
	counter int
}

func (ct *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	repository := parameters.SettingsRepositoryMemory{}
	engine := repository.GetActiveProxySettings()
	result, err := engine.Handle(r)
	if err != nil {
		fmt.Println("error response", err.Error())
		fmt.Fprintln(w, err.Error())
	}
	ct.setHeader(w, result)
	w.WriteHeader(result.StatusCode)
	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println("error  read body", err.Error())
		fmt.Fprintln(w, err.Error())
	}
	fmt.Fprintln(w, string(body))
}
func (ct *CounterHandler) setHeader(w http.ResponseWriter, resp *http.Response) {
	for key, value := range resp.Header {
			w.Header().Add(key, value[0])
	}
}

func main() {
	th := &CounterHandler{counter: 0}
	s := &http.Server{
		Addr:           ":8084",
		Handler:        th,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
	fmt.Println("server start")
}
