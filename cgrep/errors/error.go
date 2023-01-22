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

// エラーを記録するためのグローバル変数
var Store = &ErrorLogs{}

// Store に保存されたエラーを error として返す関数
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

// error を渡すとグローバル変数上に保存する関数
func Set(err error) {
	Store.errs = append(Store.errs, err)
}

// error が一つ以上保存されている場合は true を返す関数
func hasError() bool {
	return len(Store.errs) > 0
}

// Store に保存されている内容をクリアする関数
func Reset() {
	Store = &ErrorLogs{}
}
