package memorybackend

import (
	"reflect"
	"testing"
)

func TestMemoryBackend_StorePluginData(t *testing.T) {
	botID := "SOME_ID"
	pluginID := "SOME_PLUGIN_ID"
	id := "ID"
	data := "SOME DATA"

	dataStorage := botStorage{
		botID: pluginStorage{
			pluginID: storage{
				id: data,
			},
		},
	}

	type args struct {
		botID      string
		pluginID   string
		identifier string
		data       interface{}
	}
	tests := []struct {
		name    string
		b       *MemoryBackend
		args    args
		wantErr bool
	}{
		{
			name: "Store some data",
			b:    New(),
			args: args{
				botID:      botID,
				pluginID:   pluginID,
				identifier: id,
				data:       data,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.StorePluginData(tt.args.botID, tt.args.pluginID, tt.args.identifier, tt.args.data); (err != nil) != tt.wantErr {
				t.Fatalf("MemoryBackend.StorePluginData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.b.storage, dataStorage) {
				t.Errorf("MemoryBackend.StorePluginData() resulted in stored data %v, want %v", tt.b.storage, dataStorage)
			}
		})
	}
}
