package app

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/alessio/shellescape"
	"github.com/krysopath/derive/inputs"
	"github.com/krysopath/derive/kdf"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v1"
)

type Kdf struct {
	Salt    string
	Purpose string
	Version string
	Rounds  int
	KeyLen  int
	Hash    string
	Kdf     string
}
type Slot struct {
	Last  int64
	Count int64
	Kdf   Kdf
}

func (s Slot) Hashed() string {
	jsonBytes, err := json.Marshal(&s.Kdf)
	if err != nil {
		panic(err)
	}
	h := sha256.New()
	h.Write(jsonBytes)
	return fmt.Sprintf("%x", h.Sum(nil))
}

type Config struct {
	VaultAddr string
	Path      string
	Slots     map[string]Slot
}

func (c *Config) Update() {
	yamlBytes, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(c.Path, yamlBytes, 0600)

}

func NewConfig(configPath string) Config {
	var cfg Config
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		cfg = Config{
			VaultAddr: "http://localhost:8200",
			Path:      configPath,
			Slots:     make(map[string]Slot),
		}
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
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
	var slot Slot

	switch kdfFunction {
	case "pbkdf2":
		opts := kdf.PBKDF2Opts{
			Passphrase: pass,
			Salt:       salt,
			Purpose:    kdfPurpose,
			Version:    keyVersion,
			Rounds:     kdfRounds,
			KeyLen:     keyLen * 2,
			Hash:       kdfHash,
		}
		mapstructure.Decode(opts.RawData(), &slot.Kdf)
		dk = kdf.NewPBKDF2(opts)
	default:
		panic(fmt.Sprintf("err: unknown kdf: %s", kdfFunction))
	}

	if s, ok := cfg.Slots[slot.Hashed()]; !ok {
		cfg.Slots[slot.Hashed()] = slot
	} else {
		s.Count++
		s.Last = time.Now().Unix()
		cfg.Slots[slot.Hashed()] = s
	}

	cfg.Update()

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
