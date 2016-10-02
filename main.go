package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	Host  string `json:"host,omitempty"`
	Nick  string `json:"nick"`
	Pass  string `json:"pass"`
	Sleep int    `json:"sleep"`
}

var sleep time.Duration
var cfg config
var channel string

func init() {
	bs, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
		var file *os.File
		file, err = os.Create("config.json")
		cfg = config{
			Nick: "xd",
			Pass: "oauth:hiouefhiouerhjfgoirdhg",
		}
		bs, err = json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		file.Write(bs)
		fmt.Println("no config file found, i made one, just put in your username and oauth token")
		os.Exit(0)
	}
	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	sleep = time.Duration(cfg.Sleep) * time.Millisecond
	if len(os.Args) < 2 {
		fmt.Print("channel: ")
		r := bufio.NewReader(os.Stdin)
		line, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		channel = strings.ToLower(line)
		channel = strings.Replace(channel, "\n", "", -1)
		channel = strings.Replace(channel, "\r", "", -1)

	} else {
		channel = strings.ToLower(os.Args[1])
	}
	fmt.Println("channel set to", channel)
}

func main() {
	fmt.Print("usage: ")
	fmt.Println("make a pyramid: Kappa 5")
	fmt.Print("set new sleep: 1700\n\n\n")
	connect()
	go readInput()
	readChat()
}

func readInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, "\r", "", -1)
		spl := strings.Split(line, " ")
		switch len(spl) {
		case 1:
			if spl[0] == "" {
				fmt.Println("usage:")
				fmt.Println("\tmake a pyramid: Kappa 5")
				fmt.Println("\tset new sleep: 1700")
				break
			}
			slp, err := strconv.Atoi(spl[0])
			if err != nil {
				fmt.Println("error parsing sleep", err)
				break
			}
			sleep = time.Duration(slp) * time.Millisecond
			fmt.Println("sleep is now set to", sleep.String())
		case 2:
			count, err := strconv.Atoi(spl[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			emote := spl[0]
			go buildPyramid(count, emote)
		}
	}
}

func buildPyramid(count int, emote string) {
	if !mod && sleep.Seconds() < 1.2 {
		sleep = time.Millisecond * 1200
		fmt.Println("set sleep to", sleep.String(), "because you are not mod")
	}
	emote += " "
	for i := 1; i < count+1; i++ {
		say(strings.Repeat(emote, i))
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
	for i := count - 1; i > 0; i-- {
		say(strings.Repeat(emote, i))
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
}
