package gopjlink

import (
	"bytes"
	"context"
	"fmt"
)

var _inputs = map[string][]byte{
	"RGB1": {'1', '1'},
	"RGB2": {'1', '2'},
	"RGB3": {'1', '3'},
	"RGB4": {'1', '4'},
	"RGB5": {'1', '5'},
	"RGB6": {'1', '6'},
	"RGB7": {'1', '7'},
	"RGB8": {'1', '8'},
	"RGB9": {'1', '9'},

	"VIDEO1": {'2', '1'},
	"VIDEO2": {'2', '2'},
	"VIDEO3": {'2', '3'},
	"VIDEO4": {'2', '4'},
	"VIDEO5": {'2', '5'},
	"VIDEO6": {'2', '6'},
	"VIDEO7": {'2', '7'},
	"VIDEO8": {'2', '8'},
	"VIDEO9": {'2', '9'},

	"DIGITAL1": {'3', '1'},
	"DIGITAL2": {'3', '2'},
	"DIGITAL3": {'3', '3'},
	"DIGITAL4": {'3', '4'},
	"DIGITAL5": {'3', '5'},
	"DIGITAL6": {'3', '6'},
	"DIGITAL7": {'3', '7'},
	"DIGITAL8": {'3', '8'},
	"DIGITAL9": {'3', '9'},

	"STORAGE1": {'4', '1'},
	"STORAGE2": {'4', '2'},
	"STORAGE3": {'4', '3'},
	"STORAGE4": {'4', '4'},
	"STORAGE5": {'4', '5'},
	"STORAGE6": {'4', '6'},
	"STORAGE7": {'4', '7'},
	"STORAGE8": {'4', '8'},
	"STORAGE9": {'4', '9'},

	"NETWORK1": {'5', '1'},
	"NETWORK2": {'5', '2'},
	"NETWORK3": {'5', '3'},
	"NETWORK4": {'5', '4'},
	"NETWORK5": {'5', '5'},
	"NETWORK6": {'5', '6'},
	"NETWORK7": {'5', '7'},
	"NETWORK8": {'5', '8'},
	"NETWORK9": {'5', '9'},
}

func (p *Projector) VideoInputs(ctx context.Context) (map[string]string, error) {
	if err := p.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer p.sem.Release(1)

	cmd, err := newCommand('1', _bodyInput, []byte{'?'})
	if err != nil {
		return nil, fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to send command: %w", err)
	}

	return map[string]string{
		"": inputName(resp.Parameter()),
	}, nil
}

func (p *Projector) SetVideoInput(ctx context.Context, output, input string) error {
	if err := p.sem.Acquire(ctx, 1); err != nil {
		return err
	}
	defer p.sem.Release(1)

	cmdInput, ok := _inputs[input]
	if !ok {
		return fmt.Errorf("unknown input")
	}

	cmd, err := newCommand('1', _bodyInput, cmdInput)
	if err != nil {
		return fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd)
	switch {
	case err != nil:
		return fmt.Errorf("unable to send command: %w", err)
	case !bytes.EqualFold(resp.Parameter(), []byte{'O', 'K'}):
		return fmt.Errorf("unexpected parameter: %#x", resp.Parameter())
	}

	return nil
}

func (p *Projector) inputList(ctx context.Context) ([]string, error) {
	if err := p.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer p.sem.Release(1)

	cmd, err := newCommand('1', _bodyInputList, []byte{'?'})
	if err != nil {
		return nil, fmt.Errorf("unable to build command: %w", err)
	}

	resp, err := p.sendCommand(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to send command: %w", err)
	}

	split := bytes.Split(resp.Parameter(), []byte{' '})

	var inputs []string
	for _, input := range split {
		inputs = append(inputs, inputName(input))
	}

	return inputs, nil
}

func inputName(input []byte) string {
	for k, v := range _inputs {
		if bytes.Equal(v, input) {
			return k
		}
	}

	return string(input)
}
