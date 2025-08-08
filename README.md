# circ

A simple irc cli tool that can be used in scripts to send messages/notifications to irc servers.

here are some examples:

```bash
circ -sasl-pass "mysaslpassword" -sasl-user mysaslusername -nick mynick -address "irc.awesome.com" -port 6697 -message "hello" -target "user"
circ -sasl-pass "mysaslpassword" -sasl-user mysaslusername -nick mynick -address "irc.awesome.com" -port 6697 -message "hello" -target "#channel"
circ -sasl-pass "mysaslpassword" -sasl-user mysaslusername -nick mynick -address "irc.awesome.com" -port 6697 -message "hello" -target "#channel" -interactive
rlwrap circ -key nick.key -cert nick.cer -nick mynick -address "irc.awesome.com" -port 6697 -message "hello" -target "user" -interactive
```

If you pass the `-interactive` flag, your connection will not get cut after the message is sent and you keep communicating using the stdin.

```txt
$ circ -help
Usage of ./circ:
  -address string
        IRC server address
  -channel string
        IRC channel to join
  -interactive
        Run in interactive mode
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

## Thanks

* [girc](https://github.com/lrstanley/girc)
