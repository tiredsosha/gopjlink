package pjlink

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap/zaptest"
)

func TestPower(t *testing.T) {
	is := is.New(t)
	log := zaptest.NewLogger(t)
	proj.log = log
	proj.pool.Logger = log.Sugar()

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	is.NoErr(proj.SetPower(ctx, true))

	pow, err := proj.Power(ctx)
	is.NoErr(err)
	is.True(pow)

	is.NoErr(proj.SetPower(ctx, false))

	pow, err = proj.Power(ctx)
	is.NoErr(err)
	is.True(!pow)
}
