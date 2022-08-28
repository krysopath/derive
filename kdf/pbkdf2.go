package kdf

import (
	"crypto/sha512"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

type PBKDF2Opts struct {
	Passphrase, Salt, Purpose, Version string
	Rounds, KeyLen                     int
}

func NewPBKDF2(opts PBKDF2Opts) []byte {
	pw := []byte(opts.Passphrase)
	salt := []byte(fmt.Sprintf(
		"%s:%s:%s:",
		opts.Purpose,
		opts.Version,
		opts.Salt),
	)
	return pbkdf2.Key(pw, salt, opts.Rounds, opts.KeyLen, sha512.New)
}
