package kdf

import (
	"encoding/base64"
	"testing"
)

func TestHelloName(t *testing.T) {
	target := "ldfRGv9neG9cU2OTdshIUI8ZAAmI65mikG3wXuqT0l4"
	opts := PBKDF2Opts{
		"abc",
		"abc",
		"abc",
		"abc",
		1000,
		32,
	}
	this := base64.RawStdEncoding.EncodeToString(NewPBKDF2(opts))
	if this != target {
		t.Fatalf(`NewPBKDF2(opts) = %q, want match for %#q`, this, target)
	}
}
