package worker

import (
	"context"
	"errors"
	"strings"
	"sync"

	"csvproc/internal/progress"
	"csvproc/internal/reader"
	"csvproc/internal/result"
)

var ErrInvalidRow = errors.New("invalid csv row")

func StartPool(
	ctx context.Context,
	n int,
	jobs <-chan reader.Row,
	results chan<- result.Item,
	errs chan<- error,
	tracker *progress.Tracker,
	wg *sync.WaitGroup,
) {
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}

					if len(job.Data) == 0 {
						errs <- ErrInvalidRow
						tracker.Done()
						continue
					}

					results <- result.Item{
						Value: strings.ToUpper(job.Data[0]),
					}
					tracker.Done()
				}
			}
		}()
	}
}
