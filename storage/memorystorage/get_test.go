package memorystorage

import (
	"reflect"
	"testing"

	"github.com/torlenor/redseligg/storagemodels"
)

func TestMemoryStorage_GetQuotesPluginQuote(t *testing.T) {
	botID := "SOME_ID"
	pluginID := "SOME_PLUGIN_ID"
	validID := "ID"
	data := storagemodels.QuotesPluginQuote{
		Author:    "some author",
		AuthorID:  "some id",
		ChannelID: "some id",
		Text:      "some text",
	}

	dataStorage := botStorage{
		botID: pluginStorage{
			pluginID: memoryStorage{
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
		b       *MemoryStorage
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

			got, err := tt.b.GetQuotesPluginQuote(tt.args.botID, tt.args.pluginID, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryStorage.GetQuotesPluginQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MemoryStorage.GetQuotesPluginQuote() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMemoryStorage_GetQuotesPluginQuotesList(t *testing.T) {
	botID := "SOME_ID"
	pluginID := "SOME_PLUGIN_ID"
	validID := "ID"
	data := storagemodels.QuotesPluginQuotesList{
		UUIDs: []string{"something", "something else"},
	}

	dataStorage := botStorage{
		botID: pluginStorage{
			pluginID: memoryStorage{
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
		b       *MemoryStorage
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

			got, err := tt.b.GetQuotesPluginQuotesList(tt.args.botID, tt.args.pluginID, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryStorage.GetQuotesPluginQuotesList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MemoryStorage.GetQuotesPluginQuotesList() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
