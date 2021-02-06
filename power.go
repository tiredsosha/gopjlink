package pjlink

import (
	"bytes"
	"context"
	"fmt"
)

func (p *Projector) Power(ctx context.Context) (bool, error) {
	cmd, err := newCommand('1', _bodyPower, []byte{'?'})
	if err != nil {
		return false, fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd)
	switch {
	case err != nil:
		return false, fmt.Errorf("unable to send command: %w", err)
	case resp.Error() != nil:
		return false, resp.Error()
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
