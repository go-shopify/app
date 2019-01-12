package app

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const stateSize = 16

func generateRandomState() (string, error) {
	var stateData [stateSize]byte

	if _, err := rand.Read(stateData[:]); err != nil {
		return "", fmt.Errorf("could not generate random state: %s", err)
	}

	return base64.URLEncoding.EncodeToString(stateData[:]), nil
}
