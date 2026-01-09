package result

import "testing"

func TestCollect(t *testing.T) {
	results := make(chan Item, 3)
	errs := make(chan error, 2)

	results <- Item{}
	results <- Item{}
	results <- Item{}

	errs <- errDummy
	errs <- errDummy

	close(results)
	close(errs)

	summary := Collect(results, errs)

	if summary.Success != 3 {
		t.Fatalf("expected 3 successes, got %d", summary.Success)
	}

	if summary.Errors != 2 {
		t.Fatalf("expected 2 errors, got %d", summary.Errors)
	}
}

var errDummy = dummyError("dummy")

type dummyError string

func (e dummyError) Error() string {
	return string(e)
}
