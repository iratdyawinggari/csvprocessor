package reader

import (
	"context"
	"encoding/csv"
	"os"

	"csvproc/internal/progress"
)

type Row struct {
	File string
	Data []string
}

func ReadCSV(
	ctx context.Context,
	path string,
	jobs chan<- Row,
	errs chan<- error,
	tracker *progress.Tracker,
) {
	file, err := os.Open(path)
	if err != nil {
		errs <- err
		return
	}
	defer file.Close()

	r := csv.NewReader(file)

	// Skip header
	_, _ = r.Read()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			record, err := r.Read()
			if err != nil {
				return
			}
			tracker.Increment()
			jobs <- Row{File: path, Data: record}
		}
	}
}
