package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	rpcclient "github.com/james-ray/unielon-indexer/package/github.com/HcashOrg/hcrpcclient"
	"github.com/james-ray/unielon-indexer/storage"
	"github.com/james-ray/unielon-indexer/utils"
	"github.com/james-ray/unielon-indexer/verifys"
)

type Router struct {
	dbc        *storage.DBClient
	node       *rpcclient.Client
	verify     *verifys.Verifys
	feeAddress string
}

func NewRouter(dbc *storage.DBClient, node *rpcclient.Client, feeAddress string) *Router {
	return &Router{
		dbc:        dbc,
		node:       node,
		verify:     verifys.NewVerifys(dbc, feeAddress),
		feeAddress: feeAddress,
	}
}

func (r *Router) LastNumber(c *gin.Context) {
	last, err := r.dbc.LastBlock()
	if err != nil {
		result := &utils.HttpResult{}
		result.Code = 500
		result.Msg = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}

	result := &utils.HttpResult{}
	result.Code = 200
	result.Msg = "success"
	result.Data = last
	c.JSON(http.StatusOK, result)
}
