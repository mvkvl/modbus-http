package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/jeanphorn/log4go"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/controller"
	"github.com/mvkvl/modbus-http/model"
	"github.com/mvkvl/modbus-http/service"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {

	log.LoadConfiguration("./conf/logger.json")
	defer log.Close()

	config, err := readConfig("./conf/channels.json")
	if nil != err {
		log.Info("error reading config: %s\n", err)
		return
	}

	printConfig(config)
	startServer(config)
	//schedulerTestAfter()
	//schedulerTestFixedRate()
	//schedulerTestFixedDelay()
}

func readConfig(path string) (*model.Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config model.Config
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}
	// add back-reference from register to device
	for i := 0; i < len(config.Channels); i++ {
		c := config.Channels[i]
		for j := 0; j < len(c.Devices); j++ {
			d := c.Devices[j]
			for k := 0; k < len(d.Registers); k++ {
				d.Registers[k].Device = &d
			}
		}
	}
	return &config, nil
}
func printConfig(config *model.Config) {
	for _, c := range config.Channels {
		log.Info("title: %s, conn: %s, mode: %s, cpause: %d, rpause: %d",
			c.Title, c.Connection, c.Mode, c.GetCyclePause(), c.GetRegisterPause())
		for _, d := range c.Devices {
			log.Info("\t%s:%d", d.Title, d.SlaveId)
			for _, r := range d.Registers {
				log.Info(
					"\t\taddr: %4d, size: %2d, type: %7s, mode: %s, factor: %f, dev: %s",
					r.Address, r.Size, r.Type, r.Mode, r.Factor, r.Device.Title)
			}
		}
	}
}
func startServer(config *model.Config) {

	poller := service.CreateModbusPoller(modbusHandlerFactory, config)
	poller.Start()

	cachedModbusService := controller.NewCachedModbusClient(poller)

	r := mux.NewRouter()
	r.HandleFunc("/start", cachedModbusService.Start).Methods("POST")
	r.HandleFunc("/stop", cachedModbusService.Stop).Methods("POST")
	r.HandleFunc("/cycle", cachedModbusService.Cycle).Methods("POST")
	r.HandleFunc("/metrics", cachedModbusService.Metrics).Methods("GET")
	r.HandleFunc("/metric/{metric}", cachedModbusService.Get).Methods("GET")
	r.HandleFunc("/metric/{metric}", cachedModbusService.Write).Methods("POST")

	// Bind to a port and pass our router in
	log.Warn(http.ListenAndServe(":8080", r))
}

func schedulerTestAfter() {
	var wg = &sync.WaitGroup{}
	s := service.NewScheduler()
	wg.Add(1)

	s.RunAfter(func() {
		blockedScheduledPayloadTest(service.RandomString(10), wg)
		wg.Done()
	}, 1*time.Second)

	wg.Wait()
}
func schedulerTestFixedRate() {
	var wg = &sync.WaitGroup{}
	s := service.NewScheduler()
	time.AfterFunc(7*time.Second, s.Stop)
	s.RunAtFixedRate(func() { blockedScheduledPayloadTest(service.RandomString(10), wg) }, 1*time.Second)
	time.Sleep(1 * time.Second)
	s.RunAfter(func() { blockedScheduledPayloadTest(service.RandomString(10), wg) }, 1*time.Second)
	wg.Wait()
}
func schedulerTestFixedDelay() {
	var wg = &sync.WaitGroup{}
	s := service.NewScheduler()
	wg.Add(1)
	//time.AfterFunc(16*time.Second, s.Stop)
	s.RunWithFixedDelay(func() { blockedScheduledPayloadTest(service.RandomString(10), wg) }, 1*time.Second)
	wg.Wait()
}

func blockedScheduledPayloadTest(id string, wg *sync.WaitGroup) {
	wg.Add(1)
	scheduledPayloadTest(id)
	wg.Done()
}
func scheduledPayloadTest(id string) {
	log.Info(fmt.Sprintf("%s: started", id))
	time.Sleep(5 * time.Second)
	log.Info(fmt.Sprintf("%s: finished", id))
}

func modbusHandlerFactory(connection string, mode model.Mode) modbus.ClientHandler {
	switch mode {
	case model.ENC:
		_handler := modbus.NewEncClientHandler(connection)
		_handler.IdleTimeout = 2 * time.Second
		_handler.Timeout = 1 * time.Second
		//_handler.Logger = log.New(os.Stdout, fmt.Sprintf("[%s]: ", connection), log.LstdFlags|log.Lmicroseconds)
		return _handler
	case model.TCP:
		_handler := modbus.NewTCPClientHandler(connection)
		//_handler.Logger = log.New(os.Stdout, fmt.Sprintf("[%s]: ", connection), log.LstdFlags|log.Lmicroseconds)
		return _handler
	case model.RTU:
		_handler := modbus.NewRTUClientHandler(connection)
		//_handler.Logger = log.New(os.Stdout, fmt.Sprintf("[%s]: ", connection), log.LstdFlags|log.Lmicroseconds)
		return _handler
	}
	return nil
}
