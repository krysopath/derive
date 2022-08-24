package app

import (
	"encoding/base64"
	"fmt"

	"github.com/alessio/shellescape"
	"github.com/krysopath/derive/inputs"
	"github.com/krysopath/derive/kdf"
)

type Config struct {
	VaultAddr string
}

func NewConfig(configPath string) Config {
	return Config{
		VaultAddr: "http://localhost:8200",
	}
}

func (cfg *Config) Run(
	kdfFunction, kdfHash, kdfPurpose, keyVersion, outputFormat string,
	kdfRounds, keyLen int,
) string {
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
	default:
		panic(fmt.Sprintf("err: unknown kdf: %s", kdfFunction))
	}
	switch outputFormat {
	case "ascii@escape":
		return fmt.Sprintf("%s", shellescape.Quote(Coerce(dk, keyLen)))
	case "ascii":
		return fmt.Sprintf("%s", Coerce(dk, keyLen))
	case "hex":
		return fmt.Sprintf("%X", dk[:keyLen])
	case "base64":
		return fmt.Sprintf("%s", base64.RawStdEncoding.EncodeToString(dk[:keyLen]))
	default:
		return fmt.Sprintf("%s", dk[:keyLen])
	}
}
