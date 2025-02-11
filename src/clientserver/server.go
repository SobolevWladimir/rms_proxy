package clientserver

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"rms_proxy/v2/src/localstore"
	"rms_proxy/v2/src/parameters"
	"rms_proxy/v2/src/proxyserver"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ClientServer struct {
	// ReadChanLog chan parameters.LogItem
	upgrader          websocket.Upgrader
	readChanLog       chan parameters.LogItem
	restartChanSignal chan bool
	Messages          []parameters.LogItem
	MessagesMutex     *sync.Mutex
	storeConfig       *localstore.ConfigStore
	server            *proxyserver.ProxyServer
}

func (cs *ClientServer) LisenChan() {
	cs.Messages = []parameters.LogItem{}
	cs.MessagesMutex = new(sync.Mutex)
	for {
		val, ok := <-cs.readChanLog
		cs.MessagesMutex.Lock()
		if !ok {
			fmt.Println("ошибка чтения с канала")
			fmt.Println(val, ok, "<-- loop broke!")
			break // exit break loop
		}

		fmt.Println("Запись в канал ", val.ClientRequest.URL)

		cs.Messages = append(cs.Messages, val)
		cs.MessagesMutex.Unlock()
	}
}

func (cs *ClientServer) startProxyServers() {
	engine := cs.storeConfig.GetActiveProxySettings()
	pServer := &proxyserver.ProxyServer{
		Engine:       engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ChanLog:      cs.readChanLog,
	}
	cs.server = pServer
	go pServer.Start()
}

func (cs *ClientServer) ListenChangeConfiguration() {
	for {
		<-cs.restartChanSignal
		engine := cs.storeConfig.GetActiveProxySettings()
		cs.server.Engine = engine
	}
}

func (cs *ClientServer) StartServer() {
	chanLog := make(chan parameters.LogItem)
	cs.readChanLog = chanLog
	cs.restartChanSignal = make(chan bool)

	settingPath := os.Getenv("RMS_FILE_SETTING")
	if len(settingPath) == 0 {
		settingPath = "./"
		fmt.Println("Сохраняем котфигурацию в ", settingPath)
		fmt.Println("Для изменения задайте  RMS_FILE_SETTING")
	}

	// Переменная для хранения файлов
	cs.storeConfig = &localstore.ConfigStore{Path: settingPath}
	cs.startProxyServers()

	go cs.LisenChan()
	go cs.ListenChangeConfiguration()

	cs.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	cs.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	time.Sleep(1 * time.Second)
	r := gin.Default()
	cs.route(r)

	err := r.Run(":9090")
	if err != nil {
		fmt.Println("server error:")
		fmt.Println(err.Error())

	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
}
