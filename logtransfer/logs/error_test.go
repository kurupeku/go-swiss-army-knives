package logs

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "",
			err:       errors.New("error occurred"),
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			errc := make(chan error, 1)
			defer cancel()
			go Error(ctx, errc)

			errc <- tt.err
			time.Sleep(1 * time.Second)

			f, err := os.Open(errorFilePath)
			if err != nil {
				t.Fatal(err)
			}

			defer os.Remove(errorFilePath)
			defer f.Close()
			b, err := io.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.err.Error(), string(b))

		})
	}
}
