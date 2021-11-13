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
			fmt.Println(err)
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
	return
}

func (action BotActions) sendChannelMessage(cID string, content string) *discordgo.Message {
	action.checkSession()

	msg, _ := bot.session.ChannelMessageSend(cID, content)

	return msg
}

func (action BotActions) deleteChannelMessages(cID, message string) bool {
	err := bot.session.ChannelMessageDelete(cID, message)

	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Deleted")
	return true
}
