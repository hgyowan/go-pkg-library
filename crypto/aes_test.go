package aes

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Decrypt(t *testing.T) {
	res, err := Decrypt([]byte("FNySn3d7HMNE0fBsJeZy8vXuZ8RY31XbUIPMxhHNZViInb_ZO7pYyh79YBC1KaK7_cDMUFMDgOl5znunfdJXIo0z8tKUFSF2uz8l-n2DxDsfmrJ3AJeyweGvKyUOO814mXAgU-7oUae5s7rn7z2L2HlbP77eOruT0lRsxe3Co_bv_urK9B0="))
	require.NoError(t, err)
	t.Log(res)
}
