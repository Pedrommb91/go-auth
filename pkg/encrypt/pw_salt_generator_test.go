package encrypt

import (
	"regexp"
	"strings"
	"testing"
)

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
			if got := GeneratePassword(tt.args.length, tt.args.hasNumbers, tt.args.hasSymbols); len(got) != tt.args.length &&
				regexp.MustCompile(`\d`).MatchString(got) == tt.args.hasNumbers &&
				strings.ContainsAny(got, symbols) == tt.args.hasSymbols {
				t.Errorf("GeneratePassword() = %v", got)
			}
		})
	}
}
