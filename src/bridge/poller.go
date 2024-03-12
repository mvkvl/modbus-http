package bridge

import (
	"mbridge/model"
	"mbridge/util"
	"sync"
	"time"
)

type Poller interface {
	Start(title string)
	Stop(title string)
}

type pollerImpl struct {
	stopped    bool
	channel    *model.Channel
	readCmdChn chan<- Command
	quitChn    chan struct{}
	logger     util.Logger
	started    bool
	mutex      sync.Mutex
}

func CreatePoller(readCmdChn chan Command, channel *model.Channel) Poller {
	return &pollerImpl{
		stopped:    false,
		readCmdChn: readCmdChn,
		channel:    channel,
		logger:     util.GetLogger("poller"),
	}
}

func (p *pollerImpl) Start(title string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.started {
		return
	}
	p.started = true
	p.quitChn = make(chan struct{})

	go func() {
		p.logger.Info("start poller %s", title)
		defer func() {
			close(p.quitChn)
			p.logger.Info("shutdown poller %s", title)
		}()
		for {
			select {
			case <-time.After(time.Second):
				p.cycle()
			case <-p.quitChn:
				p.stopped = true
				time.Sleep(time.Millisecond * 100)
				return
				//default:
				//	time.Sleep(delay)
			}
		}
	}()
}
func (p *pollerImpl) Stop(title string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.started {
		return
	}
	p.started = false
	p.logger.Info("stop poller %s", title)
	p.quitChn <- struct{}{}
}

func (p *pollerImpl) cycle() {
	p.logger.Debug("polling channel %s (%d devices)", p.channel.Title, len(p.channel.Devices))
	for _, d := range p.channel.Devices {
		if p.stopped {
			p.logger.Debug("polling disabled; exit")
			break
		}
		p.logger.Debug("polling device: %s:%s", p.channel.Title, d.Title)
		for _, r := range d.Registers {
			if p.stopped {
				break
			}
			if r.Mode == model.RO || r.Mode == model.RW {
				cmd := NewReadCommand(p.channel, &d, &r)
				p.logger.Trace("polling register: %s:%s:%s", cmd.GetChannel().Title, cmd.GetDevice().Title, cmd.GetRegister().Title)
				p.readCmdChn <- cmd
				p.logger.Trace("sent read command for: %s:%s:%s", cmd.GetChannel().Title, cmd.GetDevice().Title, cmd.GetRegister().Title)
			}
			time.Sleep(p.channel.GetRegisterPause())
		}
		time.Sleep(p.channel.GetCyclePause())
	}
}
