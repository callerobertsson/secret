// Package main secret fetcher command
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

const defaultConfigFileName = ".secret.json"

// Config holds the secret fetcher command configuration
type Config struct {
	Secrets []Secret `json:"secrets"`
}

// Secret is a tuple containing a key and a value
type Secret struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var (
	helpFlag         bool
	verboseFlag      bool
	configFileOption string
)

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Printf("could not get current user")
		os.Exit(1)
	}

	defaultConfig := fmt.Sprintf("%v%c%v", u.HomeDir, os.PathSeparator, defaultConfigFileName)
	flag.BoolVar(&helpFlag, "h", false, "Show usage information")
	flag.BoolVar(&verboseFlag, "v", false, "Print verbose info")
	flag.StringVar(&configFileOption, "f", defaultConfig, "Config file path")
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

	if len(flag.Args()) < 1 {
		errorf("error: missing key\n")
	}

	key := flag.Arg(0)

	for i := 0; i < len(config.Secrets); i++ {
		if config.Secrets[i].Key == key {
			fmt.Printf("%v\n", config.Secrets[i].Value)
			os.Exit(0)
		}
	}

	errorf("unknown key: %v\n", key)
}

func readConfigFile(f string) (c *Config, err error) {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	c = &Config{}
	err = json.Unmarshal(bs, c)

	return c, err
}

func verbosef(format string, a ...interface{}) {
	if verboseFlag {
		fmt.Printf(format, a...)
	}
}

func errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func usage() {
	fmt.Printf(`Secret Fetcher

Usage: %s [-h] [-f <config-file>] [key]

Fetch the value for key in the JSON formatted <config-file>.

Flags:
`, os.Args[0])

	flag.PrintDefaults()
}
