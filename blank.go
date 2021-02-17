package pjlink

import (
	"bytes"
	"context"
	"fmt"
)

func (p *Projector) Blank(ctx context.Context) (bool, error) {
	cmd, err := newCommand('1', _bodyMute, []byte{'?'})
	if err != nil {
		return false, fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd, 0)
	if err != nil {
		return false, fmt.Errorf("unable to send command: %w", err)
	}

	param := resp.Parameter()
	switch {
	case bytes.EqualFold(param, []byte{'1', '1'}):
		// video mute on, audio mute off
		return true, nil
	case bytes.Equal(param, []byte{'2', '1'}):
		// audio mute on, video mute off
		return false, nil
	case bytes.Equal(param, []byte{'3', '1'}):
		// video and audio mute on
		return true, nil
	case bytes.Equal(param, []byte{'3', '0'}):
		// video and audio mute off
		return false, nil
	}

	return false, fmt.Errorf("unknown blank state: %#x", param)
}

func (p *Projector) SetBlank(ctx context.Context, blank bool) error {
	var state []byte
	switch {
	case p.avOnlyMute && !blank:
		state = []byte{'1', '0'}
	case p.avOnlyMute && blank:
		state = []byte{'1', '1'}
	case !blank:
		state = []byte{'3', '0'}
	case blank:
		state = []byte{'3', '1'}
	}

	cmd, err := newCommand('1', _bodyMute, state)
	if err != nil {
		return fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd, 0)
	if err != nil {
		return fmt.Errorf("unable to send command: %w", err)
	}

	if !bytes.EqualFold(resp.Parameter(), []byte{'O', 'K'}) {
		return fmt.Errorf("unknown response: %#x", resp.Parameter())
	}

	return nil
}
