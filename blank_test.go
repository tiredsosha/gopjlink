package gopjlink

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap/zaptest"
)

func TestBlank(t *testing.T) {
	is := is.New(t)
	log := zaptest.NewLogger(t)
	proj.log = log
	proj.pool.Logger = log.Sugar()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	is.NoErr(proj.SetPower(ctx, true))

	is.NoErr(proj.SetBlank(ctx, true))
	blank, err := proj.Blank(ctx)
	is.NoErr(err)
	is.True(blank)

	is.NoErr(proj.SetBlank(ctx, false))
	blank, err = proj.Blank(ctx)
	is.NoErr(err)
	is.True(!blank)

	is.NoErr(proj.SetPower(ctx, false))
}
