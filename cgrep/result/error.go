package result

import (
	"errors"
	"strings"
	"sync"
)

type errorLogs struct {
	sync.Mutex
	errs []error
}

var es = &errorLogs{}

func Error() error {
	if hasError() {
		ss := make([]string, 1, len(es.errs)+1)
		ss[0] = "[Error]"
		for _, e := range es.errs {
			ss = append(ss, e.Error())
		}

		return errors.New(strings.Join(ss, "\n"))
	}

	return nil
}

func SetError(err error) {
	es.errs = append(es.errs, err)
}

func hasError() bool {
	return len(es.errs) > 0
}
