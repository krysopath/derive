package inputs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func readPassword(prompt string) (pw []byte, err error) {
	fd := int(os.Stdin.Fd())
	if terminal.IsTerminal(fd) {
		fmt.Fprint(os.Stderr, prompt)
		pw, err = terminal.ReadPassword(fd)
		fmt.Fprintln(os.Stderr)
		return
	}

	var b [1]byte
	for {
		n, err := os.Stdin.Read(b[:])
		// terminal.ReadPassword discards any '\r', so we do the same
		if n > 0 && b[0] != '\r' {
			if b[0] == '\n' {
				return pw, nil
			}
			pw = append(pw, b[0])
			// limit size, so that a wrong input won't fill up the memory
			if len(pw) > 1024 {
				err = errors.New("password too long")
			}
		}
		if err != nil {
			// terminal.ReadPassword accepts EOF-terminated passwords
			// if non-empty, so we do the same
			if err == io.EOF && len(pw) > 0 {
				err = nil
			}

			return pw, err
		}
	}
}

func Credentials() (string, string, error) {
	bytesSalt := []byte(os.Getenv("DERIVE_SALT"))
	if len(bytesSalt) < 16 {
		return "", "", errors.New("env[DERIVE_SALT] value not sufficient (>16chars) or just unset")
	}
	pw, err := readPassword("! Enter Secret Token (hold Yubikey 5secs)")
	if err != nil {
		return "", "", err
	}
	fmt.Fprint(os.Stderr, " ...OK\n")

	return strings.TrimSpace(string(bytesSalt)), strings.TrimSpace(string(pw)), nil
}
