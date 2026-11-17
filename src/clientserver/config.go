package clientserver

import (
	"net/http"

	"rms_proxy/v2/src/localstore"

	"github.com/gin-gonic/gin"
)

func (cs *ClientServer) GetListRms(c *gin.Context) {
	list := cs.storeConfig.GetRMSList()
	c.JSON(http.StatusOK, list)
}

func (cs *ClientServer) SaveListRms(c *gin.Context) {
	list := localstore.ConfigRmsList{}
	c.ShouldBindJSON(&list)
	cs.storeConfig.SaveRmsList(list)

	c.JSON(http.StatusOK, list)
	cs.restartChanSignal <- true
}

func (cs *ClientServer) GetListProxy(c *gin.Context) {
	list := cs.storeConfig.GetProxyItems()
	c.JSON(http.StatusOK, list)
}

func (cs *ClientServer) SaveListProxy(c *gin.Context) {
	list := []localstore.ConfigReplacedItemInput{}
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := cs.storeConfig.SaveProxyItems(list); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
	cs.restartChanSignal <- true
}

func (cs *ClientServer) GetProxyContent(c *gin.Context) {
	urlPath := c.Query("path")
	content, err := cs.storeConfig.GetProxyContent(urlPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, content)
}
