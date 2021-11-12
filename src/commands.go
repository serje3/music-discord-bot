package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"io"
	"log"
	"reflect"
	"strings"
	"time"
)

type Commands struct{}

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

func (command Commands) Join(s *discordgo.Session, m *discordgo.MessageCreate, args commandArgs) {
	var channel *discordgo.Channel

	if len(args) > 0 {
		channel, err = utils.GetChannelByName(m.GuildID, strings.Join(args, " "))
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

func (command Commands) Stop(_ *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	var err error
	guildsInfo[m.GuildID].stopMusic <- true

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

func (command Commands) Play(s *discordgo.Session, m *discordgo.MessageCreate, searchTextArgs commandArgs) {
	voiceConnection, ok := bot.session.VoiceConnections[m.GuildID]
	if !ok {
		command.Join(s, m, searchTextArgs)
		if voiceConnection, ok = bot.session.VoiceConnections[m.GuildID]; !ok {
			bot.actions.sendChannelMessage(m.ChannelID, "Не получается:( Ты должен быть в голосовом канале")
			return
		}
	}

	url, err := utils.GetAudioURL(strings.Join(searchTextArgs, " "), false)
	if err != nil {
		bot.actions.sendChannelMessage(m.ChannelID, err.Error())
		return
	}

	streamUrl, options, err := youtubeClient.StreamAudioCreate(url)
	if err != nil {
		bot.actions.sendChannelMessage(m.ChannelID, err.Error())
		return
	}
	PlayAudioStream(voiceConnection, streamUrl, options, guildsInfo[m.GuildID].stopMusic)
	//PlayAudioFile(voiceConnection, url, guildsInfo[m.GuildID].stopMusic)
	bot.actions.sendChannelMessage(m.ChannelID, "работает???")

}

func PlayAudioStream(voiceConnection *discordgo.VoiceConnection, url string, options *dca.EncodeOptions, music chan bool) {
	log.Println(url)
	encodingSession, err := dca.EncodeFile("audio/yAF9XlluONA.webm", options)

	if err != nil {
		// Handle the error
		fmt.Println(err)
		return
	}
	defer encodingSession.Cleanup()

	done := make(chan error)
	streamingSessiong := dca.NewStream(encodingSession, voiceConnection, done)
	defer func() {
		finished, err := streamingSessiong.Finished()
		if err != nil {
			return
		}
		fmt.Println("ЗАКОНЧИЛОСЬ - ", finished)
	}()

	streamingSessiong.SetPaused(false)
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				log.Fatal("An error occured", err)
			}

			// Clean up incase something happened and ffmpeg is still running
			encodingSession.Truncate()
			return
		case <-ticker.C:
			stats := encodingSession.Stats()
			playbackPosition := streamingSessiong.PlaybackPosition()

			fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
		}
	}
}
