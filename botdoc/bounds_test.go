package botdoc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_intBounds(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bound
	}{
		{
			"BetweenHyphen",
			"Values between 1-100 are accepted.",
			bound{
				Min: 1,
				Max: 100,
			},
		},
		{
			"BetweenAnd",
			"Must be between 1 and 100000 if specified",
			bound{
				Min: 1,
				Max: 100000,
			},
		},
		{
			"SemicolonHyphen",
			"measured in meters; 0-1500",
			bound{
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
		want  bound
	}{
		{
			"Characters",
			"Text of the message to be sent, 1-4096 characters after entities parsing",
			bound{
				Min: 1,
				Max: 4096,
			},
		}, {
			"Bytes",
			"Bot-defined invoice payload, 1-128 bytes.",
			bound{
				Min: 1,
				Max: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, stringBounds(tt.input))
		})
	}
}
