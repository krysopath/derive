package inputs

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func Credentials() (string, string, error) {
	bytesSalt := []byte(os.Getenv("DERIVE_SALT"))
	if len(bytesSalt) < 16 {
		return "", "", errors.New("env[DERIVE_SALT] value not sufficient or unset")
	}

	fmt.Fprint(os.Stderr, "Enter Secret Token (hold Yubikey 5secs)")
	bytesPassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Fprint(os.Stderr, "OK!\n")

	return strings.TrimSpace(string(bytesSalt)), strings.TrimSpace(string(bytesPassword)), nil
}
