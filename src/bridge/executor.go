package bridge

import (
	"mbridge/model"
	"mbridge/util"
	"sync"
	"time"
)

// per-channel command executor

type Executor interface {
	Start(title string)
	Stop(title string)
}

type executorImpl struct {
	modbusChn    <-chan Command
	quitChn      chan struct{}
	logger       util.Logger
	modbusClient ModbusClient
	cache        MetricCache
	started      bool
	mutex        sync.Mutex
}

func CreateExecutor(modbusChn chan Command, modbusClient ModbusClient, cache MetricCache) Executor {
	return &executorImpl{
		modbusChn:    modbusChn,
		logger:       util.GetLogger("executor"),
		modbusClient: modbusClient,
		cache:        cache,
	}
}

func (e *executorImpl) Start(title string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.started {
		return
	}
	e.started = true
	e.quitChn = make(chan struct{})

	go func() {
		e.logger.Info("start executor %s", title)
		defer func() {
			close(e.quitChn)
			e.logger.Info("shutdown executor %s", title)
		}()
		for {
			select {
			case cmd, ok := <-e.modbusChn:
				if ok {
					e.handleCommand(cmd)
				}
			case <-e.quitChn:
				return
				//default:
				//	time.Sleep(delay)
			}
		}
	}()
}
func (e *executorImpl) Stop(title string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if !e.started {
		return
	}
	e.started = false
	e.logger.Info("stop executor %s", title)
	e.quitChn <- struct{}{}
}

func (e *executorImpl) handleCommand(cmd Command) {
	switch cmd.GetType() {
	case CTRead:
		e.readRegister(cmd)
	case CTWrite:
		e.writeRegister(cmd)
	}
}

func (e *executorImpl) readRegister(cmd Command) {
	raw, val, err := e.modbusClient.Read(cmd.GetRegister())
	if err != nil {
		e.logger.Warning("read error: %v", err)
	} else {
		e.logger.Trace("%v : %v : %s", raw, val, model.MetricKey(cmd.GetRegister()))
		e.cache.Set(e.cache.Key(cmd.GetChannel(), cmd.GetRegister()), &model.Metric{
			Key:       model.MetricKey(cmd.GetRegister()),
			Channel:   cmd.GetRegister().Device.Channel.Title,
			Device:    cmd.GetRegister().Device.Title,
			Alias:     cmd.GetRegister().Device.Alias,
			Register:  cmd.GetRegister().Title,
			RawValue:  raw,
			Value:     val,
			Timestamp: time.Now(),
		})
	}
}
func (e *executorImpl) writeRegister(cmd Command) {
	err := e.modbusClient.Write(cmd.GetRegister(), uint16(cmd.GetValue()))
	if err != nil {
		e.logger.Warning("write error: %v", err)
	}
}
