package storage

import (
	"crypto/rand"
	"encoding/hex"
)

func randomID(prefix string) (string, error) {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	return prefix + "_" + hex.EncodeToString(b[:]), nil
}
