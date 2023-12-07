package parallel

import (
	"kroseida.org/slixx/pkg/statustype"
	"time"
)

type Executor[T any] struct {
	Contexts       []Context[T]
	Error          chan error
	StatusCallback func(status ExecutorStatus)
}

type Context[T any] struct {
	Items    []T
	Status   float64
	Finished bool
	Data     map[string]any
}

type ExecutorStatus struct {
	Percentage float64
	Message    string
	StatusType string
}

func NewExecutor[T any](items [][]T, callback func(status ExecutorStatus)) *Executor[T] {
	contexts := make([]Context[T], len(items))
	for index, item := range items {
		contexts[index] = Context[T]{Items: item}
	}
	return &Executor[T]{Contexts: contexts, Error: make(chan error), StatusCallback: callback}
}

func (executor Executor[T]) Run(execute func(index *int, ctx *Context[T])) {
	for index := range executor.Contexts {
		atomicIndex := index
		executor.Contexts[atomicIndex].Finished = false
		executor.Contexts[atomicIndex].Status = 0
		executor.Contexts[atomicIndex].Data = map[string]any{}

		go execute(&atomicIndex, &executor.Contexts[atomicIndex])
	}
	executor.waitWatchdog(executor.StatusCallback)
}

func (executor Executor[T]) waitWatchdog(callback func(status ExecutorStatus)) error {
	for {
		allFinished := true
		for _, context := range executor.Contexts {
			if !context.Finished {
				allFinished = false
				break
			}
		}
		if allFinished {
			break
		}

		percentage := 0.0
		for _, context := range executor.Contexts {
			percentage += context.Status
		}
		percentage /= float64(len(executor.Contexts))
		callback(ExecutorStatus{
			Percentage: percentage * 100,
			Message:    "Executing...",
			StatusType: statustype.Info,
		})

		var err error
		// Check for errors in the parallelError channel
		select {
		case res := <-executor.Error:
			err = res
		case <-time.After(1000 * time.Millisecond):
			err = nil
		}

		if err != nil {
			callback(ExecutorStatus{
				Percentage: 0,
				Message:    err.Error(),
				StatusType: statustype.Error,
			})
			return err
		}
	}
	return nil
}
