// Package main secret fetcher command
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/callerobertsson/secret/secrets"
)

const rootSecretsFilePath = "/etc/secrets/config.json"
const userSecretsFileName = ".secrets.json"

var userSecretsFilePath = ""

var (
	helpFlag    bool
	listFlag    bool
	verboseFlag bool
	fileOption  string
)

func init() {
	u, err := user.Current()
	if err != nil {
		errorf("could not get current user")
	}

	flag.BoolVar(&helpFlag, "h", false, "Show usage information")
	flag.BoolVar(&listFlag, "l", false, "List available keys")
	flag.BoolVar(&verboseFlag, "v", false, "Print verbose info")
	flag.StringVar(&fileOption, "f", "", "Secrets file path")
	flag.Parse()

	userSecretsFilePath = fmt.Sprintf("%v%c%v", u.HomeDir, os.PathSeparator, userSecretsFileName)
}

func main() {

	if helpFlag {
		usage()
		os.Exit(0)
	}

	cfg, err := readSecretsFiles()
	if err != nil {
		errorf("%v", err)
		os.Exit(1)
	}

	if listFlag {
		list(cfg)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		usage()
		errorf("\nerror: missing key\n")
	}

	key := flag.Arg(0)

	for i := 0; i < len(cfg.Secrets); i++ {
		if cfg.Secrets[i].Key == key {
			fmt.Printf("%v\n", cfg.Secrets[i].Value)
			os.Exit(0)
		}
	}

	errorf("unknown key: %v\n", key)
}

func readSecretsFiles() (*secrets.Secrets, error) {

	verbosef("reading root secrets file (if present) %q\n", rootSecretsFilePath)

	// TODO: add option to igore global secrets
	rootSecrets, err := secrets.NewFromPath(rootSecretsFilePath)
	if err != nil {
		verbosef("could not get root secrets: %v\n", err)
		rootSecrets = &secrets.Secrets{}
	}

	secretsFilePath := fileOption
	if secretsFilePath == "" {
		secretsFilePath = userSecretsFilePath
	}

	verbosef("reading user secrets %q\n", secretsFilePath)

	userSecrets, err := secrets.NewFromPath(secretsFilePath)
	if err != nil {
		verbosef("could not get secrets from %q: %v\n", secretsFilePath, err)
		userSecrets = &secrets.Secrets{}
	}

	secrets := secrets.Secrets{}
	secrets.Merge(rootSecrets)
	secrets.Merge(userSecrets)

	return &secrets, nil
}

func list(s *secrets.Secrets) {

	if len(s.Secrets) < 1 {
		verbosef("No secrets\n")
		return
	}

	verbosef("Secrets keys:\n")

	for i := 0; i < len(s.Secrets); i++ {
		fmt.Printf("  %v\t\t%v\n", s.Secrets[i].Key, s.Secrets[i].Value)
	}
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

Usage: %s [<flags>] [key]

Fetch the value for key in a JSON formatted config file.

Flags:
`, os.Args[0])

	flag.PrintDefaults()
}
