package customcommandsplugin

import (
	"testing"
)

func Test_splitCommand(t *testing.T) {

	customCommand := "SomeCommand"

	type args struct {
		text string
	}
	tests := []struct {
		name              string
		args              args
		wantC             string
		wantCustomCommand string
		wantMsg           string
		wantErr           bool
	}{
		{
			name: "Valid add command",
			args: args{
				text: "add " + customCommand + " some text",
			},
			wantC:             "add",
			wantCustomCommand: customCommand,
			wantMsg:           "some text",
		},
		{
			name: "Valid remove command",
			args: args{
				text: "remove " + customCommand,
			},
			wantC:             "remove",
			wantCustomCommand: customCommand,
		},
		{
			name: "Invalid command",
			args: args{
				text: "blub",
			},
			wantErr: true,
		},
		{
			name: "Invalid add command - empty message",
			args: args{
				text: "add " + customCommand,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotCustomCommand, gotMsg, err := splitCommand(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitCommand() test = %s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if gotC != tt.wantC {
				t.Errorf("splitCommand() test = %s gotC = %v, want %v", tt.name, gotC, tt.wantC)
			}
			if gotCustomCommand != tt.wantCustomCommand {
				t.Errorf("splitCommand() test = %s gotCustomCommand = %v, want %v", tt.name, gotCustomCommand, tt.wantCustomCommand)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("splitCommand() test = %s gotMsg = %v, want %v", tt.name, gotMsg, tt.wantMsg)
			}
		})
	}
}
