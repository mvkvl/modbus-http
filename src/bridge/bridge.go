package bridge

import (
	"fmt"
	"mbridge/model"
	"slices"
	"strings"
	"sync"
)

type Bridge interface {
	Start()
	Stop()
	Get(reference string) (*model.Metric, error)
	Set(reference string, value uint16) error
	List() []*model.Metric
	Regs() []*model.Register
	Flush()
}

type bridgeImpl struct {
	config     *model.Config
	started    bool
	mutex      sync.Mutex
	processors map[string]ChannelProcessor
}

func CreateBridge(config *model.Config) Bridge {
	br := &bridgeImpl{
		config: config,
	}
	return br
}

func (b *bridgeImpl) Start() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.started {
		return
	}
	b.started = true

	b.processors = make(map[string]ChannelProcessor, 0)
	for _, chn := range b.config.Channels {
		b.processors[chn.Title] = CreateProcessor(&chn, b.config)
	}
	for _, p := range b.processors {
		p.Start()
	}

}
func (b *bridgeImpl) Stop() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if !b.started {
		return
	}
	b.started = false
	for _, p := range b.processors {
		p.Stop()
	}
}

func (b *bridgeImpl) Get(reference string) (*model.Metric, error) {
	p, err := b.getProcessor(reference)
	if err != nil {
		return nil, err
	}
	return p.Cache().Get(reference), err
}
func (b *bridgeImpl) Set(reference string, value uint16) error {
	p, err := b.getProcessor(reference)
	if err != nil {
		return err
	}
	return p.Commander().WriteRef(reference, value)
}
func (b *bridgeImpl) List() []*model.Metric {
	var result []*model.Metric
	for _, p := range b.processors {
		chnList := p.Cache().List()
		for _, m := range chnList {
			result = append(result, m)
		}
	}
	return result
}
func (b *bridgeImpl) Regs() []*model.Register {
	var registers map[string]*model.Register = make(map[string]*model.Register, 0)
	for _, c := range b.config.Channels {
		for _, d := range c.Devices {
			for _, r := range d.Registers {
				registers[model.MetricKey(&r)] = &r
			}
		}
	}
	var keys = make([]string, 0)
	for k, _ := range registers {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	var result []*model.Register
	for _, k := range keys {
		result = append(result, registers[k])
	}
	return result
}
func (b *bridgeImpl) Flush() {
	for _, p := range b.processors {
		p.Cache().Flush()
	}
}
func (b *bridgeImpl) getProcessor(reference string) (ChannelProcessor, error) {
	channel := strings.Split(reference, ":")[0]
	processor, ok := b.processors[channel]
	if !ok {
		return nil, fmt.Errorf("could not find processor for %s", reference)
	}
	return processor, nil
}
