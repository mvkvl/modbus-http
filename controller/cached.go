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

type CachedModbusController struct {
	service service.PollerServiceAPI // poller service
}

func NewCachedModbusClient(service service.PollerServiceAPI) CachedModbusAPI {
	return &CachedModbusController{service: service}
}

// endregion

//region - controller API

type CachedModbusAPI interface {
	Start(w http.ResponseWriter, r *http.Request)
	Stop(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

func (c *CachedModbusController) Start(w http.ResponseWriter, r *http.Request) {
	c.service.Start()
	w.Write([]byte(fmt.Sprintf("ok\n")))
}

func (c *CachedModbusController) Stop(w http.ResponseWriter, r *http.Request) {
	c.service.Stop()
	w.Write([]byte(fmt.Sprintf("ok\n")))
}

func (c *CachedModbusController) Get(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.Read(getMetricKey(r))
	if nil != err {
		w.Write([]byte(fmt.Sprintf("%q\n", err)))
	} else {
		if nil == result {
			w.Write([]byte(fmt.Sprintf("none")))
		} else {
			w.Write([]byte(fmt.Sprintf("% x\n", result)))
		}
	}
}

// endregion

// region - private methods

func getMetricKey(r *http.Request) string {
	return mux.Vars(r)["metric"]
}

// endregion
