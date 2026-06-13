package botapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBadRequestError_Error(t *testing.T) {
	msg := "hello"
	require.Equal(t, (&BadRequestError{Message: msg}).Error(), msg)
}

func TestNotImplementedError_Error(t *testing.T) {
	msg := "hello"
	require.Equal(t, (&NotImplementedError{Message: msg}).Error(), msg)
	require.NotEmpty(t, (&NotImplementedError{}).Error())
}
