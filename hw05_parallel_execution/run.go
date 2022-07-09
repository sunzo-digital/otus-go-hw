package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	tasksChanel := make(chan Task, len(tasks))
	var errorCount int32

	for _, task := range tasks {
		tasksChanel <- task
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(tasksChanel chan Task, errorsCount *int32, wg *sync.WaitGroup) {
			defer wg.Done()

			for {
				select {
				case task := <-tasksChanel:
					if task() != nil {
						atomic.AddInt32(errorsCount, 1)
						if atomic.LoadInt32(errorsCount) >= int32(m) {
							return
						}
					}
				default:
					return
				}
			}
		}(tasksChanel, &errorCount, wg)
	}

	wg.Wait()

	if errorCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
