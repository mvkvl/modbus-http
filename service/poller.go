// region - header

package service

import (
	"errors"
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/mvkvl/modbus"
	"github.com/mvkvl/modbus-http/model"
	"github.com/mvkvl/modbus-http/queue"
	"strings"
	"sync"
	"time"
)

// endregion

// region - links

// modbus: https://github.com/mvkvl/modbus
// chrono: https://github.com/procyon-projects/chrono

// endregion

// region - poller service

type writeRequest struct {
	channel   string
	reference string
	value     uint16
}

type poller struct {
	clients   map[string]*modbus.Client  // channels' modbus client
	config    *model.Config              // channels config
	m         sync.Mutex                 // mutex for concurrent access synchronization
	dc        map[string]*model.Metric   // data cache
	cq        map[string]queue.FifoQueue // commands queue (per channel)
	scheduler *Scheduler                 // poller task scheduler
}

type Poller interface {
	Start()
	Stop()
	Cycle()
	Read(key string) (*model.Metric, error)
	WriteByte(key string, value uint8, callback func()) error
	WriteWord(key string, value uint16, callback func()) error
	Metrics() []string
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
		dc:      make(map[string]*model.Metric),
		cq:      make(map[string]queue.FifoQueue),
	}
}

// endregion

// region - poller API

func (s *poller) Metrics() []string {
	var result []string
	for _, c := range s.config.Channels {
		for _, d := range c.Devices {
			for _, r := range d.Registers {
				metric := fmt.Sprintf("%s:%s:%s:%s", c.Title, d.Title, r.Title, r.Mode)
				if r.Mode == model.RO || r.Mode == model.RW {
					v, e := s.Read(cacheKey(c.Title, d.Title, r.Title))
					if nil != e {
						metric = fmt.Sprintf("%s : %s", metric, e)
					} else {
						if v.Timestamp.After(time.Now().Add(time.Second * time.Duration(-s.config.Ttl))) {
							metric = fmt.Sprintf("%s : 0x%04x %.2f %s", metric, v.RawValue, v.Value, v.Timestamp.Format(time.RFC3339))
						}
					}
				}
				result = append(result, metric)
			}
		}
	}
	return result
}

func (s *poller) Start() {
	if nil == s.scheduler {
		log.Info("poller start")
		sh := NewScheduler()
		s.scheduler = &sh
		(*s.scheduler).RunWithFixedDelay(s.cycle, 100*time.Millisecond)
	}
}
func (s *poller) Stop() {
	if nil != s.scheduler {
		log.Info("poller stop")
		(*s.scheduler).Stop()
		s.scheduler = nil
	}
}
func (s *poller) Cycle() {
	if nil == s.scheduler {
		s.cycle()
	}
}

func (s *poller) Read(key string) (*model.Metric, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if val, ok := s.dc[key]; ok {
		return val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no metric found for key '%s'", key))
	}
}
func (s *poller) WriteByte(reference string, value uint8, callback func()) error {
	err := s.write(reference, uint16(value))
	if nil == err && nil != callback {
		callback()
	}
	return err
}
func (s *poller) WriteWord(reference string, value uint16, callback func()) error {
	err := s.write(reference, value)
	//if nil == err && nil != callback {
	//	callback()
	//}
	return err
}

// endregion

// region - private methods

func (s *poller) cycle() {
	var wg = &sync.WaitGroup{}
	for _, chn := range s.config.Channels {
		wg.Add(1)
		go func(chn model.Channel) {
			if err := s.processChannel(chn); err != nil {
				//time.Sleep(time.Duration(chn.CyclePause) * 5 * time.Millisecond)
			}
			wg.Done()
		}(chn)
	}
	wg.Wait()
}
func (s *poller) processChannel(chn model.Channel) error {
	//log.Info("~~~~~~~~ read channel '%s' [cp: %d, rp: %d]", chn.Title, chn.GetCyclePause(), chn.GetRegisterPause())
	client := s.clients[chn.Title]
	mb := NewModbusClient(s.config, client)
	s.readChannel(&mb, chn)
	s.writeChannel(&mb, chn)
	time.Sleep(time.Duration(chn.GetCyclePause()) * time.Millisecond)
	return nil
}

func (s *poller) readChannel(mb *ModbusClient, chn model.Channel) {
	// read channel
	for _, dev := range chn.Devices {
		for _, reg := range dev.Registers {
			if reg.Mode == model.RO || reg.Mode == model.RW {
				time.Sleep(time.Duration(chn.GetRegisterPause()) * time.Millisecond)
				regKey := fmt.Sprintf("%s:%s:%s", chn.Title, dev.Title, reg.Title)
				l := fmt.Sprintf("read register: %s ->", regKey)
				raw, value, err := (*mb).Read(&reg)
				if nil == err {
					//log.Info("%s %f", l, result)
					m := &model.Metric{
						Key:       regKey,
						RawValue:  raw,
						Value:     value,
						Timestamp: time.Now(),
					}
					s.dc[cacheKey(chn.Title, dev.Title, reg.Title)] = m
				} else {
					log.Info("%s Error: %s", l, err)
				}
			}
		}
	}
}
func (s *poller) writeChannel(mb *ModbusClient, chn model.Channel) {
	if queue, exists := s.cq[chn.Title]; exists {
		for {
			if 0 == queue.Size() {
				break
			}
			rq, error := queue.Remove()
			if nil != error {
				log.Warn("write error: %s", error)
			}
			request := rq.(*writeRequest)
			err := (*mb).WriteRef(request.reference, request.value)
			if nil != err {
				log.Warn("write error: %s", err)
			}
			time.Sleep(time.Duration(chn.GetRegisterPause()) * time.Millisecond)
		}
	}
}

func (s *poller) channelKey(reference string) (string, error) {
	ref := strings.Split(reference, ":")
	if 3 != len(ref) {
		return "", errors.New(fmt.Sprintf("invalid reference passed: '%s'", reference))
	}
	c, err := s.config.FindChannelByTitle(strings.TrimSpace(ref[0]))
	if nil != err {
		return "", err
	}
	if nil == c {
		return "", errors.New("no channel found")
	}
	return c.Title, nil
}
func cacheKey(chn, dev, reg string) string {
	return fmt.Sprintf("%s:%s:%s", chn, dev, reg)
}

func (s *poller) write(reference string, value uint16) error {
	ck, err := s.channelKey(reference)
	if err != nil {
		return err
	}
	if _, exists := s.cq[ck]; !exists {
		s.cq[ck] = queue.CreateQueue()
	}
	return s.cq[ck].Insert(&writeRequest{
		channel:   ck,
		reference: reference,
		value:     value,
	})
}

// endregion
