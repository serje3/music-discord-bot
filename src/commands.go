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
		method.Call(
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

func (command Commands) Join(_ *discordgo.Session, m *discordgo.MessageCreate, args commandArgs) {
	var ok bool
	if len(args) > 0 {
		ok = command.utils.JoinByChannelName(m, strings.Join(args, " "))
	} else {
		ok = command.utils.JoinByVoiceState(m)
	}

	fmt.Println(ok)

	//bot.actions.deleteChannelMessages(m.ChannelID, m.ID)
}

func (command Commands) Stop(_ *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	guildsInfo[m.GuildID].stopMusic <- true

	if voiceConnection, ok := command.utils.GetVoiceConnection(m); ok && voiceConnection.Ready {
		err = bot.actions.quitVoiceChannel(m.GuildID)
	} else {
		err = errors.New("i cannot leave a channel to which you are not connected")
	}

	commandErrors.SimpleCommandErrorCheck(m.ChannelID, "Не получается:(", err)
	//bot.actions.deleteChannelMessages(m.ChannelID, m.ID)
}

func (command Commands) Play(_ *discordgo.Session, m *discordgo.MessageCreate, searchTextArgs commandArgs) {
	query := strings.Join(searchTextArgs, " ")

	voiceConnection, ok := command.utils.GetVoiceConnectionsOrJoin(m)
	if !ok {
		return
	}

	videoDetails, err := youtubeClient.GetVideoDetails(query)
	if commandErrors.SimpleCommandErrorCheck(m.ChannelID, err.Error(), err) {
		return
	}

	url, err := videoDetails.GetAudioPath()
	if commandErrors.SimpleCommandErrorCheck(m.ChannelID, err.Error(), err) {
		return
	}
	bot.actions.sendChannelMessage(m.ChannelID, "[**Музыка**]"+videoDetails.Name+"\n"+youtubeVideoUrlPattern+videoDetails.ID)
	dgvoice.PlayAudioFile(voiceConnection, url, guildsInfo[m.GuildID].stopMusic)
	//bot.actions.deleteChannelMessages(m.ChannelID, m.ID)
}
