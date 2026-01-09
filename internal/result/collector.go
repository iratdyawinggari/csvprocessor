package result

type Item struct {
	Value string
}

type Summary struct {
	Success int
	Errors  int
}

func Collect(results <-chan Item, errs <-chan error) Summary {
	summary := Summary{}

	resultsOpen := true
	errsOpen := true

	for resultsOpen || errsOpen {
		select {
		case _, ok := <-results:
			if !ok {
				resultsOpen = false
				continue
			}
			summary.Success++

		case _, ok := <-errs:
			if !ok {
				errsOpen = false
				continue
			}
			summary.Errors++
		}
	}

	return summary
}
