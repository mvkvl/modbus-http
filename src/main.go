package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"mbridge/bridge"
	"mbridge/controller"
	"mbridge/model"
	"mbridge/util"
	"mbridge/util/env"
	"net/http"
	"os"
	"syscall"
)

var defaultConfigFile = "/etc/mbridge/mbridge.properties"
var defaultChannelsFile = "/etc/mbridge/channels.json"

//go:embed logo.txt
var logo string

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("application error: ", r)
		}
	}()

	printLogo()
	configFilePtr := flag.String("c", "", "configuration file")
	flag.Parse()

	if configFilePtr == nil {
		configFilePtr = &defaultConfigFile
	}

	loadEnv(*configFilePtr)

	var config *model.Config

	config = readConfig(env.StringOrDefault("CHANNELS_CONFIG", defaultChannelsFile))
	printConfig(config)

	bridge := bridge.CreateBridge(config)
	defer bridge.Stop()
	bridge.Start()

	go startServer(config, bridge, env.IntOrDefault("SERVICE_PORT", 8080))

	util.GetLogger("main").Info("waiting for break signal...")
	util.HandleSignals(syscall.SIGINT, syscall.SIGTERM)
	util.GetLogger("main").Info("stop program")
}
func printLogo() {
	fmt.Println("")
	fmt.Println(logo)
	fmt.Println("")
}
func loadEnv(configFile string) {
	err := godotenv.Load(configFile)
	if err != nil {
		os.Setenv("GO_ENV", "dev")
		util.GetLogger("main").Debug("could not read configuration file")
	}
}
func readConfig(path string) *model.Config {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config model.Config
	if err := json.Unmarshal(content, &config); err != nil {
		panic(err)
	}
	// add back-reference from register to device to channel
	for i := 0; i < len(config.Channels); i++ {
		c := config.Channels[i]
		for j := 0; j < len(c.Devices); j++ {
			d := c.Devices[j]
			d.Channel = &c
			for k := 0; k < len(d.Registers); k++ {
				d.Registers[k].Device = &d
			}
		}
	}
	return &config
}
func printConfig(config *model.Config) {
	fmt.Printf("ttl: %s,\nprometheus enabled: %t\nchannels:\n", *config.Ttl, config.PrometheusExport)
	for _, c := range config.Channels {
		fmt.Printf("\ttitle: %s, conn: %s, mode: %s, cpause: %d, rpause: %d\n",
			c.Title, c.Connection, c.Mode, c.GetCyclePause(), c.GetRegisterPause())
		for _, d := range c.Devices {
			fmt.Printf("\t\t%s (%s):%d\n", d.Title, d.Alias, d.SlaveId)
			for _, r := range d.Registers {
				fmt.Printf(
					"\t\t\taddr: %4d, size: %2d, type: %7s, mode: %s, factor: %.2f, dev: %s\n",
					r.Address, r.Size, r.Type, r.Mode, r.Factor, r.Device.Title)
			}
		}
	}
	fmt.Println()
}
func startServer(config *model.Config, bridge bridge.Bridge, port int) {

	controller := controller.NewBridgeController(bridge)

	r := mux.NewRouter()
	r.HandleFunc("/start", controller.Start).Methods("POST")
	r.HandleFunc("/stop", controller.Stop).Methods("POST")
	r.HandleFunc("/registers", controller.Registers).Methods("GET")
	r.HandleFunc("/metrics", controller.Metrics).Methods("GET")
	r.HandleFunc("/flush", controller.Flush).Methods("POST")
	r.HandleFunc("/metric/{metric}", controller.Get).Methods("GET")
	r.HandleFunc("/metric/{metric}", controller.Write).Methods("POST")

	if config.PrometheusExport {
		r.HandleFunc("/metrics/prometheus", controller.PrometheusMetrics).Methods("GET")
	}

	// Bind to a port and pass our router in
	util.GetLogger("main").Error("%v", http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
