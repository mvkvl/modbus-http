package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/jeanphorn/log4go"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/controller"
	"github.com/mvkvl/modbus-http/model"
	"github.com/mvkvl/modbus-http/service"
	"net/http"
	"os"
	"time"
)

func main() {

	channelConfigFilePtr := flag.String("c", "", "channel configuration file")
	loggingConfigFilePtr := flag.String("l", "", "logging configuration file")
	portPtr := flag.Int("p", 8080, "http port to bind to")
	debugFlag := flag.Bool("d", false, "debug output")

	flag.Parse()

	if nil == channelConfigFilePtr || "" == *channelConfigFilePtr {
		fmt.Println("Error: no channel configuration file passed")
		return
	}

	if nil != loggingConfigFilePtr && "" != *loggingConfigFilePtr {
		log.LoadConfiguration(*loggingConfigFilePtr)
	} else {
		log.Close()
	}
	defer log.Close()

	config, err := readConfig(*channelConfigFilePtr)
	if nil != err {
		fmt.Sprintf("Error: could not read config file: %s\n", err)
		return
	}

	if *debugFlag {
		printConfig(config)
	}
	startServer(config, *portPtr)
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
	fmt.Printf("ttl: %d seconds,\nprometheus enabled: %t\nchannels:\n", config.Ttl, config.PrometheusExport)
	for _, c := range config.Channels {
		fmt.Printf("\ttitle: %s, conn: %s, mode: %s, cpause: %d, rpause: %d\n",
			c.Title, c.Connection, c.Mode, c.GetCyclePause(), c.GetRegisterPause())
		for _, d := range c.Devices {
			fmt.Printf("\t\t%s:%d\n", d.Title, d.SlaveId)
			for _, r := range d.Registers {
				fmt.Printf(
					"\t\t\taddr: %4d, size: %2d, type: %7s, mode: %s, factor: %.2f, dev: %s\n",
					r.Address, r.Size, r.Type, r.Mode, r.Factor, r.Device.Title)
			}
		}
	}
}
func startServer(config *model.Config, port int) {

	poller := service.CreateModbusPoller(modbusHandlerFactory, config)
	poller.Start()

	cachedModbusService := controller.NewCachedModbusClient(poller)

	r := mux.NewRouter()
	r.HandleFunc("/status", cachedModbusService.Status).Methods("GET")
	r.HandleFunc("/start", cachedModbusService.Start).Methods("POST")
	r.HandleFunc("/stop", cachedModbusService.Stop).Methods("POST")
	r.HandleFunc("/cycle", cachedModbusService.Cycle).Methods("POST")
	r.HandleFunc("/metrics", cachedModbusService.Metrics).Methods("GET")
	r.HandleFunc("/metric/{metric}", cachedModbusService.Get).Methods("GET")
	r.HandleFunc("/metric/{metric}", cachedModbusService.Write).Methods("POST")

	if config.PrometheusExport {
		r.HandleFunc("/metrics/prometheus", cachedModbusService.PrometheusMetrics).Methods("GET")
	}

	// Bind to a port and pass our router in
	log.Warn(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
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
