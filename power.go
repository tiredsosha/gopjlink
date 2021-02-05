package pjlink

import (
	"bytes"
	"context"
	"fmt"
)

func (p *Projector) Power(ctx context.Context) (bool, error) {
	cmd := command{
		Class:     '1',
		Body:      _bodyPower,
		Parameter: []byte{'?'},
	}

	resp, err := p.sendCommand(ctx, cmd)
	switch {
	case err != nil:
		return false, fmt.Errorf("unable to send command: %w", err)
	case resp.Error() != nil:
		return false, resp.Error()
	case !bytes.EqualFold(resp.Body[:], _bodyPower[:]):
		return false, fmt.Errorf("unexpected body in response: %#x", resp.Body)
	case bytes.EqualFold(resp.Parameter, []byte{'0'}):
		return false, nil
	case bytes.EqualFold(resp.Parameter, []byte{'1'}):
		return true, nil
	case bytes.EqualFold(resp.Parameter, []byte{'2'}):
		return false, nil
	case bytes.EqualFold(resp.Parameter, []byte{'3'}):
		return true, nil
	}

	return false, fmt.Errorf("unknown power state: %#x", resp.Parameter)
}
