package pjlink

import (
	"bytes"
	"errors"
	"fmt"
)

type line []byte

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

func (l line) Header() byte {
	if len(l) == 0 {
		return 0x00
	}

	return l[0]
}

func (l line) Body() [4]byte {
	if len(l) < 6 {
		return [4]byte{}
	}

	return [4]byte{l[2], l[3], l[4], l[5]}
}

func (l line) Parameter() []byte {
	if len(l) < _minCommandBytes || len(l) > _minCommandBytes+_maxParameterBytes {
		return nil
	}

	return l[7 : len(l)-1]
}

func (l line) IsAuth() bool {
	lower := bytes.ToLower(l)
	return bytes.HasPrefix(lower, []byte{'p', 'j', 'l', 'i', 'n', 'k'})
}

var (
	ErrUndefinedCommand = errors.New("undefined command")
	ErrInvalidInput     = errors.New("nonexistent input source")
	ErrOutOfParameter   = errors.New("out of parameter")
	ErrUnavailableTime  = errors.New("unavailable time")
	ErrProjectorFailure = errors.New("projector/display failure")
	ErrAuth             = errors.New("invalid password")
)

func (l line) Error() error {
	param := l.Parameter()

	switch {
	case !bytes.HasPrefix(bytes.ToUpper(param), []byte{'E', 'R', 'R'}):
		return nil
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', '1'}):
		return ErrUndefinedCommand
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', '2'}):
		if l.Body() == _bodyInput {
			return ErrInvalidInput
		}

		return ErrOutOfParameter
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', '3'}):
		return ErrUnavailableTime
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', '4'}):
		return ErrProjectorFailure
	case bytes.EqualFold(param, []byte{'E', 'R', 'R', 'A'}):
		return ErrAuth
	}

	return fmt.Errorf("unknown error: %#x", param)
}
