package app

import (
	"encoding/base64"
	"testing"

	"github.com/krysopath/derive/kdf"
)

func TestHelloName(t *testing.T) {
	target := "ldfRGv9neG9cU2OTdshIUI8ZAAmI65mikG3wXuqT0l4"
	opts := kdf.PBKDF2Opts{
		"abc",
		"abc",
		"abc",
		"abc",
		"sha512",
		1000,
		32,
	}
	this := base64.RawStdEncoding.EncodeToString(kdf.NewPBKDF2(opts))
	if this != target {
		t.Fatalf(`NewPBKDF2(opts) = %q, want match for %#q`, this, target)
	}
}
