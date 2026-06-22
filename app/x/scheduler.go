package x

import (
	"fmt"

	"github.com/hibiken/asynq"
)

type SchedulerClient struct {
	Scheduler *asynq.Scheduler
}

func NewSchedulerClient(scheduler *asynq.Scheduler) *SchedulerClient {
	return &SchedulerClient{Scheduler: scheduler}
}

func Scheduler() *SchedulerClient {
	scheduler, err := app[*SchedulerClient]()
	if err != nil {
		Logger().Error(fmt.Sprintf("SchedulerClient can't be resolved: %s", err.Error()))
		return nil
	}
	return scheduler
}

// Register registers a task with the scheduler
func (s *SchedulerClient) Register(cron string, task *asynq.Task, opts ...asynq.Option) (string, error) {
	return s.Scheduler.Register(cron, task, opts...)
}

// Unregister removes a task from the scheduler by entry ID
func (s *SchedulerClient) Unregister(entryID string) error {
	return s.Scheduler.Unregister(entryID)
}

// Shutdown implements the Shutdownable interface
func (s *SchedulerClient) Shutdown() error {
	if s.Scheduler != nil {
		s.Scheduler.Shutdown()
	}
	return nil
}
