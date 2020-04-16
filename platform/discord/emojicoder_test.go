package discord

import "testing"

func Test_getAbyleBotterEmojiFromDiscordEmoji(t *testing.T) {
	type args struct {
		discordEmoji string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Request valid Emoji",
			args: args{discordEmoji: "7️⃣"},
			want: "seven",
		},
		{
			name:    "Request invalid Emoji",
			args:    args{discordEmoji: "Ä"},
			want:    "Ä",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAbyleBotterEmojiFromDiscordEmoji(tt.args.discordEmoji)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAbyleBotterEmojiFromDiscordEmoji() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAbyleBotterEmojiFromDiscordEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}
