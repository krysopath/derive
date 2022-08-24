package inputs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/term"
)

func readFromTerminal(fd int, prompt string) (pw []byte, err error) {
	stdin := int(syscall.Stdin)
	oldState, errTerm := term.GetState(stdin)
	if errTerm != nil {
		return []byte(""), err
	}

	// restores terminal explicitly (echo input to stdout again)
	defer term.Restore(stdin, oldState)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		for _ = range sigch {
			// if SIGINT is received we must restore terminal
			term.Restore(stdin, oldState)
			os.Exit(1)
		}
	}()

	fmt.Fprint(os.Stderr, prompt)
	pw, err = terminal.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	return pw, err
}

func readFromStdin() (pw []byte, err error) {
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

func readPassword(prompt string) ([]byte, error) {
	fd := int(os.Stdin.Fd())
	if terminal.IsTerminal(fd) {
		return readFromTerminal(fd, prompt)
	}
	return readFromStdin()

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
