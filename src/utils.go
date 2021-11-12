package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

type Utils struct{}

var utils Utils

func (_ Utils) GetChannelByName(gID string, channelName string) (channel *discordgo.Channel, err error) {
	channels, err := bot.session.GuildChannels(gID)
	if err == nil {
		for _, _channel := range channels {
			if _channel.Name == channelName {
				channel = _channel
			}
		}
	} else {
		return
	}

	if channel != nil {
		return
	} else {
		return channel, errors.New("channel not found")
	}

}
