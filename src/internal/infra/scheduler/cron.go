package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

// Deprecated: Use StopWithTimeout to prevent indefinite blocking.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
}

func (s *Scheduler) StopWithTimeout(timeout time.Duration) error {
	stopCtx := s.cron.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-stopCtx.Done():
		return nil
	case <-timer.C:
		return fmt.Errorf("scheduler shutdown timeout after %v", timeout)
	}
}

// Spec follows cron expression format or predefined schedules like "@every 1h".
func (s *Scheduler) AddFunc(spec string, cmd func()) error {
	_, err := s.cron.AddFunc(spec, cmd)
	return err
}
