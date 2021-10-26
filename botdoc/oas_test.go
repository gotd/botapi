package botdoc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBounds(t *testing.T) {
	require.Equal(t, bound{Max: 32, Min: 1}, stringBounds(`1-32 characters`))
	require.Equal(t, bound{}, stringBounds(`F-D characters`))
	require.Equal(t, bound{}, stringBounds(`F-1 characters`))
}
