package pjlink

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestPower(t *testing.T) {
	is := is.New(t)
	log := zaptest.NewLogger(t)
	proj.log = log
	proj.pool.Logger = log.Sugar()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pow, err := proj.Power(ctx)
	is.NoErr(err)
	log.Info("power", zap.Bool("power", pow))
}
