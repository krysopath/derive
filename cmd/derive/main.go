package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/krysopath/derive/app"
	"github.com/krysopath/derive/inputs"
	"github.com/krysopath/derive/kdf"
	"gopkg.in/alessio/shellescape.v1"
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
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(
			w,
			"derive %s: This is not helpful.\n\nderive [FLAGS] [topic]\n\nFlags:\n", version,
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

func main() {
	salt, pass, err := inputs.Credentials()
	if err != nil {
		panic(err)
	}
	var dk []byte
	switch kdfFunction {
	case "pbkdf2":
		dk = kdf.NewPBKDF2(kdf.PBKDF2Opts{
			Passphrase: pass,
			Salt:       salt,
			Purpose:    kdfPurpose,
			Version:    keyVersion,
			Rounds:     kdfRounds,
			KeyLen:     keyLen * 2,
		})
	}
	switch outputFormat {
	case "ascii@escape":
		fmt.Fprintf(os.Stdout, "%s", shellescape.Quote(app.Coerce(dk, keyLen)))
	case "ascii":
		fmt.Fprintf(os.Stdout, "%s", app.Coerce(dk, keyLen))
	case "hex":
		fmt.Fprintf(os.Stdout, "%X", dk[:keyLen])
	case "base64":
		fmt.Fprintf(os.Stdout, "%s", base64.RawStdEncoding.EncodeToString(dk[:keyLen]))
	default:
		fmt.Fprintf(os.Stdout, "%s", dk[:keyLen])
	}

}
