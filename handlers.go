package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func (bot Bot) DiscordAddHandlers() {
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
