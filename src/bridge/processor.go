package bridge

import (
	"github.com/mvkvl/modbus"
	"mbridge/model"
	"mbridge/util"
	"strings"
	"sync"
	"time"
)

type ChannelProcessor interface {
	Start()
	Stop()
	Commander() Commander
	Cache() MetricCache
}

type channelProcessorImpl struct {
	channelTitle  string
	commander     Commander
	poller        Poller
	demultiplexer Demultiplexer
	executor      Executor
	cache         MetricCache
	logger        util.Logger
	started       bool
	mutex         sync.Mutex
}

func CreateProcessor(channel *model.Channel, config *model.Config) ChannelProcessor {

	readCmdQueue := make(chan Command)
	writeCmdQueue := make(chan Command)
	modbusCmdQueue := make(chan Command)

	modbusClient := createModbusClient(createModbusHandlerFactory, channel, config)
	cache := CreateMetricCache(config.GetTTL())
	channelTitle := strings.ToLower(channel.Title)
	return &channelProcessorImpl{
		channelTitle:  channelTitle,
		logger:        util.GetLogger("processor-" + channelTitle),
		demultiplexer: CreateDemultiplexer(readCmdQueue, writeCmdQueue, modbusCmdQueue),
		executor:      CreateExecutor(modbusCmdQueue, modbusClient, cache),
		poller:        CreatePoller(readCmdQueue, channel),
		commander:     CreateCommander(writeCmdQueue, channel, config),
		cache:         cache,
	}
}
func createModbusClient(handlerFactory func(connection string, mode model.Mode) modbus.ClientHandler,
	channel *model.Channel, config *model.Config) ModbusClient {
	handler := handlerFactory(channel.Connection, channel.Mode)
	client := modbus.NewClient(handler)
	return NewModbusClient(config, &client)
}
func createModbusHandlerFactory(connection string, mode model.Mode) modbus.ClientHandler {
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

func (p *channelProcessorImpl) Start() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.started {
		return
	}
	p.started = true
	p.logger.Info("start %s processor", p.channelTitle)
	p.demultiplexer.Start(p.channelTitle)
	p.executor.Start(p.channelTitle)
	p.poller.Start(p.channelTitle)
}
func (p *channelProcessorImpl) Stop() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.started {
		return
	}
	p.started = false
	p.logger.Info("stop channel processor %s", p.channelTitle)
	p.poller.Stop(p.channelTitle)
	p.executor.Stop(p.channelTitle)
	p.demultiplexer.Stop(p.channelTitle)
	p.logger.Info("stopped %s processor", p.channelTitle)
}
func (p *channelProcessorImpl) Cache() MetricCache {
	return p.cache
}
func (p *channelProcessorImpl) Commander() Commander {
	return p.commander
}
