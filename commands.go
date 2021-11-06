package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"reflect"
	"strings"
)

type Commands struct{}

var commands Commands

func DiscordExecuteCommand(args string, s *discordgo.Session, m *discordgo.MessageCreate) {
	arrayArgs := strings.Split(args, " ")
	commandName := strings.Title(arrayArgs[0])
	method := reflect.ValueOf(&commands).MethodByName(commandName)

	if method.IsValid() {
		method.Call(
			[]reflect.Value{
				reflect.ValueOf(s),
				reflect.ValueOf(m),
			})
	} else {
		s.ChannelMessageSend(m.ChannelID, "[**Ошибка**] Такой команды нет")
	}
}

// commands list

func (command Commands) Join(s *discordgo.Session, m *discordgo.MessageCreate) {
	voiceState, err := s.State.VoiceState(m.GuildID, m.Author.ID)

	if err != nil {
		log.Println("Voice channel not found ", err)
		return
	}

	isOk := bot.actions.joinVoiceChannel(m.GuildID, voiceState.ChannelID)

	if !isOk {
		s.ChannelMessageSend(m.ChannelID, "Не удалось подключиться к голосовому чату")
	} else {
		channel, err := s.Channel(voiceState.ChannelID)
		if err != nil {
			log.Println("Cannot receive channel ", err)
		}
		s.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("Я присоединился к вашему каналу **%v**", channel.Name))
	}
}
