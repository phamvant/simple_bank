package utils

import (
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

func RandomString(length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		sb.WriteByte(alphabet[r.Intn(k)])
	}

	return sb.String()
}

func RandomName() string {
	return RandomString(6)
}

func RandomBalance() int64 {
	return RandomInt(0, 99999)
}

func RandomCurrency() string {
	currency := []string{"USD", "YEN", "VND"}
	return currency[r.Intn(3)]
}
