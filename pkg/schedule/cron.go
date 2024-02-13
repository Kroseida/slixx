package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/robfig/cron/v3"
	"reflect"
	"time"
)

var PARSER = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type CronKind struct {
	cancelFunc context.CancelFunc // Add a cancel function to the struct
}

type CronKindConfiguration struct {
	Expression string `json:"expression" slixx:"CRON" default:"0 3 * * *"` // Add the expression field to the configuration
}

func (kind *CronKind) GetName() string {
	return "CRON"
}

func (kind *CronKind) Initialize(configuration any, callback func()) error {
	var ctx context.Context
	ctx, kind.cancelFunc = context.WithCancel(context.Background()) // Initialize the cancel function

	parsedConfiguration := configuration.(*CronKindConfiguration)

	schedule, err := PARSER.Parse(parsedConfiguration.Expression)
	if err != nil {
		return err
	}

	scheduleAt := schedule.Next(time.Now())
	go func() {
		for {
			timer := time.NewTimer(scheduleAt.Sub(time.Now()))
			select {
			case <-ctx.Done(): // Listen for the cancellation signal
				timer.Stop()
				return
			case <-timer.C:
				scheduleAt = schedule.Next(time.Now())
				callback()
			}
		}
	}()
	return nil
}

func (kind *CronKind) Deactivate() error {
	if kind.cancelFunc != nil {
		kind.cancelFunc() // Cancel the context, which stops the goroutine
	}
	return nil
}

func (kind *CronKind) Parse(configurationJson string) (interface{}, error) {
	var configuration CronKindConfiguration
	err := json.Unmarshal([]byte(configurationJson), &configuration)
	if err != nil {
		return nil, err
	}
	if configuration.Expression == "" {
		return nil, errors.New("expression can not be empty")
	}
	_, err = PARSER.Parse(configuration.Expression)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func (kind *CronKind) DefaultConfiguration() interface{} {
	return reflect.New(reflect.TypeOf(CronKindConfiguration{})).Interface()
}
