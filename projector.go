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

				time.Sleep(options.delay)
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

	readLine := func() (line, error) {
		line := line{}
		buf := make([]byte, 64)
		for !bytes.Contains(buf, []byte{'\r'}) {
			n, err := conn.Read(buf)
			if err != nil {
				return nil, fmt.Errorf("unable to read: %w", err)
			}

			line = append(line, buf[:n]...)
		}

		return line, nil
	}

	line, err := readLine()
	switch {
	case err != nil:
		return err
	case !line.IsAuth():
		return nil // just go ahead?
	}

	param := line.Parameter()
	switch {
	case len(param) == 0:
		return fmt.Errorf("empty parameter on auth line")
	case param[0] == '0':
		// no auth required
		return nil
	case len(param) != 2+8:
		return fmt.Errorf("invalid auth length")
	case param[0] != '1' && param[1] != ' ':
		return fmt.Errorf("invalid first two auth characters")
	}

	rand := string(param[2:])
	sum := md5.Sum([]byte(rand + pass))
	b := []byte(hex.EncodeToString(sum[:]))

	// append on a command so that it works?
	// i feel like i shouldn't need to do this, but it doesn't work without it
	// what i _think_ should happen is that we write this part, and then
	// the command is actually written in sendCommand. maybe it just happens too
	// late so the projector just assumes we failed auth? idk
	cmd, err := newCommand('1', _bodyPower, []byte{'?'})
	if err != nil {
		return fmt.Errorf("unable to build power command: %w", err)
	}

	b = append(b, cmd...)

	// send sum + command
	n, err := conn.Write(b)
	switch {
	case err != nil:
		return fmt.Errorf("unable to write password to connection: %w", err)
	case n != len(b):
		return fmt.Errorf("unable to write password to connection: wrote %v/%v bytes", n, len(b))
	}

	line, err = readLine()
	switch {
	case err != nil:
		return err
	case line.Error() != nil:
		return line.Error()
	}

	return nil
}

func (p *Projector) sendCommand(ctx context.Context, cmd line, coolOff time.Duration) (line, error) {
	var resp line

	err := p.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("unable to set connection deadline: %w", err)
		}

		p.log.Debug("Command line", zap.String("hex", fmt.Sprintf("%#x", cmd)), zap.ByteString("str", cmd))

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

		p.log.Debug("Response line", zap.String("hex", fmt.Sprintf("%#x", data)), zap.ByteString("str", data))

		resp = line(data)
		switch {
		case resp.Error() != nil:
			return resp.Error()
		case resp.IsAuth():
			return fmt.Errorf("invalid password")
		}

		time.Sleep(coolOff)
		return nil
	})

	return resp, err
}
