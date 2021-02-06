package pjlink

import (
	"bytes"
	"errors"
	"fmt"
)

type line []byte

func newCommand(class byte, body [4]byte, param []byte) (line, error) {
	// can parameter be len(0)?
	if len(param) > 128 {
		return nil, fmt.Errorf("parameter must be less than 128 bytes")
	}

	l := make(line, _minCommandBytes+len(param))
	l[0] = '%'
	l[1] = class

	l[2] = body[0]
	l[3] = body[1]
	l[4] = body[2]
	l[5] = body[3]

	l[6] = ' '

	copy(l[7:], param)
	l[len(l)-1] = '\r'
	return l, nil
}

func (c line) Header() byte {
	if len(c) == 0 {
		return 0x00
	}

	return c[0]
}

func (c line) Body() [4]byte {
	if len(c) < 6 {
		return [4]byte{}
	}

	return [4]byte{c[2], c[3], c[4], c[5]}
}

func (c line) Parameter() []byte {
	if len(c) < _minCommandBytes || len(c) > _minCommandBytes+_maxParameterBytes {
		return nil
	}

	return c[7 : len(c)-1]
}

func (c line) IsAuth() bool {
	lower := bytes.ToLower(c)
	return bytes.HasPrefix(lower, []byte{'p', 'j', 'l', 'i', 'n', 'k'})
}

const (
	_minCommandBytes   = 1 + 1 + 4 + 1 + 1
	_maxParameterBytes = 128

	// separators
	_separatorResponse = '='
)

var (
	_bodyAuth        = [4]byte{'L', 'I', 'N', 'K'}
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

/*
func (c *command) UnmarshalBinary(data []byte) error {
	if len(data) < _minCommandBytes {
		return errors.New("data is too short")
	}

	switch {
	case data[0] != _commandHeader:
		return fmt.Errorf("invalid header %#x", data[0])
	case data[len(data)-1] != _terminator:
		return fmt.Errorf("invalid terminator %#x", data[len(data)-1])
	case data[6] != _separatorCommand && data[6] != _separatorResponse:
		return fmt.Errorf("invalid separator %#x", data[6])
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
*/

func (l line) Error() error {
	param := l.Parameter()

	switch {
	case !bytes.HasPrefix(bytes.ToUpper(param), []byte{'E', 'R', 'R'}):
		return nil
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', '1'}):
		return errors.New("undefined command")
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', '2'}):
		return errors.New("out of parameter")
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', '3'}):
		return errors.New("unavailable time")
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', '4'}):
		return errors.New("projector/display failure")
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', 'A'}):
		return errors.New("invalid password")
	}

	return fmt.Errorf("unknown error: %#x", param)
}
