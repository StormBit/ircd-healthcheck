package healthcheck

import (
	"fmt"
	"github.com/sorcix/irc"
	"math/rand"
	"net"
	"time"
)

/**
 * Returns a string to be used as a nickname.
 */
func healthcheckNameGenerator() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	byteRange := make([]byte, 10)
	for i := range byteRange {
		byteRange[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return "check_" + string(byteRange)
}

/**
 * Performs the healthcheck and returns a boolean reporting if all is well.
 * @return boolean
 */
func RunHealthcheck(conn net.Conn) (bool, error) {
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Establish a reader, and a writer.
	writer := irc.NewEncoder(conn)
	reader := irc.NewDecoder(conn)

	// Seed the PRNG source, and use it to generate a random nickname.
	rand.Seed(time.Now().UnixNano())
	nick := healthcheckNameGenerator()

	// Set up a read channel, and look for an 001 RPL as defined by RFC2812.
	authChannel := make(chan error)
	go func() {
		for {
			msg, _ := reader.Decode()
			switch msg.Command {

			// Notice is probably an AUTH notice.
			case "NOTICE":
				break

			// Yay, we authenticated successfully!
			case irc.RPL_WELCOME:
				authChannel <- nil
				return
				
			// Oh no, we didn't get an 001. :(
			default:
				authChannel <- fmt.Errorf("Received unexpected status: %s", msg.Command)
				return
			}
		}
	}()

	// Send messages to AUTH healthcheck agent with server.
	messages := []*irc.Message{}
	messages = append(messages, &irc.Message{
		Command: irc.NICK,
		Params:  []string{nick},
	})
	messages = append(messages, &irc.Message{
		Command:  irc.USER,
		Params:   []string{nick, "0", "*"},
		Trailing: "https://github.com/StormBit/ircd-healthcheck",
	})
	for _, msg := range messages {
		if err := writer.Encode(msg); err != nil {
			break
		}
	}

	// Wait for a result from authchannel
	err := <-authChannel
	conn.Close()

	return err != nil, err
}
