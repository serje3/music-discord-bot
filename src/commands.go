package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/serje3/dgvoice"
	"reflect"
	"strings"
)

type Commands struct {
	errors CommandErrors
	utils  CommandUtils
}

type commandArgs []string

var commands Commands

type GuildVars struct {
	stopMusic chan bool
}

var guildsInfo map[string]GuildVars

func DiscordExecuteCommand(commandArgs string, s *discordgo.Session, m *discordgo.MessageCreate) {
	arrayArgs := strings.Split(commandArgs, " ")
	commandName := strings.Title(arrayArgs[0])
	method := reflect.ValueOf(&commands).MethodByName(commandName)

	if method.IsValid() {
		fmt.Println(method.String())
		go func() {
			method.Call(
				[]reflect.Value{
					reflect.ValueOf(s),
					reflect.ValueOf(m),
					reflect.ValueOf(arrayArgs[1:]),
				})
			bot.actions.deleteChannelMessages(m.ChannelID, m.ID)
		}()
	} else {
		bot.actions.sendChannelMessage(m.ChannelID, "[**Ошибка**] Такой команды нет")
	}
}

// commands list

func (command *Commands) Join(s *discordgo.Session, m *discordgo.MessageCreate, args commandArgs) {
	channel, err := command.utils.GetChannel(s, m, args)

	err = bot.actions.joinVoiceChannel(m.GuildID, channel.ID)

	if err != nil {
		bot.actions.sendChannelMessage(m.ChannelID, "Не удалось подключиться к голосовому каналу")
	} else {
		bot.actions.sendChannelMessage(m.ChannelID,
			fmt.Sprintf("Я присоединился к вашему каналу **%v**", channel.Name))
	}
}

func (command *Commands) Stop(s *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	go func() {
		guildsInfo[m.GuildID].stopMusic <- true
	}()

	if voiceConnection, ok := command.utils.GetVoiceConnection(s, m); ok && voiceConnection.Ready {
		err = bot.actions.quitVoiceChannel(m.GuildID)
	} else {
		err = errors.New("i cannot leave a channel to which you are not connected")
	}

	commandErrors.SimpleCommandErrorCheck(m.ChannelID, "Не получается:(", err)

}

func (command *Commands) Play(s *discordgo.Session, m *discordgo.MessageCreate, searchTextArgs commandArgs) {
	query := strings.Join(searchTextArgs, " ")

	voiceConnection, ok := command.utils.GetVoiceConnectionsOrJoin(s, m)
	if !ok {
		return
	}

	videoDetails, err := youtubeClient.GetVideoDetails(query)
	if commandErrors.SimpleCommandErrorCheck(m.ChannelID, "Не получилось получить информацию о видео", err) {
		return
	}

	url, err := videoDetails.GetAudioPath()
	if commandErrors.SimpleCommandErrorCheck(m.ChannelID, "Не получилось получить информацию о аудио", err) {
		return
	}
	bot.actions.sendChannelMessage(
		m.ChannelID,
		"[**Музыка**] "+videoDetails.Name+"\n"+youtubeVideoUrlPattern+videoDetails.ID,
	)
	dgvoice.PlayAudioFile(voiceConnection, url, guildsInfo[m.GuildID].stopMusic)
}
