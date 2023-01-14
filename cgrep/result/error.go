package result

import (
	"errors"
	"strings"
	"sync"
)

type ErrorLogs struct {
	sync.Mutex
	errs []error
}

var GlobalError = &ErrorLogs{}

func Error() error {
	if hasError() {
		ss := make([]string, 1, len(GlobalError.errs)+1)
		ss[0] = "[Error]"
		for _, e := range GlobalError.errs {
			ss = append(ss, e.Error())
		}

		return errors.New(strings.Join(ss, "\n"))
	}

	return nil
}

func SetError(err error) {
	GlobalError.errs = append(GlobalError.errs, err)
}

func hasError() bool {
	return len(GlobalError.errs) > 0
}

func ResetError() {
	GlobalError = &ErrorLogs{}
}
