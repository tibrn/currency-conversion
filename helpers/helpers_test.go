package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomString(t *testing.T) {

	req := require.New(t)

	str, err := GenerateRandomString(-5)
	req.ErrorIs(err, ErrNegativeLength)

	str, err = GenerateRandomString(0)
	req.NoError(err)
	req.Len(str, 0)

	str, err = GenerateRandomString(5)
	req.NoError(err)
	req.Len(str, 5)

	str, err = GenerateRandomString(50)
	req.NoError(err)
	req.Len(str, 50)
}
