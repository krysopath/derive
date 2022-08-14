package main

import (
	"crypto/sha512"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh/terminal"
)

func credentials() (string, string, error) {
	fmt.Fprint(os.Stderr, "Enter Salt: ")
	bytesSalt, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	fmt.Fprint(os.Stderr, "\nEnter Password: ")
	bytesPassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	return strings.TrimSpace(string(bytesSalt)), strings.TrimSpace(string(bytesPassword)), nil
}

var (
	keyLen    int
	kdfRounds int
	kdfHash   string
)

func init() {
	flag.IntVar(&keyLen, "b", 32, "length of derived key in bytes")
	flag.IntVar(&kdfRounds, "c", 4096, "rounds for deriving key")
	flag.StringVar(&kdfHash, "h", "sha512", "hash function for deriving key")
	flag.Parse()
}

func main() {
	salt, pass, err := credentials()
	if err != nil {
		panic(err)
	}
	var dk []byte
	switch kdfHash {
	default:
		dk = pbkdf2.Key([]byte(pass), []byte(salt), kdfRounds, keyLen, sha512.New)

	}
	fmt.Fprintf(os.Stdout, "\n%s", base64.RawStdEncoding.EncodeToString(dk))
}
