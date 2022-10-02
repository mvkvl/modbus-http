package service

import (
	log "github.com/jeanphorn/log4go"
	"time"
)

// region - scheduler

type scheduler struct {
	started bool
}

type Scheduler interface {
	Stop()
	RunAfter(f func(), d time.Duration)
	RunAtFixedRate(f func(), d time.Duration)
	RunWithFixedDelay(f func(), d time.Duration)
}

func NewScheduler() Scheduler {
	return &scheduler{
		started: false,
	}
}

// endregion

// region - public API

func (s *scheduler) Stop() {
	s.started = false
}
func (s *scheduler) RunAfter(f func(), d time.Duration) {
	if !s.started {
		s.started = true
		time.AfterFunc(d, func() {
			f()
			s.started = false
		})
	} else {
		panic("already started")
	}
}
func (s *scheduler) RunAtFixedRate(f func(), d time.Duration) {
	if !s.started {
		s.started = true
		for {
			if !s.started {
				log.Info("stop fixed rate task")
				break
			}
			go f()
			time.Sleep(d)
		}
	} else {
		panic("already started")
	}
}
func (s *scheduler) RunWithFixedDelay(f func(), d time.Duration) {
	if !s.started {
		s.started = true
		go s.runWithFixedDelay(f, d)
	} else {
		panic("already started")
	}
}

// endregion

// region - private methods

func (s *scheduler) runWithFixedDelay(f func(), d time.Duration) {
	for {
		f()
		if !s.started {
			log.Info("stop fixed delay task")
			break
		}
		time.Sleep(d)
	}
}

// endregion
