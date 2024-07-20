package random

import "math/rand"

// NewRandomURL creates a new random URL with length 8 characters.
func NewRandomURL() string {
	b := make([]rune, 8)
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
