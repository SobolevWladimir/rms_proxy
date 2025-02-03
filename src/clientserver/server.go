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
	upgrader    websocket.Upgrader
	readChanLog chan parameters.LogItem
	Messages    []parameters.LogItem
	storeConfig *localstore.ConfigStore
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

func (cs *ClientServer) StartServer() {
	chanLog := make(chan parameters.LogItem)
	cs.readChanLog = chanLog

	settingPath := os.Getenv("RMS_FILE_SETTING")
	if len(settingPath) == 0 {
		settingPath = "./"
		fmt.Println("Сохраняем котфигурацию в ", settingPath)
		fmt.Println("Для изменения задайте  RMS_FILE_SETTING")
	}

	// Переменная для хранения файлов
	cs.storeConfig = &localstore.ConfigStore{Path: settingPath}

	pServer := &proxyserver.ProxyServer{
		Port:         ":8084",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ChanLog:      chanLog,
	}
	go pServer.Start()
	go cs.LisenChan()

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
