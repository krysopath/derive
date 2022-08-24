package kdf

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

var (
	Hashes = map[string]func() hash.Hash{
		"sha512": sha512.New,
		"sha256": sha256.New,
		//		"blake2b-512": blake2b.New512,
	}
)

//func Blake2b_512() func() hash.Hash {
//
//}

type PBKDF2Opts struct {
	Passphrase, Salt, Purpose, Version, Hash string
	Rounds, KeyLen                           int
}

func NewPBKDF2(opts PBKDF2Opts) []byte {
	hash, ok := Hashes[opts.Hash]
	if !ok {
		panic(fmt.Sprintf("ERR: unkown hash: %s", opts.Hash))
	}
	pw := []byte(opts.Passphrase)
	salt := []byte(fmt.Sprintf(
		"%s:%s:%s:",
		opts.Purpose,
		opts.Version,
		opts.Salt),
	)
	return pbkdf2.Key(pw, salt, opts.Rounds, opts.KeyLen, hash)
}
