package main

import (
	"context"
	"github.com/bigkevmcd/go-configparser"
	"github.com/bwmarrin/discordgo"
	YT "github.com/kkdai/youtube/v2"
	"google.golang.org/api/youtube/v3"
)

type GuildVars struct {
	speaking  *bool
	stopMusic chan bool
	queue     *SongsQueue
	skipSong  chan bool
}

type Commands struct {
	errors CommandErrors
	utils  CommandUtils
}

type CommandErrors struct{}

type CommandUtils struct{}

type Bot struct {
	session *discordgo.Session
	actions BotActions
}

type BotActions struct{}

type YoutubeAPI struct {
	ctx     context.Context
	service *youtube.Service
	client  YT.Client
}

type YoutubeVideoDetails struct {
	Name        string
	ID          string
	Thumbnail   string
	Description string
}

type Config struct {
	parser  *configparser.ConfigParser
	options map[string]map[string]*string
}

type Utils struct{}

type Song struct {
	details YoutubeVideoDetails
	stream  string
}

type SongsQueue struct {
	songs []Song
}
