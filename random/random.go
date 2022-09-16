package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func String(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRune[rand.Intn(len(letterRune))]
	}

	return string(b)
}
