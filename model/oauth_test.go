package model

import (
	"encoding/json"
	"testing"
)

func TestOidcUser_ToOauthUser_EmailVerifiedBoolTrue(t *testing.T) {
	u := &OidcUser{VerifiedEmail: json.RawMessage(`true`)}
	if !u.ToOauthUser().VerifiedEmail {
		t.Fatal("expected true for JSON boolean true")
	}
}

func TestOidcUser_ToOauthUser_EmailVerifiedBoolFalse(t *testing.T) {
	u := &OidcUser{VerifiedEmail: json.RawMessage(`false`)}
	if u.ToOauthUser().VerifiedEmail {
		t.Fatal("expected false for JSON boolean false")
	}
}

func TestOidcUser_ToOauthUser_EmailVerifiedStringTrue(t *testing.T) {
	u := &OidcUser{VerifiedEmail: json.RawMessage(`"true"`)}
	if !u.ToOauthUser().VerifiedEmail {
		t.Fatal("expected true for JSON string \"true\"")
	}
}

func TestOidcUser_ToOauthUser_EmailVerifiedStringFalse(t *testing.T) {
	u := &OidcUser{VerifiedEmail: json.RawMessage(`"false"`)}
	if u.ToOauthUser().VerifiedEmail {
		t.Fatal("expected false for JSON string \"false\"")
	}
}

func TestOidcUser_ToOauthUser_EmailVerifiedStringCaseInsensitive(t *testing.T) {
	for _, s := range []string{`"True"`, `"TRUE"`} {
		u := &OidcUser{VerifiedEmail: json.RawMessage(s)}
		if !u.ToOauthUser().VerifiedEmail {
			t.Fatalf("expected true for %s", s)
		}
	}
}

func TestOidcUser_ToOauthUser_EmailVerifiedNull(t *testing.T) {
	u := &OidcUser{VerifiedEmail: json.RawMessage(`null`)}
	if u.ToOauthUser().VerifiedEmail {
		t.Fatal("expected false for JSON null")
	}
}

func TestOidcUser_ToOauthUser_EmailVerifiedAbsent(t *testing.T) {
	u := &OidcUser{}
	if u.ToOauthUser().VerifiedEmail {
		t.Fatal("expected false when email_verified field is absent")
	}
}
