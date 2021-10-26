package botdoc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBounds(t *testing.T) {
	require.Equal(t, stringBound{Max: 32, Min: 1}, stringBounds(`1-32 characters`))
	require.Equal(t, stringBound{}, stringBounds(`F-D characters`))
	require.Equal(t, stringBound{}, stringBounds(`F-1 characters`))
}
