package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lrstanley/girc"
	"golang.org/x/net/proxy"
)

const (
	PingDelay      = 60 * time.Second
	PingTimeout    = 30 * time.Second
	ProxyTimeout   = 10 * time.Second
	DefaultTLSPort = 6697
)

func InteractiveCLI(client *girc.Client, channel string, isRaw bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if isRaw {
			err := client.Cmd.SendRaw(scanner.Text())
			if err != nil {
				log.Println("Error sending raw command:", err)

				continue
			}
		}

		client.Cmd.Message(channel, scanner.Text())
	}
}

func main() {
	Server := flag.String("address", "", "IRC server address")
	Port := flag.Int("port", DefaultTLSPort, "IRC server port")
	TLS := flag.Bool("tls", true, "Use TLS for IRC connection")
	SkipTLSVerify := flag.Bool("skip-tls-verify", false, "Skip TLS certificate verification")
	Channel := flag.String("channel", "", "IRC channel to join")
	Nick := flag.String("nick", "botnick", "IRC nickname")
	SASLName := flag.String("sasl-user", "", "IRC SASL username")
	SASLPassword := flag.String("sasl-pass", "", "IRC SASL password")
	Message := flag.String("message", "Hello, IRC!", "Message to send to the channel")
	Target := flag.String("target", "", "Target for the message (channel or user)")
	ProxyURL := flag.String("proxy", "", "Proxy URL (e.g., socks5://user:pass@host:port)")
	SendRaw := flag.Bool("send-raw", false, "Send raw command to the server")
	Interactive := flag.Bool("interactive", false, "Run in interactive mode (not implemented)")
	CertFile := flag.String("cert", "", "Path to TLS certificate file (optional)")
	KeyFile := flag.String("key", "", "Path to TLS key file (optional)")

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
		PingDelay:    PingDelay,
		PingTimeout:  PingTimeout,
		GlobalFormat: true,
	})

	if *SASLName != "" && *SASLPassword != "" {
		irc.Config.SASL = &girc.SASLPlain{
			User: *SASLName,
			Pass: *SASLPassword,
		}
	}

	if *CertFile != "" && *KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(*CertFile, *KeyFile)
		if err != nil {
			log.Printf("Failed to load TLS certificate: %v", err)
			panic(err)
		}

		irc.Config.TLSConfig.Certificates = []tls.Certificate{cert}
	}

	irc.Handlers.AddBg(girc.CONNECTED, func(client *girc.Client, _ girc.Event) {
		if *SendRaw {
			err := client.Cmd.SendRaw(*Message)
			if err != nil {
				log.Println(err)
			}

			return
		}

		if *Target != "" {
			if strings.HasPrefix(*Target, "#") {
				client.Cmd.Join(*Target)
			}
			client.Cmd.Message(*Target, *Message)
		} else {
			client.Cmd.Join(*Channel)
			client.Cmd.Message(*Channel, *Message)
		}

		if !*Interactive {
			client.Quit("")
		} else {
			if *Target != "" {
				go InteractiveCLI(client, *Target, *SendRaw)
			} else {
				go InteractiveCLI(client, *Channel, *SendRaw)
			}
		}
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

		dialer, err = proxy.FromURL(proxyURL, &net.Dialer{Timeout: ProxyTimeout})
		if err != nil {
			log.Printf("Failed to create proxy dialer: %v", err)
			panic(err)
		}
	}

	if err := irc.DialerConnect(dialer); err != nil {
		log.Printf("Failed to connect to IRC server: %v", err)
		panic(err)
	}

	if *Interactive {
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

		<-quitChannel
	}
}
