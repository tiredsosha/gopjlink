package pjlink

import (
	"bytes"
	"context"
	"fmt"
	"time"
)

func (p *Projector) Power(ctx context.Context) (bool, error) {
	cmd, err := newCommand('1', _bodyPower, []byte{'?'})
	if err != nil {
		return false, fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd, 0)
	if err != nil {
		return false, fmt.Errorf("unable to send command: %w", err)
	}

	param := resp.Parameter()
	body := resp.Body()

	switch {
	case !bytes.EqualFold(body[:], _bodyPower[:]):
		return false, fmt.Errorf("unexpected body in response: %#x", body)
	case bytes.EqualFold(param, []byte{'0'}):
		return false, nil
	case bytes.EqualFold(param, []byte{'1'}):
		return true, nil
	case bytes.EqualFold(param, []byte{'2'}):
		return false, nil
	case bytes.EqualFold(param, []byte{'3'}):
		return true, nil
	}

	return false, fmt.Errorf("unknown power state: %#x", param)
}

func (p *Projector) SetPower(ctx context.Context, power bool) error {
	state := []byte{'0'}
	delay := time.Duration(0)
	if power {
		state = []byte{'1'}
		delay = (10 * time.Second) - p.pool.Delay
	}

	cmd, err := newCommand('1', _bodyPower, state)
	if err != nil {
		return fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd, delay)
	if err != nil {
		return fmt.Errorf("unable to send command: %w", err)
	}

	param := resp.Parameter()
	body := resp.Body()

	switch {
	case !bytes.EqualFold(body[:], _bodyPower[:]):
		return fmt.Errorf("unexpected body in response: %#x", body)
	case !bytes.EqualFold(param, []byte{'O', 'K'}):
		return fmt.Errorf("unknown response: %#x", param)
	}

	return nil
}
