package model

import (
	"encoding/json"
	"testing"
)

func TestOidcUser_ToOauthUser_EmailVerified(t *testing.T) {
	tests := []struct {
		name          string
		emailVerified string // raw JSON value
		want          bool
	}{
		{"boolean true", `true`, true},
		{"boolean false", `false`, false},
		{"string true", `"true"`, true},
		{"string false", `"false"`, false},
		{"string True (case-insensitive)", `"True"`, true},
		{"string TRUE (case-insensitive)", `"TRUE"`, true},
		{"null", `null`, false},
		{"absent (nil)", ``, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var raw json.RawMessage
			if tt.emailVerified != "" {
				raw = json.RawMessage(tt.emailVerified)
			}
			u := &OidcUser{
				OauthUserBase:     OauthUserBase{Name: "Test", Email: "test@example.com"},
				Sub:               "sub123",
				VerifiedEmail:     raw,
				PreferredUsername: "testuser",
			}
			got := u.ToOauthUser()
			if got.VerifiedEmail != tt.want {
				t.Errorf("VerifiedEmail = %v, want %v (input: %s)", got.VerifiedEmail, tt.want, tt.emailVerified)
			}
		})
	}
}
