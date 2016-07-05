package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const DEFAULT_KEYPATH = `~/.ssh/id_rsa`
const DEFAULT_PWD = `$RSA_KEY_PWD`

// Represents this app's possible configuration values
type Config struct {
	KeyPath string
	Pwd     string
	Server  bool
	Port    int
	Wait    int
}

// Generates and returns a new Config based on the command-line
func newConfig() Config {
	var (
		keyArg  = flag.String("key", DEFAULT_KEYPATH, "path to key file")
		print_v = flag.Bool("version", false, "print version and exit")
		server  = flag.Bool("server", false, "run key server instead of command")
		port    = flag.Int("port", SERVER_RECV_PORT, "server receiving port")
		wait    = flag.Int("wait", CLIENT_TIMEOUT, "client timeout, in seconds")
		pwd     = flag.String("pwd", DEFAULT_PWD, "password for encrypted RSA key")
	)
	flag.Parse()
	if *print_v {
		fmt.Printf("docker-ssh-exec version %s, built %s\n", VERSION, SOURCE_DATE)
		os.Exit(0)
	}
	// check arguments for validity
	if (len(flag.Args()) < 1) && (*server == false) {
		fmt.Println("ERROR: A command to execute is required:",
			" docker-ssh-exec [options] [command]")
		os.Exit(1)
	}
	keyPath := *keyArg
	if keyPath == DEFAULT_KEYPATH {
		home := os.Getenv(`HOME`)
		if home == `` {
			home = `/root`
		}
		keyPath = filepath.Join(home, `.ssh`, `id_rsa`)
	}
	rsaPasswd := *pwd
	if *pwd == DEFAULT_PWD {
		rsaPasswd = os.Getenv(`RSA_KEY_PWD`)
	}
	return Config{
		Server:  *server,
		KeyPath: keyPath,
		Pwd:     rsaPasswd,
		Port:    *port,
		Wait:    *wait,
	}
}
