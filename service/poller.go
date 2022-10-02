// region - header

package service

import (
	"errors"
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/model"
	"github.com/mvkvl/modbus-http/queue"
	"github.com/procyon-projects/chrono"
	"sync"
	"time"
)

// endregion

// region - links

// modbus: https://github.com/mvkvl/modbus
// chrono: https://github.com/procyon-projects/chrono

// endregion

// region - poller service

type PollerServiceAPI interface {
	Start()
	Stop()
	Read(key string) (any, error)
	WriteByte(key string, value uint8, callback func()) error
	WriteWord(key string, value uint16, callback func()) error
	WriteValue(key string, value float32, callback func()) error
	destroy()
}

type PollerService struct {
	clients   map[string]*modbus.Client // channels' modbus client
	config    *model.Config             // channels config
	m         sync.Mutex                // mutex for concurrent access synchronization
	dc        map[string]any            // data cache
	cq        queue.FifoQueue           // commands queue
	scheduler chrono.TaskScheduler
	poller    chrono.ScheduledTask // poller task
}

func CreateModbusPoller(handlerFactory func(connection string, mode model.Mode) modbus.ClientHandler, config *model.Config) PollerServiceAPI {
	var clients = make(map[string]*modbus.Client)
	for _, chn := range config.Channels {
		handler := handlerFactory(chn.Connection, chn.Mode)
		client := modbus.NewClient(handler)
		clients[chn.Title] = &client
	}
	return &PollerService{
		clients: clients,
		config:  config,
		dc:      make(map[string]any),
		cq:      queue.CreateQueue(),
	}
}

func DestroyModbusPoller(poller *PollerServiceAPI) {
	(*poller).destroy()
}

// endregion

// region - poller API

func (s *PollerService) Start() {
	//if nil == s.scheduler || s.scheduler.IsShutdown() || nil == s.poller || s.poller.IsCancelled() {
	//	s.scheduler = chrono.NewDefaultTaskScheduler()
	//	//s.poller, _ = s.scheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
	//	s.poller, _ = s.scheduler.ScheduleAtFixedRate(func(ctx context.Context) {
	//		s.cycle()
	//	}, 10000*time.Millisecond)
	//}
}
func (s *PollerService) Stop() {
	//if nil != s.poller && !s.poller.IsCancelled() {
	//	s.poller.Cancel()
	//	s.poller = nil
	//}
	//if nil != s.scheduler && !s.scheduler.IsShutdown() {
	//	s.scheduler.Shutdown()
	//	s.scheduler = nil
	//}
}
func (s *PollerService) Read(key string) (any, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if val, ok := s.dc[key]; ok {
		return val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no metric found for key '%s'", key))
	}
}
func (s *PollerService) WriteByte(key string, value uint8, callback func()) error {
	//TODO implement me
	panic("implement me")
}
func (s *PollerService) WriteWord(key string, value uint16, callback func()) error {
	//TODO implement me
	panic("implement me")
}
func (s *PollerService) WriteValue(key string, value float32, callback func()) error {
	//TODO implement me
	panic("implement me")
}

// endregion

// region - private methods

func (s *PollerService) destroy() {
}

func (s *PollerService) cycle() {
	var wg = &sync.WaitGroup{}
	for _, chn := range s.config.Channels {
		wg.Add(1)
		go func(chn model.Channel) {
			if err := s.readChannel(chn); err != nil {
				//time.Sleep(time.Duration(chn.CyclePause) * 5 * time.Millisecond)
			}
			//if err := s.writeChannel(chn); err != nil {
			//}
		}(chn)
	}
	wg.Wait()
}

func (s *PollerService) readChannel(chn model.Channel) error {
	log.Info("~~~~~~~~ read channel '%s' [cp: %d, rp: %d]", chn.Title, chn.CyclePause, chn.RegisterPause)
	client := s.clients[chn.Title]
	for _, dev := range chn.Devices {
		for _, reg := range dev.Registers {
			if reg.Mode == model.RO || reg.Mode == model.RW {
				l := fmt.Sprintf("read register: %s:%s:%s ->", chn.Title, dev.Title, reg.Title)
				var result []uint8
				var err error
				switch reg.Type {
				case model.COIL:
					result, err = (*client).ReadCoils(dev.SlaveId, reg.Address, reg.Size)
				case model.DISCRETE:
					result, err = (*client).ReadDiscreteInputs(dev.SlaveId, reg.Address, reg.Size)
				case model.INPUT:
					result, err = (*client).ReadInputRegisters(dev.SlaveId, reg.Address, reg.Size)
				case model.HOLDING:
					result, err = (*client).ReadHoldingRegisters(dev.SlaveId, reg.Address, reg.Size)
				}
				time.Sleep(time.Duration(chn.RegisterPause) * time.Millisecond)
				if nil == err {
					log.Info("%s % x", l, result)
					s.dc[cacheKey(chn.Title, dev.Title, reg.Title)] = result
				} else {
					log.Info("%s Error: %s", l, err)
					//return err
				}
			}
		}
	}
	time.Sleep(time.Duration(chn.CyclePause) * time.Millisecond)
	return nil
}

func cacheKey(chn, dev, reg string) string {
	return fmt.Sprintf("%s:%s:%s", chn, dev, reg)
}

// endregion
