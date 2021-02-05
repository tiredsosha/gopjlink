package pjlink

import (
	"errors"
	"fmt"
)

type command struct {
	Class     byte
	Body      [4]byte
	Response  bool
	Parameter []byte
}

const (
	_minCommandBytes = 1 + 1 + 4 + 1 + 1
	_terminator      = '\r'
	_commandHeader   = '%'

	// separators
	_separatorCommand  = ' '
	_separatorResponse = '='
)

var (
	// bodies
	_bodyPower       = [4]byte{'P', 'O', 'W', 'R'}
	_bodyInput       = [4]byte{'I', 'N', 'P', 'T'}
	_bodyMute        = [4]byte{'A', 'V', 'M', 'T'}
	_bodyError       = [4]byte{'E', 'R', 'S', 'T'}
	_bodyLamp        = [4]byte{'L', 'A', 'M', 'P'}
	_bodyInputList   = [4]byte{'I', 'N', 'S', 'T'}
	_bodyName        = [4]byte{'N', 'A', 'M', 'E'}
	_bodyManufacture = [4]byte{'I', 'N', 'F', '1'}
	_bodyProduct     = [4]byte{'I', 'N', 'F', '2'}
	_bodyInfo        = [4]byte{'I', 'N', 'F', 'O'}
	_bodyClass       = [4]byte{'C', 'L', 'S', 'S'}
)

func (c command) MarshalBinary() ([]byte, error) {
	// can parameter be len(0)?
	if len(c.Parameter) > 128 {
		return nil, fmt.Errorf("parameter must be less than 128 bytes")
	}

	data := make([]byte, _minCommandBytes+len(c.Parameter))

	data[0] = _commandHeader
	data[1] = c.Class

	data[2] = c.Body[0]
	data[3] = c.Body[1]
	data[4] = c.Body[2]
	data[5] = c.Body[3]

	data[6] = _separatorCommand
	if c.Response {
		data[6] = _separatorResponse
	}

	copy(data[7:], c.Parameter)
	data[len(data)-1] = _terminator

	return data, nil
}

func (c *command) UnmarshalBinary(data []byte) error {
	if len(data) < _minCommandBytes {
		return errors.New("data is too short")
	}

	switch {
	case data[0] != _commandHeader:
		return fmt.Errorf("invalid header %#x", data[0])
	case data[len(data)-1] != _terminator:
		return fmt.Errorf("invalid terminator %#x", data[0])
	case data[6] != _separatorCommand && data[6] != _separatorResponse:
		return fmt.Errorf("invalid separator %#x", data[0])
	}

	c.Class = data[1]
	c.Body = [4]byte{data[2], data[3], data[4], data[5]}
	c.Response = data[6] == _separatorResponse

	c.Parameter = make([]byte, len(data)-_minCommandBytes)
	if len(c.Parameter) > 0 {
		copy(c.Parameter, data[7:len(data)-1])
	}

	return nil
}
