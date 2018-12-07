package timeout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/slok/goresilience"
	grerrors "github.com/slok/goresilience/errors"
	"github.com/slok/goresilience/timeout"
)

func TestStaticLatency(t *testing.T) {
	err := errors.New("wanted error")

	tests := []struct {
		name   string
		cfg    timeout.StaticConfig
		f      goresilience.Func
		expErr error
	}{
		{
			name: "A command that has been run without timeout shouldn't return and error.",
			cfg: timeout.StaticConfig{
				Timeout: 1 * time.Second,
			},
			f: func(ctx context.Context) error {
				return nil
			},
			expErr: nil,
		},
		{
			name: "A command that has been run without timeout should return aerror result).",
			cfg: timeout.StaticConfig{
				Timeout: 1 * time.Second,
			},
			f: func(ctx context.Context) error {
				return err
			},
			expErr: err,
		},
		{
			name: "A command that has been run with timeout should return a fallback and don't let the function finish and return the err result.",
			cfg: timeout.StaticConfig{
				Timeout: 1,
			},
			f: func(ctx context.Context) error {
				time.Sleep(1 * time.Millisecond)
				return errors.New("wanted error")
			},
			expErr: grerrors.ErrTimeout,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			cmd := timeout.NewStatic(test.cfg, nil)
			err := cmd.Run(context.TODO(), test.f)

			assert.Equal(test.expErr, err)
		})
	}
}
