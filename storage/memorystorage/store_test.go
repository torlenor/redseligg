package memorystorage

import (
	"reflect"
	"testing"

	"github.com/torlenor/abylebotter/storagemodels"
)

func TestMemoryStorage_StoreQuotesPluginQuote(t *testing.T) {
	botID := "SOME_ID"
	pluginID := "SOME_PLUGIN_ID"
	id := "ID"
	data := storagemodels.QuotesPluginQuote{
		Author:    "some author",
		AuthorID:  "some id",
		ChannelID: "some id",
		Text:      "some text",
	}

	dataStorage := botStorage{
		botID: pluginStorage{
			pluginID: memoryStorage{
				id: data,
			},
		},
	}

	type args struct {
		botID      string
		pluginID   string
		identifier string
		data       storagemodels.QuotesPluginQuote
	}
	tests := []struct {
		name    string
		b       *MemoryStorage
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
			if err := tt.b.StoreQuotesPluginQuote(tt.args.botID, tt.args.pluginID, tt.args.identifier, tt.args.data); (err != nil) != tt.wantErr {
				t.Fatalf("MemoryStorage.StoreQuotesPluginQuote() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.b.storage, dataStorage) {
				t.Errorf("MemoryStorage.StoreQuotesPluginQuote() resulted in stored data %v, want %v", tt.b.storage, dataStorage)
			}
		})
	}
}

func TestMemoryStorage_StoreQuotesPluginQuotesList(t *testing.T) {
	botID := "SOME_ID"
	pluginID := "SOME_PLUGIN_ID"
	id := "ID"
	data := storagemodels.QuotesPluginQuotesList{
		UUIDs: []string{"something", "something else"},
	}

	dataStorage := botStorage{
		botID: pluginStorage{
			pluginID: memoryStorage{
				id: data,
			},
		},
	}

	type args struct {
		botID      string
		pluginID   string
		identifier string
		data       storagemodels.QuotesPluginQuotesList
	}
	tests := []struct {
		name    string
		b       *MemoryStorage
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
			if err := tt.b.StoreQuotesPluginQuotesList(tt.args.botID, tt.args.pluginID, tt.args.identifier, tt.args.data); (err != nil) != tt.wantErr {
				t.Fatalf("MemoryStorage.StoreQuotesPluginQuotesList() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.b.storage, dataStorage) {
				t.Errorf("MemoryStorage.StoreQuotesPluginQuotesList() resulted in stored data %v, want %v", tt.b.storage, dataStorage)
			}
		})
	}
}
