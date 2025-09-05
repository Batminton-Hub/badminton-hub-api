package util

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

func RandomString(length int, alpha, numeric, symbol bool) string {
	const alphaSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numericSet = "0123456789"
	const symbolSet = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	var charset string
	if !alpha && !numeric && !symbol {
		charset = alphaSet
	}
	if alpha {
		charset = charset + alphaSet
	}
	if numeric {
		charset = charset + numericSet
	}
	if symbol {
		charset = charset + symbolSet
	}

	state := make([]byte, length)
	for i := range state {
		state[i] = charset[rand.Intn(len(charset))]
	}
	return string(state)
}

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateUUIDWithoutHyphens() string {
	u := uuid.New()
	return strings.ReplaceAll(u.String(), "-", "")
}
