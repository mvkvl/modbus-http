package service

import "time"

// region - scheduler type

type Scheduler struct {
	started bool
}

// endregion

// region - public API

type SchedulerAPI interface {
	RunAfter(f func(), d time.Duration)
	RunAtFixedRate(f func(), d time.Duration)
	RunWithFixedDelay(f func(), d time.Duration)
}

func (s *Scheduler) RunAfter(f func(), d time.Duration) {
	time.AfterFunc(d, f)
}

func (s *Scheduler) RunAtFixedRate(f func(), d time.Duration) {
	for {
		if !s.started {
			break
		}
		go f()
		time.Sleep(d)
	}
}

func (s *Scheduler) RunWithFixedDelay(f func(), d time.Duration) {
	go f()
	go s.runWithFixedDelay(f, d)
}

// endregion

// region - private methods

func (s *Scheduler) runWithFixedDelay(f func(), d time.Duration) {
	for {
		if !s.started {
			break
		}
		time.AfterFunc(d, f)
	}
}

// endregion

// region - constructors

func NewScheduler() SchedulerAPI {
	return &Scheduler{
		started: false,
	}
}

// endregion
