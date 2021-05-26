package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/thoj/go-ircevent"
)

func main() {
	// Config data
	serverssl := os.Getenv("IRC_SERVER")
	_, useSSL := os.LookupEnv("IRC_SSL")
	channels := os.Getenv("IRC_CHANNELS")
	ircnick1 := os.Getenv("IRC_NICK")
	msgs := os.Getenv("IRC_MSGS")
	ignoreRegexp := os.Getenv("IRC_IGNOREHOSTS_REGEXP")

	// Compile the ignore regexp
	reIgnore, err := regexp.Compile(ignoreRegexp)
	if err != nil {
		log.Println("Bad regexp", err)
		log.Println("Set IRC_IGNOREHOSTS_REGEXP to a regexp that should ignore hostmasks")
		os.Exit(2)
	}

	// Setup the responses
	responses := make(map[string]string)
	for _, resp := range strings.Split(msgs, ";") {
		if len(resp) == 0 {
			continue
		}
		parts := strings.SplitN(resp, ":", 2)
		if len(parts) != 2 {
			log.Printf("Ignoring malformed message: %s", resp)
			continue
		}
		responses[parts[0]] = parts[1]
	}

	// Set up the connection
	conn := irc.IRC(ircnick1, ircnick1)
	if conn == nil {
		log.Println("conn is nil!  Did you set IRC_NICK?")
		return
	}

	// IRC auth config
	if _, use := os.LookupEnv("IRC_SASL"); use {
		conn.UseSASL = true
		conn.SASLLogin = os.Getenv("IRC_USER")
		conn.SASLPassword = os.Getenv("IRC_PASS")
		conn.SASLMech = "PLAIN"
	}

	// IRC startup
	conn.QuitMessage = "I've probably crashed..."
	conn.UseTLS = useSSL
	conn.AddCallback("001", func(e *irc.Event) {
		log.Println("Connected to Server")
		for _, channel := range strings.Split(channels, ",") {
			conn.Join(channel)
		}
	})
	conn.AddCallback("366", func(e *irc.Event) {
		log.Printf("Connected to Channel (%s)", e.Arguments[0])
	})
	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		switch e.Message() {
		case ircnick1 + ": hello?":
			conn.Privmsgf(e.Arguments[0], "%s: go away, I'm busy", e.Nick)
		}
	})

	conn.AddCallback("JOIN", func(e *irc.Event) {
		log.Println(e)
		if reIgnore.Match([]byte(e.Host)) && false {
			// Ignore hosts, this allows you to ensure
			// certain people don't get the notices sent
			// by the bot.
			log.Printf("Nick %s matched hostregexp (%s)", e.Nick, e.Host)
			return
		}
		m, ok := responses[e.Arguments[0]]
		if !ok {
			m = "You joined a channel I don't know about, maybe tell my handler what happened?"
		}
		conn.Privmsg(e.Nick, m)
	})

	err = conn.Connect(serverssl)
	if err != nil {
		log.Println(err)
		return
	}

	// Shutdown handler setup
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs

		log.Println(sig)
		done <- true
	}()

	// Startup Serving
	go conn.Loop()

	// Shut down
	<-done
	log.Println("exiting")
	for _, channel := range strings.Split(channels, ",") {
		conn.Privmsg(channel, "I'm going away now")
	}
	conn.Quit()
}
