package hash

import (
	"crypto/sha512"
	"fmt"
)

func Calc(text string, rounds int) string {
	result := []byte(text)
	for i := 0; i < rounds; i++ {
		sum := sha512.Sum512(result)
		result = sum[:]
	}
	return fmt.Sprintf("%x", result)
}
