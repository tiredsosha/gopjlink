package pjlink

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/byuoitav/connpool"
	"go.uber.org/zap"
)

type Projector struct {
	address string
	pool    *connpool.Pool
	log     *zap.Logger
}

func NewProjector(addr string, opts ...Option) *Projector {
	options := &options{
		ttl:   30 * time.Second,
		delay: 250 * time.Millisecond,
		log:   zap.NewNop(),
		port:  4352,
	}

	for _, o := range opts {
		o.apply(options)
	}

	return &Projector{
		log: options.log,
		pool: &connpool.Pool{
			TTL:   options.ttl,
			Delay: options.delay,
			NewConnection: func(ctx context.Context) (net.Conn, error) {
				dial := net.Dialer{}
				return dial.DialContext(ctx, "tcp", addr+":"+strconv.Itoa(options.port))
			},
			Logger: options.log.Sugar(),
		},
	}
}
