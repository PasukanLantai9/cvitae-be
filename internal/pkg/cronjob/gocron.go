package cronjob

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"

	"time"
)

type cronJob struct {
	scheduler gocron.Scheduler
}

type CronJob interface {
	SetupCronJob(interval time.Duration, job func() error)
}

func NewInstance() CronJob {
	s, _ := gocron.NewScheduler()

	return &cronJob{
		scheduler: s,
	}
}

func (c *cronJob) SetupCronJob(interval time.Duration, job func() error) {
	_, err := c.scheduler.NewJob(gocron.DurationJob(interval), gocron.NewTask(func() {
		if err := job(); err != nil {
			fmt.Println(err)
		}
	}))
	if err != nil {
		panic(err)
	}

	go func() {
		c.scheduler.Start()
	}()
}
