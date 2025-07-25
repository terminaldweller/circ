# circ
irc cli

```txt
$ circ -help
Usage of ./circ:
  -address string
        IRC server address
  -channel string
        IRC channel to join
  -interactive
        Run in interactive mode (not implemented)
  -message string
        Message to send to the channel (default "Hello, IRC!")
  -nick string
        IRC nickname (default "botnick")
  -port int
        IRC server port (default 6697)
  -proxy string
        Proxy URL (e.g., socks5://user:pass@host:port)
  -sasl-pass string
        IRC SASL password
  -sasl-user string
        IRC SASL username
  -send-raw
        Send raw command to the server
  -skip-tls-verify
        Skip TLS certificate verification
  -target string
        Target for the message (channel or user)
  -tls
        Use TLS for IRC connection (default true)
```
