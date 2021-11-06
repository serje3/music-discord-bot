package main

import "flag"

var token string

func init() {
	flag.StringVar(&token, "t", "", "Insert your Discord bot token here")
	flag.Parse()
}

func main() {
	bot.DiscordConnect()
}
