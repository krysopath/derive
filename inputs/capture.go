package inputs

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func Credentials() (string, string, error) {
	bytesSalt := []byte(os.Getenv("DERIVE_SALT"))
	if len(bytesSalt) < 16 {
		fmt.Fprint(os.Stderr, "ENV value not sufficient. Enter Salt: ")
		var err error
		bytesSalt, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", err
		}
	}

	fmt.Fprintln(os.Stderr, "Enter Password: ")
	bytesPassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	return strings.TrimSpace(string(bytesSalt)), strings.TrimSpace(string(bytesPassword)), nil
}
