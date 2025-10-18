// Package main secret fetcher command
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"
)

const defaultConfigFileName = ".secret.json"

type Encoding string

const (
	encDefault = "plain"
	encRot13   = "rot13"
)

// Config holds the secret fetcher command configuration
type Config struct {
	Secrets        []Secret `json:"secrets"`
	configFilePath string
}

// Secret is a tuple containing a key and a value
type Secret struct {
	Key   string   `json:"key"`
	Value string   `json:"value"`
	Enc   Encoding `json:"enc"` // plain (default), rot13
}

var (
	helpFlag         bool
	listFlag         bool
	verboseFlag      bool
	configFileOption string
)

func init() {
	u, err := user.Current()
	if err != nil {
		errorf("could not get current user")
	}

	defaultConfig := fmt.Sprintf("%v%c%v", u.HomeDir, os.PathSeparator, defaultConfigFileName)

	flag.BoolVar(&helpFlag, "h", false, "Show usage information")
	flag.BoolVar(&listFlag, "l", false, "List available keys")
	flag.BoolVar(&verboseFlag, "v", false, "Print verbose info")
	flag.StringVar(&configFileOption, "c", defaultConfig, "Config file path")
	flag.Parse()
}

func main() {

	if helpFlag {
		usage()
		os.Exit(0)
	}

	verbosef("using config file %q\n", configFileOption)

	config, err := readConfigFile(configFileOption)
	if err != nil {
		errorf("error reading config file %q: %v\n", configFileOption, err)
	}
	config.configFilePath = configFileOption

	if listFlag {
		list(config)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		usage()
		errorf("\nerror: missing key\n")
	}

	key := flag.Arg(0)

	for _, s := range config.Secrets {
		if s.Key == key {
			val, err := decodedValue(s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding secret %q: %s\n", s.Key, err)
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "%v\n", val)
			os.Exit(0)
		}
	}

	errorf("unknown key: %v\n", key)
}

func decodedValue(s Secret) (string, error) {
	switch strings.ToLower(string(s.Enc)) {
	case "", encDefault:
		return s.Value, nil
	case encRot13:
		return rot13(s.Value), nil
	}

	return "", fmt.Errorf("unknown encoding: %q", s.Enc)
}

func readConfigFile(f string) (c *Config, err error) {
	bs, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	c = &Config{}
	err = json.Unmarshal(bs, c)

	return c, err
}

func rot13(s string) string {
	return strings.Map(rot13Rune, s)
}

func rot13Rune(r rune) rune {
	if r >= 'a' && r <= 'z' {
		if r >= 'm' {
			return r - 13
		}
		return r + 13
	} else if r >= 'A' && r <= 'Z' {
		if r >= 'M' {
			return r - 13
		}
		return r + 13
	}
	return r
}

func list(config *Config) {
	verbosef("Listing keys in config %v\n", config.configFilePath)

	for _, s := range config.Secrets {
		enc := s.Enc
		if enc == "" {
			enc = encDefault
		}

		// fmt.Printf("%s [%s]\n", s.Key, enc)
		fmt.Printf("%s\n", s.Key)
	}
}

func verbosef(format string, a ...any) {
	if verboseFlag {
		fmt.Printf(format, a...)
	}
}

func errorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func usage() {
	fmt.Printf(`Secret Fetcher

Usage: %s [<flags>] [key]

Fetch the value for key in a JSON formatted config file.

Values with enc "rot13" will be rot13 decrypted.

Flags:
`, os.Args[0])

	flag.PrintDefaults()
}
