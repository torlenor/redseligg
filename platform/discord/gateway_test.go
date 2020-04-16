package discord

import (
	"fmt"
	"testing"

	"github.com/torlenor/abylebotter/webclient"
)

func TestBot_getGateway(t *testing.T) {
	wsGatewayURL := "ws://something"

	type fields struct {
		api              *webclient.MockClient
		apiResponse      webclient.APIResponse
		apiResponseError error
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{name: "Successful API call with valid gateway URL",
			fields: fields{
				api: webclient.NewMock(),
				apiResponse: webclient.APIResponse{
					StatusCode: 200,
					Body:       []byte(`{"url": "` + wsGatewayURL + `"}`),
				},
			},
			want: wsGatewayURL,
		},
		{name: "Failed API call",
			fields: fields{
				api:              webclient.NewMock(),
				apiResponseError: fmt.Errorf("Some error"),
			},
			want:    "",
			wantErr: true,
		},
		{name: "Successful API call with invalid json response",
			fields: fields{
				api: webclient.NewMock(),
				apiResponse: webclient.APIResponse{
					StatusCode: 200,
					Body:       []byte(`{{{"url": "` + wsGatewayURL + `"}`),
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := tt.fields.api
			api.ReturnOnCall = tt.fields.apiResponse
			api.ReturnOnCallError = tt.fields.apiResponseError
			b := &Bot{
				api: tt.fields.api,
			}
			got, err := b.getGateway()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.getGateway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Bot.getGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}
