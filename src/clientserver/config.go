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
}
