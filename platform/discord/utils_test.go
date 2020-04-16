package discord

import "testing"

func Test_combineUsernameAndDiscriminator(t *testing.T) {
	name1 := "user1"
	disc1 := "1234"
	name2 := "user2"
	disc2 := "45"

	type args struct {
		username      string
		discriminator string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Combine an username with a discriminator",
			args: args{
				username:      name1,
				discriminator: disc1,
			},
			want: name1 + "#" + disc1,
		},
		{
			name: "Combine another username with a discriminator",
			args: args{
				username:      name2,
				discriminator: disc2,
			},
			want: name2 + "#" + disc2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineUsernameAndDiscriminator(tt.args.username, tt.args.discriminator); got != tt.want {
				t.Errorf("combineUsernameAndDiscriminator() = %v, want %v", got, tt.want)
			}
		})
	}
}
