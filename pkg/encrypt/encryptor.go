package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	mathRand "math/rand"
	"strings"
	"time"

	"github.com/Pedrommb91/go-auth/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/xdg-go/pbkdf2"
)

const (
	ivSize     = 12
	iterations = 65536
	keyLen     = 32

	voc     string = "abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers string = "0123456789"
	symbols string = "!@#$%&*+_-="
)

type PasswordEncryptor struct{}

type Encryptor interface {
	Encrypt(plaintext, salt, encPass string) (string, error)
	Decrypt(ciphertext, salt, encPass string) (string, error)
	GenerateSalt(length int, hasNumbers bool, hasSymbols bool) string
}

func NewPasswordEncryptor() Encryptor {
	return &PasswordEncryptor{}
}

func (PasswordEncryptor) Encrypt(plaintext, salt, encPass string) (string, error) {
	const operation = "encryption.Encrypt"

	if plaintext == "" {
		return plaintext, nil
	}

	iv := make([]byte, ivSize)
	// http://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf
	// Section 8.2
	_, err := rand.Read(iv)
	if err != nil {
		return "", errors.Build(
			errors.WithOp(operation),
			errors.KindInternalServerError(),
			errors.WithSeverity(zerolog.WarnLevel),
			errors.WithError(err),
			errors.WithMessage("Failed to encrypt password"))
	}

	aesGCM, err := getAesGCM(salt, encPass)
	if err != nil {
		return "", err
	}

	data := aesGCM.Seal(nil, iv, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(iv) + "-" + base64.StdEncoding.EncodeToString(data), nil
}

func (PasswordEncryptor) Decrypt(ciphertext, salt, encPass string) (string, error) {
	const operation = "encryption.Decrypt"

	if ciphertext == "" {
		return ciphertext, nil
	}
	arr := strings.Split(ciphertext, "-")

	iv, err := base64.StdEncoding.DecodeString(arr[0])
	if err != nil {
		return "", errors.Build(
			errors.WithOp(operation),
			errors.KindInternalServerError(),
			errors.WithSeverity(zerolog.WarnLevel),
			errors.WithError(err),
			errors.WithMessage("Failed to decrypt password"))
	}

	data, err := base64.StdEncoding.DecodeString(arr[1])
	if err != nil {
		return "", errors.Build(
			errors.WithOp(operation),
			errors.KindInternalServerError(),
			errors.WithSeverity(zerolog.WarnLevel),
			errors.WithError(err),
			errors.WithMessage("Failed to decrypt password"))
	}

	aesGCM, err := getAesGCM(salt, encPass)
	if err != nil {
		return "", err
	}

	data, err = aesGCM.Open(nil, iv, data, nil)
	if err != nil {
		return "", errors.Build(
			errors.WithOp(operation),
			errors.KindInternalServerError(),
			errors.WithSeverity(zerolog.WarnLevel),
			errors.WithError(err),
			errors.WithMessage("Failed to decrypt password"))
	}

	return string(data), nil
}

func (PasswordEncryptor) GenerateSalt(length int, hasNumbers bool, hasSymbols bool) string {
	chars := voc
	if hasNumbers {
		chars = chars + numbers
	}
	if hasSymbols {
		chars = chars + symbols
	}
	return generatePassword(length, chars)
}

func getAesGCM(salt, encPass string) (cipher.AEAD, error) {
	const operation = "encryption.getAesGCM"

	key := pbkdf2.Key([]byte(encPass), []byte(salt), iterations, keyLen, sha256.New)

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, errors.Build(
			errors.WithOp(operation),
			errors.KindInternalServerError(),
			errors.WithError(err),
			errors.WithSeverity(zerolog.WarnLevel),
			errors.WithMessage("Failed to encrypt/decript password"))

	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Build(
			errors.WithOp(operation),
			errors.KindInternalServerError(),
			errors.WithSeverity(zerolog.WarnLevel),
			errors.WithError(err),
			errors.WithMessage("Failed to encrypt/decript password"))
	}

	return aesGCM, nil
}

func generatePassword(length int, chars string) string {
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	password := ""
	for i := 0; i < length; i++ {
		password += string([]rune(chars)[r.Intn(len(chars))])
	}
	return password
}
