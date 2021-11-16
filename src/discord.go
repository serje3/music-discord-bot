package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func (bot *Bot) DiscordConnect() {
	var err error
	bot.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return
	}

	bot.DiscordAddHandlers()
	bot.DiscordChangeIntents()
	bot.DiscordOpenWebsocket()
	defer bot.DiscordCloseWebsocket()

}

func (bot *Bot) DiscordChangeIntents() {
	bot.session.Identify.Intents = discordgo.IntentsAll
}

func (bot *Bot) DiscordOpenWebsocket() {
	err := bot.session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	suckcock := make(chan os.Signal, 1)
	signal.Notify(suckcock, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-suckcock
}

func (bot *Bot) DiscordCloseWebsocket() {
	err := bot.session.Close()
	if err != nil {
		log.Println(err)
		return
	}
}
