package encrypt

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Pedrommb91/go-auth/pkg/errors"
)

func TestPasswordEncryptor_Encrypt(t *testing.T) {
	type args struct {
		plaintext string
		salt      string
		encPass   string
	}
	tests := []struct {
		name           string
		args           args
		wantEncryptErr error
		wantDecryptErr error
	}{
		{
			name: "Success on encryting/decrypting password",
			args: args{
				plaintext: "very-strong-password",
				salt:      "very-strong-salt",
				encPass:   "very-strong-encoding-password",
			},
			wantEncryptErr: nil,
			wantDecryptErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := NewPasswordEncryptor()
			got, err := enc.Encrypt(tt.args.plaintext, tt.args.salt, tt.args.encPass)
			if !errors.Equal(err, tt.wantEncryptErr) {
				t.Errorf("PasswordEncryptor.Encrypt() error = %v, wantErr %v", err, tt.wantEncryptErr)
				return
			}
			pw, err := enc.Decrypt(got, tt.args.salt, tt.args.encPass)
			if !errors.Equal(err, tt.wantDecryptErr) {
				t.Errorf("PasswordEncryptor.Decrypt() error = %v, wantErr %v", err, tt.wantDecryptErr)
				return
			}

			if pw != tt.args.plaintext {
				t.Errorf("PasswordEncryptor.Encrypt()/Decrypt() got = %v, want %v", pw, tt.args.plaintext)
				return
			}
		})
	}
}

func TestGeneratePassword(t *testing.T) {
	type args struct {
		length     int
		hasNumbers bool
		hasSymbols bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Generate password with size 20 and no number and no symbols",
			args: args{
				length:     20,
				hasNumbers: false,
				hasSymbols: false,
			},
		},
		{
			name: "Generate password with size 20 and no number and no symbols",
			args: args{
				length:     50,
				hasNumbers: true,
				hasSymbols: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := NewPasswordEncryptor()
			if got := enc.GenerateSalt(tt.args.length, tt.args.hasNumbers, tt.args.hasSymbols); len(got) != tt.args.length &&
				regexp.MustCompile(`\d`).MatchString(got) == tt.args.hasNumbers &&
				strings.ContainsAny(got, symbols) == tt.args.hasSymbols {
				t.Errorf("GeneratePassword() = %v", got)
			}
		})
	}
}
