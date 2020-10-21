# Secret

_A simple command for extracting secret values from a JSON file_

## Synopsis

    secret [-h] [-f <CONFIG_FILE>] [<KEY>]

Argument `KEY` is the key of the value to extract.

Options:

* `-f <CONFIG_FILE>` will use `CONFIG_FILE` instead of the default
  `~/.secret.json`.
* `-v` verbose mode
* `-h` flag will print usage information and exit.

## Usage example

Assume `~/.secret.json` has the following JSON content:

    {
        "secrets": [
            { "key": "testkey", "value": "testvalue" }
        ]
    }

Then the command:

    > secret testkey

will return `testvalue`. So if you want to use it in a bash script to
define a value you can use something like:

    mysecret=$(secret testvalue)

## Build

Do 
    git clone https://github.com/callerobertsson/secret.git

    cd secret

    go build

and you have the `secret` executable in current directory.

To install it, put it somewhen in `$PATH`. I have mine in
`/usr/local/bin/secret`.


