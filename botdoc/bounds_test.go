package botdoc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_intBounds(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  intBound
	}{
		{
			"BetweenHyphen",
			"Values between 1-100 are accepted.",
			intBound{
				Min: 1,
				Max: 100,
			},
		},
		{
			"BetweenAnd",
			"Must be between 1 and 100000 if specified",
			intBound{
				Min: 1,
				Max: 100000,
			},
		},
		{
			"SemicolonHyphen",
			"measured in meters; 0-1500",
			intBound{
				Min: 0,
				Max: 1500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, intBounds(tt.input))
		})
	}
}

func Test_stringBounds(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  lengthBound
	}{
		{
			"Characters",
			"Text of the message to be sent, 1-4096 characters after entities parsing",
			lengthBound{
				Min: 1,
				Max: 4096,
			},
		},
		{
			"Bytes",
			"Bot-defined invoice payload, 1-128 bytes.",
			lengthBound{
				Min: 1,
				Max: 128,
			},
		},
		{
			name:  "BadBothLetter",
			input: `F-D characters`,
			want:  lengthBound{},
		},
		{
			name:  "BadOneLetter",
			input: `F-1 characters`,
			want:  lengthBound{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, stringBounds(tt.input))
		})
	}
}
