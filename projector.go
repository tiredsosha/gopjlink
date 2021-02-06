package pjlink

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
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

				conn, err := dial.DialContext(ctx, "tcp", addr+":"+strconv.Itoa(options.port))
				if err != nil {
					return nil, err
				}

				if err := doAuth(ctx, conn, options.password); err != nil {
					conn.Close()
					return nil, fmt.Errorf("unable to do auth: %w", err)
				}

				return conn, nil
			},
			Logger: options.log.Sugar(),
		},
	}
}

func doAuth(ctx context.Context, conn net.Conn, pass string) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(5 * time.Second)
	}

	if err := conn.SetDeadline(deadline); err != nil {
		return fmt.Errorf("unable to set deadline: %w", err)
	}

	// read the line
	line := line{}
	buf := make([]byte, 64)
	for !bytes.Contains(buf, []byte{'\r'}) {
		n, err := conn.Read(buf)
		if err != nil {
			return fmt.Errorf("unable to read: %w", err)
		}

		line = append(line, buf[:n]...)
	}

	if !line.IsAuth() {
		return nil // just go ahead?
	}

	param := line.Parameter()
	switch {
	case len(param) == 0:
		return fmt.Errorf("empty parameter on auth line")
	case param[0] == '0':
		return nil // no auth required
	case len(param) != 2+8:
		return fmt.Errorf("invalid auth length")
	case param[0] != '1' && param[1] != ' ':
		return fmt.Errorf("invalid first two auth characters")
	}

	rand := string(param[2:])
	sum := md5.Sum([]byte(rand + pass))
	b := []byte(hex.EncodeToString(sum[:]))

	// send sum
	n, err := conn.Write(b)
	switch {
	case err != nil:
		return fmt.Errorf("unable to write password to connection: %w", err)
	case n != len(b):
		return fmt.Errorf("unable to write password to connection: wrote %v/%v bytes", n, len(b))
	}

	return nil
}

func (p *Projector) sendCommand(ctx context.Context, cmd line) (line, error) {
	var resp line

	err := p.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("unable to set connection deadline: %w", err)
		}

		p.log.Debug("Command line", zap.String("line", fmt.Sprintf("%#x", cmd)))

		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return fmt.Errorf("unable to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("unable to write to connection: wrote %v/%v bytes", n, len(cmd))
		}

		data, err := conn.ReadUntil('\r', deadline)
		if err != nil {
			return fmt.Errorf("unable to read from connection: %w", err)
		}

		p.log.Debug("Response line", zap.String("line", fmt.Sprintf("%#x", data)))

		resp = line(data)
		if resp.IsAuth() {
			return fmt.Errorf("invalid password")
		}

		return nil
	})

	return resp, err
}
