package gopjlink

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestInput(t *testing.T) {
	is := is.New(t)
	log := zaptest.NewLogger(t)
	proj.log = log
	proj.pool.Logger = log.Sugar()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	is.NoErr(proj.SetPower(ctx, true))

	input := "RGB1"
	is.NoErr(proj.SetVideoInput(ctx, "", input))

	inputs, err := proj.VideoInputs(ctx)
	is.NoErr(err)
	is.True(inputs[""] == input)

	input = "DIGITAL1"
	is.NoErr(proj.SetVideoInput(ctx, "", input))

	inputs, err = proj.VideoInputs(ctx)
	is.NoErr(err)
	is.True(inputs[""] == input)

	is.NoErr(proj.SetPower(ctx, false))
}

func TestInputList(t *testing.T) {
	is := is.New(t)
	log := zaptest.NewLogger(t)
	proj.log = log
	proj.pool.Logger = log.Sugar()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	inputs, err := proj.inputList(ctx)
	is.NoErr(err)
	log.Info("Input list", zap.Strings("inputs", inputs))
}
