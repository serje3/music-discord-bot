package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"reflect"
	"strings"
)

type Commands struct{}

type commandArgs []string

var commands Commands

func DiscordExecuteCommand(commandArgs string, s *discordgo.Session, m *discordgo.MessageCreate) {
	arrayArgs := strings.Split(commandArgs, " ")
	commandName := strings.Title(arrayArgs[0])
	method := reflect.ValueOf(&commands).MethodByName(commandName)

	if method.IsValid() {
		go method.Call(
			[]reflect.Value{
				reflect.ValueOf(s),
				reflect.ValueOf(m),
				reflect.ValueOf(arrayArgs[1:]),
			})
	} else {
		bot.actions.sendChannelMessage(m.ChannelID, "[**Ошибка**] Такой команды нет")
	}
}

// commands list

func (command *Commands) Join(s *discordgo.Session, m *discordgo.MessageCreate, args commandArgs) {
	var channel *discordgo.Channel

	if len(args) > 0 {
		channel, err = utils.GetChannelByName(m.GuildID, args[0])
	} else {
		if voiceState, err := s.State.VoiceState(m.GuildID, m.Author.ID); err == nil {
			channel, err = s.Channel(voiceState.ChannelID)
		}
	}

	if channel == nil || err != nil {
		bot.actions.sendChannelMessage(m.ChannelID, "Голосовой канал не найден")
		return
	}

	err = bot.actions.joinVoiceChannel(m.GuildID, channel.ID)

	if err != nil {
		bot.actions.sendChannelMessage(m.ChannelID, "Не удалось подключиться к голосовому каналу")
	} else {
		bot.actions.sendChannelMessage(m.ChannelID,
			fmt.Sprintf("Я присоединился к вашему каналу **%v**", channel.Name))
	}
}

func (command *Commands) Stop(_ *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	var err error
	bot.session.RLock()
	if voiceConnection, ok := bot.session.VoiceConnections[m.GuildID]; ok && voiceConnection.Ready {
		err = bot.actions.quitVoiceChannel(m.GuildID)
	} else {
		err = errors.New("cannot quit channel, which u not connected")
	}
	bot.session.RUnlock()

	if err != nil {
		bot.actions.sendChannelMessage(m.ChannelID, "Не получается:(")
		return
	}
}
