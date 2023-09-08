// Ref: https://gist.github.com/fracasula/38aa1a4e7481f9cedfa78a0cdd5f1865

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const isBase64 = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"

func TestEncryptDecryptMessage(t *testing.T) {
	key := "0123456789abcdef" // must be of 16 bytes for this example to work
	message := "Lorem ipsum dolor sit amet"

	encrypted, err := EncryptMessage(message, key)
	require.Nil(t, err)
	require.Regexp(t, isBase64, encrypted)

	decrypted, err := DecryptMessage(encrypted, key)
	require.Nil(t, err)
	require.Equal(t, message, decrypted)
}
