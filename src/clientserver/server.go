package clientserver

import (
	"fmt"
	"net/http"
	"os"
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
	storeConfig       *localstore.ConfigStore
	proxyServers      []*proxyserver.ProxyServer
}

func (cs *ClientServer) LisenChan() {
	cs.Messages = []parameters.LogItem{}
	for {
		val, ok := <-cs.readChanLog
		if !ok {
			fmt.Println(val, ok, "<-- loop broke!")
			break // exit break loop
		}
		cs.Messages = append(cs.Messages, val)
	}
}

func (cs *ClientServer) startProxyServers() {
	cs.proxyServers = []*proxyserver.ProxyServer{}
	engineList := cs.storeConfig.GetEngines()
	for _, conf := range engineList {
		engine := conf.GetActiveProxySettings()
		pServer := &proxyserver.ProxyServer{
			Engine:       engine,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			ChanLog:      cs.readChanLog,
		}
		cs.proxyServers = append(cs.proxyServers, pServer)
		go pServer.Start()
	}
}

func (cs *ClientServer) stopProxyServers() {
	for _, server := range cs.proxyServers {
		server.Stop()
	}
}

func (cs *ClientServer) ListenChangeConfiguration() {
	for {
		<-cs.restartChanSignal
		cs.stopProxyServers()
		cs.startProxyServers()
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

	r := gin.Default()
	cs.route(r)

	r.Run(":9090")
}
