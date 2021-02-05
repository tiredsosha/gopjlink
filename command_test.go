package pjlink

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

type commandBinaryMarshalTest struct {
	name           string
	command        command
	data           []byte
	marshalError   error
	unmarshalError error
}

var _commandBinaryMarshalTests = []commandBinaryMarshalTest{
	{
		name: "PowerOn",
		command: command{
			Class:     '1',
			Body:      _bodyPower,
			Parameter: []byte{'1'},
		},
		data: []byte{0x25, 0x31, 0x50, 0x4f, 0x57, 0x52, 0x20, 0x31, 0x0d},
	},
	{
		name: "PowerOff",
		command: command{
			Class:     '1',
			Body:      _bodyPower,
			Parameter: []byte{'0'},
		},
		data: []byte{0x25, 0x31, 0x50, 0x4f, 0x57, 0x52, 0x20, 0x30, 0x0d},
	},
	{
		name: "TooLong",
		command: command{
			Class:     '1',
			Body:      _bodyPower,
			Parameter: make([]byte, 130),
		},
		marshalError: errors.New("parameter must be less than 128 bytes"),
	},
}

func TestCommandBinaryMarshal(t *testing.T) {
	for _, tt := range _commandBinaryMarshalTests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			cmd := command{}

			data, err := tt.command.MarshalBinary()
			if err != nil {
				is.Equal(err, tt.marshalError)
				return
			}

			is.Equal(data, tt.data)

			err = cmd.UnmarshalBinary(data)
			if err != nil {
				is.Equal(err, tt.unmarshalError)
				return
			}

			is.Equal(cmd, tt.command)
		})
	}
}
