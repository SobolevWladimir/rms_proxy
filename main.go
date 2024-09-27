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
const settinsRmsPage ="/settings/rms.json"
const settinsRoutePage ="/settings/route.json"
// test
func (ct *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if(r.URL.Path == settinsRmsPage || r.URL.Path == settinsRoutePage) {
		ct.ServeHTTPSettings(w, r);
		return; 
	}

	ct.ServeHTTPProxy(w, r);
}

func (ct *CounterHandler) ServeHTTPSettings(w http.ResponseWriter, r *http.Request) {
}
func (ct *CounterHandler) ServeHTTPProxy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n----------------")
	fmt.Println(r.Method, r.URL.String())
	repository := parameters.SettingsRepositoryMemory{}
	engine := repository.GetActiveProxySettings()
	fmt.Println("result:")
	result, err := engine.Handle(r)
	if err != nil {
		fmt.Println("error response", err.Error())
		fmt.Fprintln(w, err.Error())
	}
	ct.setHeader(w, result)
	w.WriteHeader(result.StatusCode)

	fmt.Println("status code", result.StatusCode)
	body, err := io.ReadAll(result.Body)
	fmt.Println("body:")
	fmt.Println(string(body))
	if err != nil {
		fmt.Println("error  read body", err.Error())
		fmt.Fprintln(w, err.Error())
	}
	fmt.Fprintln(w, string(body))
	fmt.Println("---------------- end ---------------------")
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
