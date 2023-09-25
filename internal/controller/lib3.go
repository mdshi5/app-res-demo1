package controller

import (
	"math/rand"
	"time"
)

func StringWithCharset(length int, charset string) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomStringGenerator(stringLength int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return StringWithCharset(stringLength, charset)
}

func getPointerRefernce[V int32 | string](value V) *V {
	return &value
}
