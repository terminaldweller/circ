package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/lrstanley/girc"
	"golang.org/x/net/proxy"
)

func main() {
	Server := flag.String("address", "", "IRC server address")
	Port := flag.Int("port", 6697, "IRC server port")
	TLS := flag.Bool("tls", true, "Use TLS for IRC connection")
	SkipTLSVerify := flag.Bool("skip-tls-verify", false, "Skip TLS certificate verification")
	Channel := flag.String("channel", "", "IRC channel to join")
	Nick := flag.String("nick", "botnick", "IRC nickname")
	SASLName := flag.String("sasl-user", "", "IRC SASL username")
	SASLPassword := flag.String("sasl-pass", "", "IRC SASL password")
	Message := flag.String("message", "Hello, IRC!", "Message to send to the channel")
	Target := flag.String("target", "", "Target for the message (channel or user)")
	ProxyURL := flag.String("proxy", "", "Proxy URL (e.g., socks5://user:pass@host:port)")

	flag.Parse()

	irc := girc.New(girc.Config{
		Server: *Server,
		Port:   *Port,
		SSL:    *TLS,
		Nick:   *Nick,
		User:   *Nick,
		Name:   *Nick,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: *SkipTLSVerify,
			ServerName:         *Server,
		},
		PingDelay:    60,
		PingTimeout:  30,
		GlobalFormat: true,
	})

	if *SASLName != "" && *SASLPassword != "" {
		irc.Config.SASL = &girc.SASLPlain{
			User: *SASLName,
			Pass: *SASLPassword,
		}
	}

	irc.Handlers.AddBg(girc.CONNECTED, func(c *girc.Client, _ girc.Event) {
		if *Target != "" {
			c.Cmd.Message(*Target, *Message)
		} else {
			c.Cmd.Join(*Channel)
			c.Cmd.Message(*Channel, *Message)
		}

		c.Quit("")
	})

	irc.Handlers.AddBg(girc.PRIVMSG, func(_ *girc.Client, e girc.Event) {
		log.Println("Received message:", e.String())
	})

	var dialer proxy.Dialer

	if *ProxyURL != "" {
		proxyURL, err := url.Parse(*ProxyURL)
		if err != nil {
			log.Printf("Invalid proxy URL: %v", err)
			panic(err)
		}

		dialer, err = proxy.FromURL(proxyURL, &net.Dialer{Timeout: 10 * time.Second})
		if err != nil {
			log.Printf("Failed to create proxy dialer: %v", err)
			panic(err)
		}
	}

	if err := irc.DialerConnect(dialer); err != nil {
		log.Printf("Failed to connect to IRC server: %v", err)
		panic(err)
	}
}
