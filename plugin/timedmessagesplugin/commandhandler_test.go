package timedmessagesplugin

import (
	"testing"
	"time"
)

func Test_splitTmCommand(t *testing.T) {
	validIntervalStr := "1m"
	validInterval, _ := time.ParseDuration(validIntervalStr)

	validOtherIntervalStr := "2h"
	validOtherInterval, _ := time.ParseDuration(validOtherIntervalStr)

	type args struct {
		text string
	}
	tests := []struct {
		name         string
		args         args
		wantC        string
		wantInterval time.Duration
		wantMsg      string
		wantErr      bool
	}{
		{
			name: "Valid add command",
			args: args{
				text: "!tm add " + validIntervalStr + " some text",
			},
			wantC:        "add",
			wantInterval: validInterval,
			wantMsg:      "some text",
		},
		{
			name: "Valid remove command",
			args: args{
				text: "!tm remove " + validOtherIntervalStr + " some other text",
			},
			wantC:        "remove",
			wantInterval: validOtherInterval,
			wantMsg:      "some other text",
		},
		{
			name: "Invalid command",
			args: args{
				text: "!tm blub",
			},
			wantErr: true,
		},
		{
			name: "Invalid add command",
			args: args{
				text: "!tm add 4dfdfd ssdsd",
			},
			wantC:   "add",
			wantErr: true,
		},
		{
			name: "Invalid remove command",
			args: args{
				text: "!tm remove 1mmm ssdsd",
			},
			wantC:   "remove",
			wantErr: true,
		},
		{
			name: "Invalid add command - empty message",
			args: args{
				text: "!tm add 1m",
			},
			wantErr: true,
		},
		{
			name: "Invalid remove command - empty message",
			args: args{
				text: "!tm remove 2h",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotInterval, gotMsg, err := splitTmCommand(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitTmCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotC != tt.wantC {
				t.Errorf("splitTmCommand() gotC = %v, want %v", gotC, tt.wantC)
			}
			if gotInterval != tt.wantInterval {
				t.Errorf("splitTmCommand() gotInterval = %v, want %v", gotInterval, tt.wantInterval)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("splitTmCommand() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
		})
	}
}
