package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

type CommandUtils struct{}

// GetVoiceConnectionsOrJoin - thread-safe function that checks for the existence
// and returns a *discordgo.VoiceConnection, or if it does not exist,
// it executes the join function and checks again
func (cUtils CommandUtils) GetVoiceConnectionsOrJoin(m *discordgo.MessageCreate) (*discordgo.VoiceConnection, bool) {
	voiceConnection, ok := cUtils.GetVoiceConnection(m)
	if !ok {
		cUtils.JoinByVoiceState(m)
		if voiceConnection, ok = cUtils.GetVoiceConnection(m); !ok {
			bot.actions.sendChannelMessage(m.ChannelID, "Не получается:( Ты должен быть в голосовом канале.\nПопробуй команду !join")
			return nil, false
		}
	}
	return voiceConnection, ok
}

// GetVoiceConnection - thread-safe function  that returns a *discordgo.VoiceConnection
func (cUtils CommandUtils) GetVoiceConnection(m *discordgo.MessageCreate) (*discordgo.VoiceConnection, bool) {
	bot.session.RLock()
	voiceConnection, ok := bot.session.VoiceConnections[m.GuildID]
	bot.session.RUnlock()

	return voiceConnection, ok
}

// JoinByVoiceState joins to voice chat with no parameters(except event). Joins to current voice state
func (cUtils CommandUtils) JoinByVoiceState(m *discordgo.MessageCreate) bool {
	if voiceState, err := bot.session.State.VoiceState(m.GuildID, m.Author.ID); err == nil {
		channel, err := bot.session.Channel(voiceState.ChannelID)

		if commands.errors.SimpleCommandErrorCheck(m.ChannelID, "Голосовой канал не найден", err) {
			return false
		}

		if _join(m.GuildID, m.ChannelID, channel.Name) {
			return true
		}

	}

	return false
}

// JoinByChannelName joins to voice chat with channel name parameter
func (cUtils CommandUtils) JoinByChannelName(m *discordgo.MessageCreate, channelName string) bool {
	channel, err := utils.GetChannelByName(m.GuildID, channelName)

	if commands.errors.SimpleCommandErrorCheck(m.ChannelID, "Голосовой канал не найден", err) {
		return false
	}

	if _join(m.GuildID, m.ChannelID, channel.Name) {
		return true
	}

	return false
}

func _join(gID, cID, cName string) bool {

	err = bot.actions.joinVoiceChannel(gID, cID)
	log.Println(cID, err)
	if commands.errors.SimpleCommandErrorCheck(cID, "Не удалось подключиться к голосовому каналу ", err) {
		return false
	}

	bot.actions.sendChannelMessage(cID, fmt.Sprintf("Я присоединился к вашему каналу **%v**", cName))

	return true
}
