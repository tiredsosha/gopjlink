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
	skipMarshal    bool
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
		name: "MarshalTooLong",
		command: command{
			Class:     '1',
			Body:      _bodyPower,
			Parameter: make([]byte, 130),
		},
		marshalError: errors.New("parameter must be less than 128 bytes"),
	},
	{
		name:           "UnmarshalTooShort",
		unmarshalError: errors.New("data is too short"),
		data:           []byte{0x01, 0x01},
		skipMarshal:    true,
	},
	{
		name:           "UnmarshalInvalidHeader",
		unmarshalError: errors.New("invalid header 0x1"),
		data:           []byte{0x01, 0x31, 0x50, 0x4f, 0x57, 0x52, 0x20, 0x30, 0x0d},
		skipMarshal:    true,
	},
	{
		name:           "UnmarshalInvalidTerminator",
		unmarshalError: errors.New("invalid terminator 0x1"),
		data:           []byte{0x25, 0x31, 0x50, 0x4f, 0x57, 0x52, 0x20, 0x30, 0x01},
		skipMarshal:    true,
	},
	{
		name:           "UnmarshalInvalidSeparator",
		unmarshalError: errors.New("invalid separator 0x1"),
		data:           []byte{0x25, 0x31, 0x50, 0x4f, 0x57, 0x52, 0x01, 0x30, 0x0d},
		skipMarshal:    true,
	},
	{
		name: "PowerOnResponse",
		command: command{
			Class:     '1',
			Body:      _bodyPower,
			Response:  true,
			Parameter: []byte{'O', 'K'},
		},
		data: []byte{0x25, 0x31, 0x50, 0x4f, 0x57, 0x52, 0x3d, 0x4f, 0x4b, 0x0d},
	},
}

func TestCommandBinaryMarshal(t *testing.T) {
	for _, tt := range _commandBinaryMarshalTests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			cmd := command{}

			if tt.skipMarshal {
				err := cmd.UnmarshalBinary(tt.data)
				if err != nil {
					is.Equal(err, tt.unmarshalError)
					return
				}

				is.Equal(cmd, tt.command)
				return
			}

			data, err := tt.command.MarshalBinary()
			if err != nil {
				is.Equal(err, tt.marshalError)
				return
			}

			is.Equal(data, tt.data)

			err = cmd.UnmarshalBinary(tt.data)
			if err != nil {
				is.Equal(err, tt.unmarshalError)
				return
			}

			is.Equal(cmd, tt.command)
		})
	}
}
