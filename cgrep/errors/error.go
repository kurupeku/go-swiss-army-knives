package errors

import (
	"errors"
	"strings"
	"sync"
)

type ErrorLogs struct {
	sync.Mutex
	errs []error
}

var Store = &ErrorLogs{}

func Error() error {
	if hasError() {
		ss := make([]string, 1, len(Store.errs)+1)
		ss[0] = "[Error]"
		for _, e := range Store.errs {
			ss = append(ss, e.Error())
		}

		return errors.New(strings.Join(ss, "\n"))
	}

	return nil
}

func Set(err error) {
	Store.errs = append(Store.errs, err)
}

func hasError() bool {
	return len(Store.errs) > 0
}

func Reset() {
	Store = &ErrorLogs{}
}
