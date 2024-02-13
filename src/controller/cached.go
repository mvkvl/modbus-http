// region - header

package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mvkvl/modbus-http/service"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	Status(w http.ResponseWriter, r *http.Request)
	Start(w http.ResponseWriter, r *http.Request)
	Stop(w http.ResponseWriter, r *http.Request)
	Cycle(w http.ResponseWriter, r *http.Request)
	Metrics(w http.ResponseWriter, r *http.Request)
	PrometheusMetrics(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Write(w http.ResponseWriter, r *http.Request)
}

func (c *cachedModbusController) Status(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%s\n", c.poller.Status())))
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

func (c *cachedModbusController) Metrics(w http.ResponseWriter, r *http.Request) {
	for _, m := range c.poller.Metrics() {
		w.Write([]byte(fmt.Sprintf("%s\n", m)))
	}
}

func (c *cachedModbusController) PrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	for _, m := range c.poller.Prometheus() {
		w.Write([]byte(fmt.Sprintf("%s\n", m)))
	}
}

func (c *cachedModbusController) Get(w http.ResponseWriter, r *http.Request) {
	result, err := c.poller.Read(getMetricKey(r))
	if nil != err {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
	} else {
		if nil == result {
			w.WriteHeader(204)
			w.Write([]byte(fmt.Sprintf("{}")))
		} else {
			buff, _ := json.Marshal(result)
			w.WriteHeader(200)
			w.Write(buff)
		}
	}
}

func (c *cachedModbusController) Write(w http.ResponseWriter, r *http.Request) {
	b, e := io.ReadAll(r.Body)
	defer r.Body.Close()
	if nil != e {
		w.Write([]byte(fmt.Sprintf("Error: %s", e)))
		return
	}
	bodyStr := string(b)

	ok := false
	v, e := strconv.ParseUint(bodyStr, 10, 64)
	if nil != e {
		bodyStr, e = hexaNumberToInteger(bodyStr)
		if nil != e {
			w.Write([]byte(fmt.Sprintf("Error: %s", e)))
			return
		} else {
			v, e = strconv.ParseUint(bodyStr, 16, 64)
			if nil != e {
				w.Write([]byte(fmt.Sprintf("Error: %s", e)))
				return
			} else {
				ok = true
			}
		}
	} else {
		ok = true
	}
	if ok {
		e = c.poller.WriteWord(getMetricKey(r), uint16(v), nil)
		if nil != e {
			w.Write([]byte(fmt.Sprintf("Error: %s", e)))
			return
		}
	}
	w.Write([]byte("ok"))
}

// endregion

// region - private methods

func getMetricKey(r *http.Request) string {
	return mux.Vars(r)["metric"]
}

func hexaNumberToInteger(hexaString string) (string, error) {
	// replace 0x or 0X with empty String
	if strings.Contains(hexaString, "0x") && 0 == strings.Index(hexaString, "0x") ||
		strings.Contains(hexaString, "0X") && 0 == strings.Index(hexaString, "0X") {
		numberStr := strings.Replace(hexaString, "0x", "", -1)
		numberStr = strings.Replace(numberStr, "0X", "", -1)
		return numberStr, nil
	}
	return "", errors.New("not a hexadecimal string")

}

// endregion
