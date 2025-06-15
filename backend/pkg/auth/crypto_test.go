package auth_test

import (
	"backend/pkg/auth"
	"testing"

)

func TestSha256hex(t *testing.T) {
	s := "hello"
	wanted := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	got := auth.Sha256hex(s)
	if got != wanted {
		t.Errorf("Sha256hex(%q) = %q; want %q", s, got, wanted)
	}
}

func TestSafeHash(t *testing.T) {
	s := "TEST_USE_API_KEY"
	println(auth.SafeHash(s))
}