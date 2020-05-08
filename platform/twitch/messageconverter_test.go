package twitch

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
			name: "Convert a Redseligg message containing a userid <@USERID>",
			args: args{
				text: "Some Text with <@USERID>",
			},
			want: "Some Text with USERID",
		},
		{
			name: "Convert a Redseligg message containing more than one userid <@USERID>",
			args: args{
				text: "Some Text with <@USERID> and with also a user <@SOMETHING ELSE> and text afterwards",
			},
			want: "Some Text with USERID and with also a user SOMETHING ELSE and text afterwards",
		},
		{
			name: "Convert a Redseligg message containing more than one userid <@USERID>",
			args: args{
				text: "<@test> and <@SOMETHING>",
			},
			want: "test and SOMETHING",
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
