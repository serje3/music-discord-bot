package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type commandArgs []string

func DiscordExecuteCommand(commandArgs string, s *discordgo.Session, m *discordgo.MessageCreate) {
	arrayArgs := strings.Split(commandArgs, " ")
	commandName := strings.Title(arrayArgs[0])
	method := reflect.ValueOf(&commands).MethodByName(commandName)

	if method.IsValid() {
		fmt.Println(commandName)
		go func() {
			method.Call(
				[]reflect.Value{
					reflect.ValueOf(s),
					reflect.ValueOf(m),
					reflect.ValueOf(arrayArgs[1:]),
				})
			if removeIt, err := strconv.ParseBool(executeCommandMsgDelete); removeIt {
				if err != nil {
					log.Println("Cannot parse COMMAND_EXECUTE_MSG_DELETE value. Its must be bool (True or False)")
				} else {
					bot.actions.deleteChannelMessages(m.ChannelID, m.ID)
				}
			}
		}()
	} else {
		bot.actions.sendChannelMessage(m.ChannelID, "[**Ошибка**] Такой команды нет")
	}
}

// commands list

func (command *Commands) Join(s *discordgo.Session, m *discordgo.MessageCreate, channelName commandArgs) {
	channel, err := command.utils.GetChannel(s, m, channelName)

	err = bot.actions.joinVoiceChannel(m.GuildID, channel.ID)

	if err != nil {
		bot.actions.sendChannelMessage(m.ChannelID, "Не удалось подключиться к голосовому каналу")
	} else {
		bot.actions.sendChannelMessage(m.ChannelID,
			fmt.Sprintf("Я присоединился к вашему каналу **%v**", channel.Name))
	}
}

func (command *Commands) Stop(s *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	guildInfo := guildsInfo[m.GuildID]
	var err error
	if *guildInfo.speaking {
		go func() { guildInfo.stopMusic <- true }()
	}

	if voiceConnection, ok := command.utils.GetVoiceConnection(s, m); ok && voiceConnection.Ready {
		err = voiceConnection.Speaking(false)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = bot.actions.quitVoiceChannel(m.GuildID)
	} else {
		err = errors.New("i cannot leave a channel to which you are not connected")
	}

	commandErrors.SimpleCommandErrorCheck(m.ChannelID, "Не получается:(", err)
}

func (command *Commands) Skip(_ *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	guildInfo := guildsInfo[m.GuildID]
	if *guildInfo.speaking {
		guildInfo.skipSong <- true
		bot.actions.sendChannelMessage(m.ChannelID, "Музыка попущена")
	} else {
		bot.actions.sendChannelMessage(m.ChannelID, "нет")
	}
}

func (command *Commands) Clear(_ *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	guildInfo := guildsInfo[m.GuildID]
	guildInfo.queue.Clear()
	bot.actions.sendChannelMessage(m.ChannelID, "Музыкальный ряд очищен от говна и ссанья")
}

func (command *Commands) Play(s *discordgo.Session, m *discordgo.MessageCreate, searchTextArgs commandArgs) {
	guildInfo := guildsInfo[m.GuildID]

	query := strings.Join(searchTextArgs, " ")

	voiceConnection, ok := command.utils.GetVoiceConnectionsOrJoin(s, m)
	if !ok {
		return
	}

	videoDetails, err := youtubeClient.GetVideoDetails(query)
	if commandErrors.SimpleCommandErrorCheck(m.ChannelID, "Не получилось получить информацию о видео", err) {
		return
	}

	url, err := youtubeClient.RequestAudioPath(videoDetails)
	fmt.Println("Play: stream url: ", url)
	song := Song{
		details: videoDetails,
		stream:  url,
	}

	if commandErrors.SimpleCommandErrorCheck(m.ChannelID, "Не получилось получить информацию о аудио", err) {
		return
	}

	guildInfo.queue.Push(song)
	bot.actions.sendChannelMessageEmbed(m.ChannelID, songEmbed(videoDetails))

	fmt.Println("Play: speaking: ", *guildInfo.speaking)
	if !*guildInfo.speaking {
		PlayQueue(voiceConnection, &guildInfo)
	}
}

func (command *Commands) Help(_ *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	bot.actions.sendChannelMessageEmbed(m.ChannelID, helpEmbed())
}
