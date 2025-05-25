package pjlink

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	_powerOn      = "on"
	_powerOff     = "off"
	_powerCooling = "cooling"
	_powerWarmUp  = "warm-up"
)

func (p *Projector) Power(ctx context.Context) (bool, error) {
	if err := p.sem.Acquire(ctx, 1); err != nil {
		return false, err
	}
	defer p.sem.Release(1)

	state, err := p.power(ctx)
	if err != nil {
		return false, err
	}

	return state == _powerOn || state == _powerWarmUp, nil
}

func (p *Projector) power(ctx context.Context) (string, error) {
	cmd, err := newCommand('1', _bodyPower, []byte{'?'})
	if err != nil {
		return "", fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd)
	if err != nil {
		return "", fmt.Errorf("unable to send command: %w", err)
	}

	param := resp.Parameter()
	switch {
	case bytes.EqualFold(param, []byte{'0'}):
		return _powerOff, nil
	case bytes.EqualFold(param, []byte{'1'}):
		return _powerOn, nil
	case bytes.EqualFold(param, []byte{'2'}):
		return _powerCooling, nil
	case bytes.EqualFold(param, []byte{'3'}):
		return _powerWarmUp, nil
	}

	return "", fmt.Errorf("unknown power state: %#x", param)
}

// SetPower sets the power state of the projector.
func (p *Projector) SetPower(ctx context.Context, power bool) error {
	if err := p.sem.Acquire(ctx, 1); err != nil {
		return err
	}
	defer p.sem.Release(1)

	state := []byte{'0'}
	if power {
		state = []byte{'1'}
	}

	cmd, err := newCommand('1', _bodyPower, state)
	if err != nil {
		return fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("unable to send command: %w", err)
	}

	if !bytes.EqualFold(resp.Parameter(), []byte{'O', 'K'}) {
		return fmt.Errorf("unknown response: %#x", resp.Parameter())
	}

}
