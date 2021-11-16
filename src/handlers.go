package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

func (bot Bot) DiscordAddHandlers() {
	bot.session.AddHandler(Ready)
	bot.session.AddHandler(messageCreate)
	bot.session.AddHandler(GuildCreate)
	bot.session.AddHandler(GuildDelete)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, commandPrefix) {
		DiscordExecuteCommand(m.Content[1:], s, m)
	}
}

func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	fmt.Println("Ready event called")

	// bad idea, but the only one...
	go func() {
		for {
			time.Sleep(60e+9)
			err := s.UpdateListeningStatus(fmt.Sprintf("%v guilds", guildsCount))
			if err != nil {
				log.Println(err)
				fmt.Println(err)
				return
			}
		}
	}()
}

func GuildCreate(_ *discordgo.Session, g *discordgo.GuildCreate) {
	guildsInfo[g.ID] = GuildVars{
		speaking:  func() *bool { b := false; return &b }(),
		stopMusic: make(chan bool),
		queue: &SongsQueue{
			make([]Song, 0),
		},
		skipSong: make(chan bool),
	}
	guildsCount++
}

func GuildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
	log.Println("Deleted", g.ID)
	guildsCount--
	err := s.UpdateListeningStatus(fmt.Sprintf("%v guilds", guildsCount))
	if err != nil {
		log.Println(err)
		return
	}
}
