package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
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

func (_ Utils) findAudioFormat(formats youtube.FormatList) *youtube.Format {
	var audioFormat *youtube.Format
	var audioFormats youtube.FormatList

	audioFormats = formats.Type("audio")

	if len(audioFormats) > 0 {
		audioFormats.Sort()
		audioFormat = &audioFormats[0]
	}

	return audioFormat
}
