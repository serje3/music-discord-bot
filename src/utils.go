package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
	"io"
	"os"
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

func (_ Utils) GetAudioURL(searchText string, stream bool) (url string, err error) {
	videoDetails, err := youtubeClient.searchVideo(searchText)
	if err != nil {
		fmt.Println("Fails here 1")
		return
	} else if videoDetails.ID == "" {
		return url, errors.New("got zero search results")
	}
	url = "https://www.youtube.com/watch?v=" + videoDetails.ID
	if !stream {
		url, err = youtubeClient.DownloadAudio(url)
	}
	return url, err
}

func findAudioFormat(formats youtube.FormatList) *youtube.Format {
	var audioFormat *youtube.Format
	var audioFormats youtube.FormatList

	audioFormats = formats.Type("audio")

	if len(audioFormats) > 0 {
		audioFormats.Sort()
		audioFormat = &audioFormats[0]
	}

	return audioFormat
}

func saveStream(stream *io.ReadCloser, filename string) (err error) {
	fmt.Println("1")

	file, err := os.Create(filename)
	fmt.Println("2", err)

	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("3")
	_, err = io.Copy(file, *stream)
	fmt.Println("4")
	if err != nil {
		return err
	}
	fmt.Println("5")

	fmt.Println("ended")
	return
}
