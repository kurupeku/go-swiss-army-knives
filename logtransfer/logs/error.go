package logs

import (
	"bufio"
	"context"
	"os"
)

const errorFilePath = "error.log"

func Error(ctx context.Context, errc chan error) error {
	ef, err := os.Create(errorFilePath)
	if err != nil {
		return err
	}
	defer ef.Close()

	w := bufio.NewWriter(ef)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errc:
			_, _ = w.WriteString(err.Error())
			_ = w.Flush()
		}
	}
}
