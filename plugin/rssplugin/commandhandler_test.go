package rssplugin

import (
	"testing"
)

func Test_splitRssCommand(t *testing.T) {
	validURLStr := "https://a.valid.url/something.xml"
	invalidURLStr := "dfdsfsdf.valid.url/something.xml"

	type args struct {
		text string
	}
	tests := []struct {
		name     string
		args     args
		wantC    string
		wantLink string
		wantErr  bool
	}{
		{
			name: "Valid add command",
			args: args{
				text: "add " + validURLStr,
			},
			wantC:    "add",
			wantLink: validURLStr,
		},
		{
			name: "Valid remove command",
			args: args{
				text: "remove " + validURLStr,
			},
			wantC:    "remove",
			wantLink: validURLStr,
		},
		{
			name: "Invalid command",
			args: args{
				text: "blub",
			},
			wantErr: true,
		},
		{
			name: "Invalid add command",
			args: args{
				text: "add ssdsd " + validURLStr,
			},
			wantC:    "add",
			wantLink: "ssdsd https://a.valid.url/something.xml",
			wantErr:  true,
		},
		{
			name: "Invalid remove command",
			args: args{
				text: "remove ssdsd " + validURLStr,
			},
			wantC:    "remove",
			wantLink: "ssdsd https://a.valid.url/something.xml",
			wantErr:  true,
		},
		{
			name: "Invalid add command - empty url",
			args: args{
				text: "add",
			},
			wantErr: true,
		},
		{
			name: "Invalid remove command - empty url",
			args: args{
				text: "remove",
			},
			wantErr: true,
		},
		{
			name: "Invalid add command - invalid url",
			args: args{
				text: "add " + invalidURLStr,
			},
			wantC:    "add",
			wantLink: "dfdsfsdf.valid.url/something.xml",
			wantErr:  true,
		},
		{
			name: "Invalid remove command - invalid url",
			args: args{
				text: "remove " + invalidURLStr,
			},
			wantC:    "remove",
			wantLink: "dfdsfsdf.valid.url/something.xml",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotLink, err := splitRssCommand(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitRssCommand() test = %s, error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if gotC != tt.wantC {
				t.Errorf("splitRssCommand() test = %s, gotC = %v, want %v", tt.name, gotC, tt.wantC)
			}
			if gotLink != tt.wantLink {
				t.Errorf("splitRssCommand() test = %s, gotMsg = %v, want %v", tt.name, gotLink, tt.wantLink)
			}
		})
	}
}
