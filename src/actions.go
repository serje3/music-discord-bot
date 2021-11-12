package main

import (
	"errors"
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

func (action BotActions) sendChannelMessage(cID string, content string) {
	action.checkSession()

	_, _ = bot.session.ChannelMessageSend(cID, content)
	//nothing to do
}
