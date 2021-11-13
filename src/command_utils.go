package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type CommandUtils struct{}

// GetVoiceConnectionsOrJoin - thread-safe function that checks for the existence
// and returns a *discordgo.VoiceConnection, or if it does not exist,
// it executes the join function and checks again
func (cUtil CommandUtils) GetVoiceConnectionsOrJoin(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.VoiceConnection, bool) {
	voiceConnection, ok := cUtil.GetVoiceConnection(s, m)
	if !ok {
		commands.Join(s, m, []string{})
		if voiceConnection, ok = cUtil.GetVoiceConnection(s, m); !ok {
			bot.actions.sendChannelMessage(m.ChannelID, "Не получается:( Ты должен быть в голосовом канале.\nПопробуй команду !join")
			return nil, false
		}
	}
	return voiceConnection, ok
}

// GetVoiceConnection - thread-safe function  that returns a *discordgo.VoiceConnection
func (cUtil CommandUtils) GetVoiceConnection(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.VoiceConnection, bool) {
	s.RLock()
	voiceConnection, ok := s.VoiceConnections[m.GuildID]
	s.RUnlock()
	return voiceConnection, ok
}

func (cUtil CommandUtils) GetChannel(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (*discordgo.Channel, error) {
	var channel *discordgo.Channel
	var err error

	if len(args) > 0 {
		channel, err = utils.GetChannelByName(m.GuildID, strings.Join(args, " "))
	} else {
		if voiceState, err := s.State.VoiceState(m.GuildID, m.Author.ID); err == nil {
			channel, err = s.Channel(voiceState.ChannelID)
		}
	}

	fmt.Println(err)

	if channel == nil {
		bot.actions.sendChannelMessage(m.ChannelID, "Голосовой канал не найден")
		return nil, err
	}

	return channel, err
}
