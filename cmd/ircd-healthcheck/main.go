package main

import (
	"github.com/stormbit/ircd-healthcheck"

	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
)

/**
 * Converts a boolean value to an integer.
 */
func boolToInt(val bool) int {
	if val {
		return 1
	}
	return 0
}

/**
 * Runs healthcheck and returns result as status code.
 * @return void
 */
func main() {
	secure := flag.Bool("secure", false, "Whether or not to use SSL/TLS when connecting.")
	skipVerification := flag.Bool("skip-verification", false, "Whether or not to skip verifying the certificate.")
	server := flag.String("server", "irc.stormbit.net:6667", "Server and port to connect to. (format: irc.example.org:6667)")

	// Parse active command-line flags.
	flag.Parse()

	var success bool

	if *secure == true {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: *skipVerification,
		}

		conn, err := tls.Dial("tcp", *server, tlsConfig)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}

		success, err = healthcheck.RunHealthcheck(conn)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			success = false
		}
	} else {
		conn, err := net.Dial("tcp", *server)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}

		success, err = healthcheck.RunHealthcheck(conn)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			success = false
		}
	}

	os.Exit(boolToInt(success))
}
