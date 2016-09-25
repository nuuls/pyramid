package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strings"
	"time"
)

var (
	conn net.Conn
	re   = regexp.MustCompile(`:(\w+)!\w+@\w+\.tmi\.twitch\.tv PRIVMSG #(\w+) :(.*)`)
)

var (
	messages     int
	messageLimit int = 18
	mod          bool
)

func connect() {
	var err error
	if cfg.Host != "" {
		conn, err = net.Dial("tcp", cfg.Host)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connected to", cfg.Host)
	} else {
		conn, err = tls.Dial("tcp", "irc.chat.twitch.tv:443", nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connected to irc.chat.twitch.tv:443 using tls")
	}

	send("PASS " + cfg.Pass)
	send("NICK " + cfg.Nick)
	send("CAP REQ twitch.tv/tags")
	send("CAP REQ twitch.tv/commands")
	send("JOIN #" + channel)
}

func send(msg string) {
	_, err := conn.Write([]byte(msg + "\r\n"))
	if err != nil {
		log.Fatal(err)
	}
}

func say(msg string) {
	for messages > messageLimit {
		time.Sleep(time.Second * 1)
		fmt.Println("waiting so you dont get global banned LUL")
	}
	fmt.Println("sent:", msg)
	send("PRIVMSG #" + channel + " :" + msg + " ")
	messages++
	time.AfterFunc(time.Second*35, func() { messages-- })
}

func readChat() {
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(line, "PING") {
			send(strings.Replace(line, "PING", "PONG", 1))
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) == 4 {
			user, _, msg := matches[1], matches[2], matches[3]
			fmt.Printf("%s: %s\n", user, msg)
		} else {
			if line != "" && !strings.HasPrefix(line, "@") {
				fmt.Println(line)
			}
			spl := strings.Split(line, " ")
			if strings.Contains(spl[0], "mod=1") && messageLimit != 95 {
				messageLimit = 95
				mod = true
				fmt.Println("mod detected, message limit now set to", messageLimit, " / 30 seconds")
			}
			if strings.Contains(spl[0], "mod=0") && messageLimit != 18 {
				messageLimit = 18
				mod = false
				fmt.Println("mod not detected, message limit now set to", messageLimit, " / 30 seconds")
			}
		}
	}
}
