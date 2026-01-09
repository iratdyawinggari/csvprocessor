package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"csvproc/internal/progress"
	"csvproc/internal/reader"
	"csvproc/internal/result"
	"csvproc/internal/worker"
)

func main() {
	dir := flag.String("dir", "./data", "directory containing csv files")
	workers := flag.Int("workers", 4, "number of workers")
	flag.Parse()

	files, err := filepath.Glob(filepath.Join(*dir, "*.csv"))
	if err != nil || len(files) == 0 {
		fmt.Println("No CSV files found")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan reader.Row)
	results := make(chan result.Item)
	errs := make(chan error)

	tracker := progress.New()

	var readWG sync.WaitGroup
	var workerWG sync.WaitGroup

	// Start workers
	workerWG.Add(*workers)
	worker.StartPool(ctx, *workers, jobs, results, errs, tracker, &workerWG)

	// Start CSV readers
	for _, f := range files {
		readWG.Add(1)
		go func(file string) {
			defer readWG.Done()
			reader.ReadCSV(ctx, file, jobs, errs, tracker)
		}(f)
	}

	// Close jobs after readers finish
	go func() {
		readWG.Wait()
		close(jobs)
	}()

	// Close results & errors after workers finish
	go func() {
		workerWG.Wait()
		close(results)
		close(errs)
	}()

	// Start progress printer
	go tracker.Print()

	// Collect results
	summary := result.Collect(results, errs)

	// Wait for tracker to finish
	tracker.Wait()

	fmt.Println("\nProcessing completed")
	fmt.Printf("Success: %d | Errors: %d\n", summary.Success, summary.Errors)
}
