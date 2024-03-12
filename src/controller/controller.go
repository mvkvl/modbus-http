package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"mbridge/bridge"
	"mbridge/model"
	"mbridge/util"
	"net/http"
	"strconv"
)

type ModbusBridgeController interface {
	Start(w http.ResponseWriter, r *http.Request)
	Stop(w http.ResponseWriter, r *http.Request)
	Registers(w http.ResponseWriter, r *http.Request)
	Metrics(w http.ResponseWriter, r *http.Request)
	PrometheusMetrics(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Write(w http.ResponseWriter, r *http.Request)
	Flush(w http.ResponseWriter, r *http.Request)
}

type modbusBridgeControllerImpl struct {
	bridge bridge.Bridge
}

func NewBridgeController(bridge bridge.Bridge) ModbusBridgeController {
	return &modbusBridgeControllerImpl{
		bridge: bridge,
	}
}

func (c *modbusBridgeControllerImpl) Registers(w http.ResponseWriter, r *http.Request) {
	header := "%-30s %-8s %-5s %-5s %-7s %-7s"
	out := fmt.Sprintf(header, "reference", "type", "mode", "size", "addr", "factor")
	w.Write([]byte(fmt.Sprintf("%s\n", out)))

	template := "%-30s %-8s %-5s %-5d %-7d %-5.2f"
	for _, r := range c.bridge.Regs() {
		out = fmt.Sprintf(template, model.MetricKey(r), r.Type, r.Mode, r.Size, r.Address, r.Factor)
		w.Write([]byte(fmt.Sprintf("%s\n", out)))
	}
}
func (c *modbusBridgeControllerImpl) Metrics(w http.ResponseWriter, r *http.Request) {
	for _, m := range c.bridge.List() {
		w.Write([]byte(fmt.Sprintf("%s\n", m)))
	}
}
func (c *modbusBridgeControllerImpl) PrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	template := "modbus_metric_%s{channel=\"%s\",device=\"%s\",alias=\"%s\",register=\"%s\"} %s"
	for _, m := range c.bridge.List() {
		raw := fmt.Sprintf(fmt.Sprintf(template, "raw", m.Channel, m.Device, m.Alias, m.Register, "%d"), m.RawValue)
		w.Write([]byte(fmt.Sprintf("%s\n", raw)))
		val := fmt.Sprintf(fmt.Sprintf(template, "value", m.Channel, m.Device, m.Alias, m.Register, "%f"), m.Value)
		w.Write([]byte(fmt.Sprintf("%s\n", val)))
		ts := fmt.Sprintf(fmt.Sprintf(template, "timestamp", m.Channel, m.Device, m.Alias, m.Register, "%d"), m.Timestamp.Unix())
		w.Write([]byte(fmt.Sprintf("%s\n", ts)))
	}
}
func (c *modbusBridgeControllerImpl) Get(w http.ResponseWriter, r *http.Request) {
	result, err := c.bridge.Get(getMetricKey(r))
	if nil != err {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
	} else {
		if nil == result {
			w.WriteHeader(204)
			w.Write([]byte(fmt.Sprintf("{}")))
		} else {
			buff, _ := json.Marshal(*result)
			w.WriteHeader(200)
			w.Write(buff)
		}
	}
}
func (c *modbusBridgeControllerImpl) Write(w http.ResponseWriter, r *http.Request) {
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
		bodyStr, e = util.HexaNumberToInteger(bodyStr)
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
		e = c.bridge.Set(getMetricKey(r), uint16(v))
		if nil != e {
			w.Write([]byte(fmt.Sprintf("Error: %s", e)))
			return
		}
	}
	w.Write([]byte("ok"))
}
func (c *modbusBridgeControllerImpl) Flush(w http.ResponseWriter, r *http.Request) {
	c.bridge.Flush()
	w.Write([]byte(fmt.Sprintf("ok\n")))
}
func (c *modbusBridgeControllerImpl) Start(w http.ResponseWriter, r *http.Request) {
	c.bridge.Start()
	w.Write([]byte(fmt.Sprintf("ok\n")))
}
func (c *modbusBridgeControllerImpl) Stop(w http.ResponseWriter, r *http.Request) {
	c.bridge.Stop()
	w.Write([]byte(fmt.Sprintf("ok\n")))
}

func getMetricKey(r *http.Request) string {
	return mux.Vars(r)["metric"]
}
