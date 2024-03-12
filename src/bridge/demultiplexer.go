package bridge

import (
	"mbridge/util"
	"sync"
)

type Demultiplexer interface {
	Start(title string)
	Stop(title string)
}

type demultiplexerImpl struct {
	readCmdChn  <-chan Command
	writeCmdChn <-chan Command
	modbusChn   chan<- Command
	quitChn     chan struct{}
	logger      util.Logger
	started     bool
	mutex       sync.Mutex
}

func CreateDemultiplexer(readCmdChn, writeCmdChn, modbusChn chan Command) Demultiplexer {
	return &demultiplexerImpl{
		readCmdChn:  readCmdChn,
		writeCmdChn: writeCmdChn,
		modbusChn:   modbusChn,
		logger:      util.GetLogger("demux"),
	}
}

func (d *demultiplexerImpl) Start(title string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.started {
		return
	}
	d.started = true
	d.quitChn = make(chan struct{})

	go func() {
		d.logger.Info("start demultiplexer %s", title)
		defer func() {
			close(d.modbusChn)
			close(d.quitChn)
			d.logger.Info("shutdown demultiplexer %s", title)
		}()
		for {
			select {
			case cmd, ok := <-d.writeCmdChn:
				if ok {
					d.logger.Trace("multiplexing write command: %v", cmd.GetRegister().Title)
					d.modbusChn <- cmd
				}
			case cmd, ok := <-d.readCmdChn:
				if ok {
					d.logger.Trace("multiplexing read command: %v", cmd.GetRegister().Title)
					d.modbusChn <- cmd
				}
			case <-d.quitChn:
				return
				//default:
				//	time.Sleep(delay)
			}
		}
	}()
}
func (d *demultiplexerImpl) Stop(title string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if !d.started {
		return
	}
	d.started = false
	d.logger.Info("stop demultiplexer %s", title)
	d.quitChn <- struct{}{}
}
