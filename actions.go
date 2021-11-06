package main

import (
	"log"
)

func (action BotActions) joinVoiceChannel(gID, cID string) (isOk bool) {
	if bot.session == nil {
		log.Printf("bot not active.... %v", bot.session)
		return false
	}

	_, err := bot.session.ChannelVoiceJoin(gID, cID, false, false)

	if err != nil {
		if _, ok := bot.session.VoiceConnections[gID]; ok {
			_ = bot.session.VoiceConnections[gID]
		} else {
			log.Println(err)
			return false
		}
	}

	return true

}
