package main;

import (
    "github.com/sorcix/irc"
    "log"
    "flag"
    "os"
    "net"
    "crypto/tls"
    "math/rand"
    "time"
    "errors"
)

/**
 * Converts a boolean value to an integer.
 */
func boolToInt(val bool) int {
    if val {
        return 1;
    }
    return 0;
}

/**
 * Returns true if an error is set.
 */
func isError(err error) bool {
    return err != nil;
}

/**
 * Returns a string to be used as a nickname.
 */
func healthcheckNameGenerator() string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789";
    byteRange := make([]byte, 16)
    for i := range byteRange {
        byteRange[i] = letters[rand.Int63() % int64(len(letters))]
    }
    return "healthcheck_" + string(byteRange)
}

/**
 * Performs the healthcheck and returns a boolean reporting if all is well.
 * @return boolean
 */
func runHealthcheck(conn net.Conn, err error) bool {
    if err != nil {
        log.Fatalln(err);
    } else {
        conn.SetDeadline(time.Now().Add(30 * time.Second))
        log.Println("Connected to server.")

        // Establish a reader, and a writer.
        writer := irc.NewEncoder(conn)
        reader := irc.NewDecoder(conn)

        // Seed the PRNG source, and use it to generate a random nickname.
        rand.Seed(time.Now().UnixNano());
        nick := healthcheckNameGenerator()

        // Set up a read channel, and look for an 001 RPL as defined by RFC2812.
        authChannel := make(chan error)
        go func() {
            var err error
            for {
                msg, _ := reader.Decode()
                log.Println(msg.String())
                switch msg.Command {

                    // Notice is probably an AUTH notice.
                    case "NOTICE":
                        break

                    // Yay, we authenticated successfully!
                    case irc.RPL_WELCOME:
                        authChannel <- err
                        return

                    // Oh no, we didn't get an 001. :(
                    default:
                        authChannel <- errors.New("Received unexpected status")
                        return
                }
            }
        }()

        // Send messages to AUTH healthcheck agent with server.
        messages := []*irc.Message{}
        messages = append(messages, &irc.Message{
            Command: irc.NICK,
            Params: []string{nick},
        })
        messages = append(messages, &irc.Message{
            Command: irc.USER,
            Params: []string{nick, "0", "*"},
            Trailing: "https://github.com/StormBit/ircd-healthcheck",
        })
        for _, msg := range messages {
            log.Println(msg)
            if err := writer.Encode(msg); err != nil {
                break
            }
        }

        // Wait for a result from authchannel
        err = <-authChannel
        conn.Close();
    }
    return isError(err);
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

    // We return a non-0 exit status.
    os.Exit(func() int {
        return boolToInt(func() bool {
            if *secure == true {
                return runHealthcheck(tls.Dial("tcp", *server, &tls.Config{
                    InsecureSkipVerify: *skipVerification,
                }))
            } else {
                return runHealthcheck(net.Dial("tcp", *server))
            }
        }())
    }())
}
