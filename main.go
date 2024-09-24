package main

import (
	"fmt"
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
	fmt.Fprintln(w, result)
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
