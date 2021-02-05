package pjlink

import (
	"context"
	"fmt"
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

func (p *Projector) sendCommand(ctx context.Context, cmd command) (command, error) {
	var resp command

	req, err := cmd.MarshalBinary()
	if err != nil {
		return resp, fmt.Errorf("unable to marshal command: %w", err)
	}

	err = p.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("unable to set connection deadline: %w", err)
		}

		n, err := conn.Write(req)
		switch {
		case err != nil:
			return fmt.Errorf("unable to write to connection: %w", err)
		case n != len(req):
			return fmt.Errorf("unable to write to connection: wrote %v/%v bytes", n, len(req))
		}

		data, err := conn.ReadUntil(_terminator, deadline)
		if err != nil {
			return fmt.Errorf("unable to read from connection: %w", err)
		}

		if err := resp.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("unable to unmarshal response: %w", err)
		}

		return nil
	})

	return resp, err
}
