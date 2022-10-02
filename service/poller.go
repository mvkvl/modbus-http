// region - header

package service

import (
	"errors"
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/model"
	"github.com/mvkvl/modbus-http/queue"
	"sync"
	"time"
)

// endregion

// region - links

// modbus: https://github.com/mvkvl/modbus
// chrono: https://github.com/procyon-projects/chrono

// endregion

// region - poller service

type poller struct {
	clients   map[string]*modbus.Client // channels' modbus client
	config    *model.Config             // channels config
	m         sync.Mutex                // mutex for concurrent access synchronization
	dc        map[string]any            // data cache
	cq        queue.FifoQueue           // commands queue
	scheduler *Scheduler                // poller task scheduler
}

type Poller interface {
	Start()
	Stop()
	Cycle()
	Read(key string) (any, error)
	WriteByte(key string, value uint8, callback func()) error
	WriteWord(key string, value uint16, callback func()) error
	WriteValue(key string, value float32, callback func()) error
}

func CreateModbusPoller(handlerFactory func(connection string, mode model.Mode) modbus.ClientHandler, config *model.Config) Poller {
	var clients = make(map[string]*modbus.Client)
	for _, chn := range config.Channels {
		handler := handlerFactory(chn.Connection, chn.Mode)
		client := modbus.NewClient(handler)
		clients[chn.Title] = &client
	}
	return &poller{
		clients: clients,
		config:  config,
		dc:      make(map[string]any),
		cq:      queue.CreateQueue(),
	}
}

// endregion

// region - poller API

func (s *poller) Start() {
	if nil == s.scheduler {
		sh := NewScheduler()
		s.scheduler = &sh
		(*s.scheduler).RunWithFixedDelay(s.cycle, 100*time.Millisecond)
	}
}
func (s *poller) Stop() {
	if nil != s.scheduler {
		(*s.scheduler).Stop()
		s.scheduler = nil
	}
}
func (s *poller) Cycle() {
	if nil == s.scheduler {
		s.cycle()
	}
}

func (s *poller) Read(key string) (any, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if val, ok := s.dc[key]; ok {
		return val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no metric found for key '%s'", key))
	}
}
func (s *poller) WriteByte(key string, value uint8, callback func()) error {
	//TODO implement me
	panic("implement me")
}
func (s *poller) WriteWord(key string, value uint16, callback func()) error {
	//TODO implement me
	panic("implement me")
}
func (s *poller) WriteValue(key string, value float32, callback func()) error {
	//TODO implement me
	panic("implement me")
}

// endregion

// region - private methods

func (s *poller) cycle() {
	var wg = &sync.WaitGroup{}
	for _, chn := range s.config.Channels {
		wg.Add(1)
		go func(chn model.Channel) {
			if err := s.readChannel(chn); err != nil {
				//time.Sleep(time.Duration(chn.CyclePause) * 5 * time.Millisecond)
			}
			//if err := s.writeChannel(chn); err != nil {
			//}
			wg.Done()
		}(chn)
	}
	wg.Wait()
}

func (s *poller) readChannel(chn model.Channel) error {
	log.Info("~~~~~~~~ read channel '%s' [cp: %d, rp: %d]", chn.Title, chn.GetCyclePause(), chn.GetRegisterPause())
	client := s.clients[chn.Title]
	reader := NewReader(s.config, client)
	for _, dev := range chn.Devices {
		for _, reg := range dev.Registers {
			if reg.Mode == model.RO || reg.Mode == model.RW {
				time.Sleep(time.Duration(chn.GetRegisterPause()) * time.Millisecond)
				l := fmt.Sprintf("read register: %s:%s:%s ->", chn.Title, dev.Title, reg.Title)
				result, err := reader.Read(&reg)
				if nil == err {
					log.Info("%s %f", l, result)
					s.dc[cacheKey(chn.Title, dev.Title, reg.Title)] = result
				} else {
					log.Info("%s Error: %s", l, err)
				}
			}
		}
	}
	time.Sleep(time.Duration(chn.GetCyclePause()) * time.Millisecond)
	return nil
}

func cacheKey(chn, dev, reg string) string {
	return fmt.Sprintf("%s:%s:%s", chn, dev, reg)
}

// endregion
