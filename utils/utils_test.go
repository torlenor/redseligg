package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripCmd(t *testing.T) {
	assert := assert.New(t)

	result := StripCmd("!CMD test", "CMD")
	assert.Equal("test", result)

	result = StripCmd("test !CMD", "CMD")
	assert.Equal("test !CMD", result)

	result = StripCmd("!CMDtest", "CMD")
	assert.Equal("!CMDtest", result)

	result = StripCmd("!TEST test", "CMD")
	assert.Equal("!TEST test", result)

	result = StripCmd("!CMD2 test", "CMD")
	assert.Equal("!CMD2 test", result)
}

func TestGenerateErrorResponse(t *testing.T) {
	assert := assert.New(t)

	actualError := GenerateErrorResponse("Server error, try again later")
	expectedError := `{"error": "Server error, try again later"}`
	assert.Equal(expectedError, actualError)

	actualError = GenerateErrorResponse("something else")
	expectedError = `{"error": "something else"}`
	assert.Equal(expectedError, actualError)
}

func TestStringSliceContains(t *testing.T) {
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil slice",
			args: args{s: nil, e: "something"},
			want: false,
		},
		{
			name: "regular slice, containing entry",
			args: args{s: []string{"something", "something else"}, e: "something"},
			want: true,
		},
		{
			name: "regular slice, does not contain entry",
			args: args{s: []string{"something", "something else"}, e: "blub"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceContains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractSubCommandAndArgsString(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name           string
		args           args
		wantSubcommand string
		wantArgument   string
	}{
		{
			name: "message with subcommand and args",
			args: args{
				message: "subcmd some arguments",
			},
			wantSubcommand: "subcmd",
			wantArgument:   "some arguments",
		},
		{
			name: "message with subcommand and no args",
			args: args{
				message: "subcmd",
			},
			wantSubcommand: "subcmd",
			wantArgument:   "",
		},
		{
			name: "empty message",
			args: args{
				message: "",
			},
			wantSubcommand: "",
			wantArgument:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSubcommand, gotArgument := ExtractSubCommandAndArgsString(tt.args.message)
			if gotSubcommand != tt.wantSubcommand {
				t.Errorf("ExtractSubCommandAndArgsString() gotSubcommand = %v, want %v", gotSubcommand, tt.wantSubcommand)
			}
			if gotArgument != tt.wantArgument {
				t.Errorf("ExtractSubCommandAndArgsString() gotArgument = %v, want %v", gotArgument, tt.wantArgument)
			}
		})
	}
}
