package botapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_pollOption(t *testing.T) {
	require.Equal(t, []byte{'0'}, pollOption(0))
}
