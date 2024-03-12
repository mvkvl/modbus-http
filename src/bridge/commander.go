package bridge

import (
	"mbridge/model"
	"mbridge/util"
)

type Commander interface {
	WriteRef(reference string, value uint16) error
}

type commanderImpl struct {
	channel     *model.Channel
	config      *model.Config
	writeCmdChn chan<- Command
	quitChn     chan struct{}
	logger      util.Logger
}

func CreateCommander(writeCmdChn chan Command, channel *model.Channel, config *model.Config) Commander {
	return &commanderImpl{
		channel:     channel,
		config:      config,
		writeCmdChn: writeCmdChn,
		quitChn:     make(chan struct{}),
		logger:      util.GetLogger("commander"),
	}
}

func (p *commanderImpl) WriteRef(reference string, value uint16) error {
	reg, err := p.config.FindRegister(reference)
	if err != nil {
		return err
	}
	cmd := NewWriteCommand(p.channel, reg.Device, reg, value)
	p.logger.Trace("writing register: %s:%s:%s", cmd.GetChannel().Title, cmd.GetDevice().Title, cmd.GetRegister().Title)
	p.writeCmdChn <- cmd
	return nil
}
