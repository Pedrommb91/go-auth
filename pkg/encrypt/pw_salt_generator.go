package encrypt

import (
	"math/rand"
	"time"
)

const (
	voc     string = "abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers string = "0123456789"
	symbols string = "!@#$%&*+_-="
)

func GeneratePassword(length int, hasNumbers bool, hasSymbols bool) string {
	chars := voc
	if hasNumbers {
		chars = chars + numbers
	}
	if hasSymbols {
		chars = chars + symbols
	}
	return generatePassword(length, chars)
}

func generatePassword(length int, chars string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := ""
	for i := 0; i < length; i++ {
		password += string([]rune(chars)[r.Intn(len(chars))])
	}
	return password
}
