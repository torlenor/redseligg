package giveawayplugin

import (
	"reflect"
	"testing"

	"github.com/torlenor/abylebotter/botconfig"
)

func Test_parseConfig(t *testing.T) {
	mods := []string{"user1", "user2"}

	type args struct {
		c botconfig.PluginConfig
	}
	tests := []struct {
		name    string
		args    args
		want    config
		wantErr bool
	}{
		{
			name: "Valid config, empty",
			args: args{
				c: botconfig.PluginConfig{
					Type: "giveaway",
				},
			},
			want: config{},
		},
		{
			name: "Valid config, with mods",
			args: args{
				c: botconfig.PluginConfig{
					Type: "giveaway",
					Config: map[string]interface{}{
						"mods": mods,
					},
				},
			},
			want: config{Mods: mods},
		},
		{
			name: "Valid config, with only one mod in list",
			args: args{
				c: botconfig.PluginConfig{
					Type: "giveaway",
					Config: map[string]interface{}{
						"mods": []string{"user1"},
					},
				},
			},
			want: config{Mods: []string{"user1"}},
		},
		{
			name: "Valid config, with mods and onlymods = true",
			args: args{
				c: botconfig.PluginConfig{
					Type: "giveaway",
					Config: map[string]interface{}{
						"mods":     mods,
						"onlymods": true,
					},
				},
			},
			want: config{Mods: mods, OnlyMods: true},
		},
		{
			name: "Invalid config, wrong type",
			args: args{
				c: botconfig.PluginConfig{
					Type: "something",
				},
			},
			wantErr: true,
		},
		{
			name: "Valid config, with mods and onlymods = true",
			args: args{
				c: botconfig.PluginConfig{
					Type: "giveaway",
					Config: map[string]interface{}{
						"onlymods": true,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfig(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
