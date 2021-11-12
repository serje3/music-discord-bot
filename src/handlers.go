package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func (bot Bot) DiscordAddHandlers() {
	bot.session.AddHandler(Ready)
	bot.session.AddHandler(messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, COMMAND_PREFIX) {
		DiscordExecuteCommand(m.Content[1:], s, m)
	}
}

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println("Ready event called")
	guildsInfo = make(map[string]GuildVars)
	for _, guild := range r.Guilds {
		guildsInfo[guild.ID] = GuildVars{}
	}
}
