package model

import "testing"

const (
	EXPECTED_ID   = "some_id"
	EXPECTED_NAME = "some name"
)

func TestUser_IsValid(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		Nickname  string
		FirstName string
		LastName  string
		Email     string
		Tz        string
		Locale    string
		IsAdmin   bool
		IsBot     bool
		IsOwner   bool
		IsMod     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Set all necessary information",
			fields: fields{
				ID:   EXPECTED_ID,
				Name: EXPECTED_NAME,
			},
			want: true,
		},
		{
			name: "ID not set",
			fields: fields{
				Name: EXPECTED_NAME,
			},
			want: false,
		},
		{
			name: "ID not set",
			fields: fields{
				ID: EXPECTED_ID,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				Nickname:  tt.fields.Nickname,
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
				Email:     tt.fields.Email,
				Tz:        tt.fields.Tz,
				Locale:    tt.fields.Locale,
				IsAdmin:   tt.fields.IsAdmin,
				IsBot:     tt.fields.IsBot,
				IsOwner:   tt.fields.IsOwner,
				IsMod:     tt.fields.IsMod,
			}
			if got := u.IsValid(); got != tt.want {
				t.Errorf("User.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
