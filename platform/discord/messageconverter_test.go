package discord

import "testing"

func Test_convertMessageFromRedseligg(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Convert a message with newlines",
			args: args{
				text: "Some\nText",
			},
			want: "Some\\nText",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertMessageFromRedseligg(tt.args.text); got != tt.want {
				t.Errorf("convertMessageFromRedseligg() = %v, want %v", got, tt.want)
			}
		})
	}
}
