package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func (action BotActions) checkSession() {
	// idk if it's right or not, but just in case
	if bot.session == nil {
		log.Printf("bot not active.... %v", bot.session)
		SimpleFatalErrorHandler(errors.New("bot session is nil"))
	}
}

func (action BotActions) joinVoiceChannel(gID, cID string) (err error) {
	action.checkSession()

	_, err = bot.session.ChannelVoiceJoin(gID, cID, false, false)

	if err != nil {
		if _, ok := bot.session.VoiceConnections[gID]; ok {
			log.Println(err)
			_ = bot.session.VoiceConnections[gID]
		} else {
			log.Println(err)
			return err
		}
	}

	return

}

func (action BotActions) quitVoiceChannel(gID string) (err error) {
	action.checkSession()
	_, err = bot.session.ChannelVoiceJoin(gID, "", false, false)
	if err != nil {
		log.Println(err)
	}
	return
}

func (action BotActions) sendChannelMessage(cID string, content string) *discordgo.Message {
	action.checkSession()

	msg, err := bot.session.ChannelMessageSend(cID, content)

	if err != nil {
		log.Println(err)
	}

	return msg
}

func (action BotActions) sendChannelMessageEmbed(cID string, embed *discordgo.MessageEmbed) *discordgo.Message {
	action.checkSession()

	msg, _ := bot.session.ChannelMessageSendEmbed(cID, embed)

	return msg
}

func (action BotActions) deleteChannelMessages(cID, message string) bool {
	err := bot.session.ChannelMessageDelete(cID, message)
	if err != nil {
		log.Println(err)
		return false
	}
	fmt.Println("Deleted msg id:")
	return true
}
