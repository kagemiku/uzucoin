package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Server Server
}

func (this *Handler) AddTransaction(c *gin.Context) {
	t := &Transaction{}
	if err := c.BindJSON(t); err != nil {
		log.Println("Failed to bind json in AddTransaction")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := this.Server.AddTransaction(t)
	if err != nil {
		log.Println("Failed to get res in AddTransaction")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (this *Handler) GetTask(c *gin.Context) {
	res, err := this.Server.GetTask(&GetTaskRequest{})
	if err != nil {
		log.Println("Failed to bind json in GetTask")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (this *Handler) ResolveNonce(c *gin.Context) {
	n := &Nonce{}
	if err := c.BindJSON(n); err != nil {
		log.Println("Failed to bind json in ResolveNonce")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := this.Server.ResolveNonce(n)
	if err != nil {
		log.Println("Failed to bind json in ResolveNonce")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func main() {
	server := &ServerImpl{
		Idles:            []*Idle{},
		TransactionQueue: []*FixedTransaction{},
	}
	handler := &Handler{Server: server}

	r := gin.Default()
	r.POST("/add_transaction", handler.AddTransaction)
	r.GET("/task", handler.GetTask)
	r.POST("/resolve", handler.ResolveNonce)
	r.Run()
}
