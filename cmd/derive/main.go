package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/krysopath/derive/app"
)

var (
	keyLen       int
	kdfRounds    int
	kdfFunction  string
	kdfHash      string
	kdfPurpose   string
	keyVersion   string
	outputFormat string
	version      string
	out          io.Writer = os.Stdout
	homeDir      string
)

func parseFlags() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(
			w,
			"derive %s: This is not helpful.\n\nderive [FLAGS] [slot]\n\nFlags:\n", version,
		)
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t[-%s %s] \t%v\n", f.Name, f.Value, f.Usage)
		})
	}

	flag.IntVar(&keyLen, "b", 32, "length of derived key in bytes")
	flag.IntVar(&kdfRounds, "c", 4096, "rounds for deriving key")
	flag.StringVar(&kdfHash, "h", "sha512", "hash for kdf function")
	flag.StringVar(&kdfFunction, "k", "pbkdf2", "kdf function for deriving key")
	flag.StringVar(&keyVersion, "v", "1000", "'versioned' key ")
	flag.StringVar(&outputFormat, "f", "bytes", "key output format: bytes|base64|hex|ascii|ascii@shell")
	flag.Parse()
	if flag.NArg() == 1 {
		kdfPurpose = flag.Arg(0)
	}
}

func init() {
	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	homeDir = filepath.Join(homeDir, ".derive.yaml")
	parseFlags()
}

func run() string {
	cfg := app.NewConfig(homeDir)
	return cfg.Run(
		kdfFunction, kdfHash, kdfPurpose, keyVersion, outputFormat, kdfRounds, keyLen,
	)
}

func main() {
	fmt.Fprintln(out, run())
}
