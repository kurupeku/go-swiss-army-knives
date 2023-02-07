package logs

import (
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

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errc:
			ef.WriteString(err.Error())
		}
	}
}
