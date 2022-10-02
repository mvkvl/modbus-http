// region - header

package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mvkvl/modbus-http/service"
	"net/http"
)

// endregion

// region - service

type cachedModbusController struct {
	poller service.Poller // poller service
}

func NewCachedModbusClient(poller service.Poller) CachedModbusAPI {
	return &cachedModbusController{poller: poller}
}

// endregion

//region - controller API

type CachedModbusAPI interface {
	Start(w http.ResponseWriter, r *http.Request)
	Stop(w http.ResponseWriter, r *http.Request)
	Cycle(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Metrics(w http.ResponseWriter, r *http.Request)
}

func (c *cachedModbusController) Start(w http.ResponseWriter, r *http.Request) {
	c.poller.Start()
	w.Write([]byte(fmt.Sprintf("ok\n")))
}

func (c *cachedModbusController) Cycle(w http.ResponseWriter, r *http.Request) {
	c.poller.Cycle()
	w.Write([]byte(fmt.Sprintf("ok\n")))
}

func (c *cachedModbusController) Stop(w http.ResponseWriter, r *http.Request) {
	c.poller.Stop()
	w.Write([]byte(fmt.Sprintf("ok\n")))
}

func (c *cachedModbusController) Get(w http.ResponseWriter, r *http.Request) {
	result, err := c.poller.Read(getMetricKey(r))
	if nil != err {
		w.Write([]byte(fmt.Sprintf("%q\n", err)))
	} else {
		if nil == result {
			w.Write([]byte(fmt.Sprintf("none")))
		} else {
			w.Write([]byte(fmt.Sprintf("%.2f\n", result)))
		}
	}
}

func (c *cachedModbusController) Metrics(w http.ResponseWriter, r *http.Request) {
	for _, m := range c.poller.Metrics() {
		w.Write([]byte(fmt.Sprintf("%s\n", m)))
	}
}

// endregion

// region - private methods

func getMetricKey(r *http.Request) string {
	return mux.Vars(r)["metric"]
}

// endregion
