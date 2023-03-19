package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
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

type OpenAIChoice struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}

type OpenAIResponse struct {
	Id      string         `json:"id"`
	Object  string         `json:"object"`
	Created int            `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
}

var guildsInfo = make(map[string]GuildVars)

func DiscordExecuteCommand(commandArgs string, s *discordgo.Session, m *discordgo.MessageCreate) {
	arrayArgs := strings.Split(commandArgs, " ")
	commandName := strings.Title(arrayArgs[0])
	method := reflect.ValueOf(&commands).MethodByName(commandName)

	if method.IsValid() {
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
	if channel == nil {
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

func (command *Commands) Stop(s *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	go func() {
		guildsInfo[m.GuildID].stopMusic <- true
	}()

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
	log.Println(url, voiceConnection)

	//go dgvoice.PlayAudioFile(voiceConnection, url, guildsInfo[m.GuildID].stopMusic)
}

func (command *Commands) Help(s *discordgo.Session, m *discordgo.MessageCreate, _ commandArgs) {
	var embed *discordgo.MessageEmbed
	var fields []*discordgo.MessageEmbedField
	embed = &discordgo.MessageEmbed{}
	embed.Title = "Help instruction"
	embed.Description = "I will help you understand how to tolerate me"
	embed.Author = &discordgo.MessageEmbedAuthor{}
	embed.Author.Name = "Captain Gopus"
	embed.Author.URL = "https://github.com/serje3/music-discord-bot"
	embed.Author.IconURL = "https://i.ibb.co/HVnMfxc/7-Uu-Wz-LWn-HZA.jpg"
	fields = []*discordgo.MessageEmbedField{
		{
			"[Commands] Join to voice channel",
			"!join <none | channel name>",
			true,
		},
		{
			"[Commands] Play music to voice channel",
			"!play <youtube url | search query>",
			true,
		},
		{
			"[Commands] Stop music player",
			"!stop <none>",
			true,
		},
	}
	embed.Fields = fields
	embed.Color = 0x3dea1a

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    "Dev: serje322#4196; Github: serje3",
		IconURL: "https://cdn.discordapp.com/avatars/263430624080035841/99b51ce89e05651f82910e13bec8e2b0.png",
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		fmt.Println(err)
	}
}

func (command *Commands) Ask(s *discordgo.Session, m *discordgo.MessageCreate, question commandArgs) {

	textQuestion := strings.Join(question, " ")

	log.Println(textQuestion)

	url := "https://api.openai.com/v1/completions"

	// Define the request body
	//body := json.Marshal({
	//	"model": "text-davinci-003",
	//	"prompt": %s,
	//	"max_tokens": 7,
	//	"temperature": 0
	//}, textQuestion)))

	tokens, _ := strconv.Atoi(openAIMaxLength)
	temp, _ := strconv.Atoi(openAITemp)
	log.Println(openAITemp, openAIModel, openAIMaxLength)
	body := map[string]interface{}{
		"model":       openAIModel,
		"prompt":      textQuestion,
		"max_tokens":  tokens,
		"temperature": temp,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Printf("could not marshal json: %s\n", err)
		return
	}

	// Create a new request with the body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Add the API key to the request headers
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", openAIToken))
	req.Header.Add("Content-Type", "application/json")

	// Send the request and get the response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var answer OpenAIResponse

	err = json.Unmarshal(responseBody, &answer)

	if err != nil {
		log.Println(err)
		bot.actions.sendChannelMessage(m.ChannelID, "Упс произошла ошибочька.....")
		return
	}

	if len(answer.Choices) == 0 {
		bot.actions.sendChannelMessage(m.ChannelID, "Хз че ответить).....")
		return
	}
	content := answer.Choices[0].Text

	bot.actions.sendChannelMessage(m.ChannelID, content)
}
