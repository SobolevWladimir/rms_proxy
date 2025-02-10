package clientserver

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"rms_proxy/v2/src/parameters"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (cs *ClientServer) route(r *gin.Engine) {
	r.GET("/setting/rms", cs.GetListRms)
	r.POST("/setting/rms", cs.SaveListRms)
	r.GET("/setting/proxy", cs.GetListProxy)
	r.POST("/setting/proxy", cs.SaveListProxy)

	r.GET("/ws", func(c *gin.Context) {
		conn, err := cs.upgrader.Upgrade(c.Writer, c.Request, nil)
		fmt.Println("----------connect")
		if err != nil {
			fmt.Println("Ошибка сокета")
			fmt.Println(err.Error())
			return
		}
		defer conn.Close()
		for {
			if len(cs.Messages) > 0 {
				cs.MessagesMutex.Lock()
				data, err := json.Marshal(cs.Messages)

				if err != nil {
					fmt.Println("Ошибка парсина списка")
					fmt.Println(err.Error())
					break
				}
				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("Error: %v", err)
					} else {
						log.Println("Connection closed by client")
					}
					cs.MessagesMutex.Unlock()
					break
				}
				cs.Messages = make([]parameters.LogItem, 0)
				cs.MessagesMutex.Unlock()
			}
		}
	})
	cs.routeFrontFiles(r)
}

func (cs *ClientServer) routeFrontFiles(r *gin.Engine) {
	filesFront := os.Getenv("RMS_FILE_FILES_FRONT") // Переменная для хранения файлов
	if len(filesFront) == 0 {
		slog.Warn("Не указана папка для хранения фронта (RMS_FILE_FILES_FRONT)")
		return
	}
	entries, err := os.ReadDir(filesFront)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			pwd := path.Join(filesFront, e.Name())
			r.Static("/"+e.Name(), pwd)
		} else {
			pwd := path.Join(filesFront, e.Name())
			r.StaticFile("/"+e.Name(), pwd)
		}
	}
	pwd := path.Join(filesFront, "index.html")
	r.StaticFile("", pwd)
}
