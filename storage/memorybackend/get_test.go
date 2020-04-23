package memorybackend

import (
	"reflect"
	"testing"
)

func TestMemoryBackend_GetPluginData(t *testing.T) {
	botID := "SOME_ID"
	pluginID := "SOME_PLUGIN_ID"
	validID := "ID"
	data := "SOME DATA"

	dataStorage := botStorage{
		botID: pluginStorage{
			pluginID: storage{
				validID: data,
			},
		},
	}

	type args struct {
		botID      string
		pluginID   string
		identifier string
	}
	tests := []struct {
		name    string
		b       *MemoryBackend
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Get stored data",
			b:    New(),
			args: args{
				botID:      botID,
				pluginID:   pluginID,
				identifier: validID,
			},
			want: data,
		},
		{
			name: "Try to get data which does not exist",
			b:    New(),
			args: args{
				botID:      botID,
				pluginID:   pluginID,
				identifier: "ID_NOT_EXISTING",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set data storage
			tt.b.storage = dataStorage

			got, err := tt.b.GetPluginData(tt.args.botID, tt.args.pluginID, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryBackend.GetPluginData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryBackend.GetPluginData() = %v, want %v", got, tt.want)
			}
		})
	}
}
