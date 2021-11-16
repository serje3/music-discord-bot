package main

import (
	"github.com/AvraamMavridis/randomcolor"
	"github.com/bwmarrin/discordgo"
)

func helpEmbed() *discordgo.MessageEmbed {
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
			false,
		},
		{
			"[Commands] Play music to voice channel",
			"!play <youtube url | search query>",
			false,
		},
		{
			"[Commands] Skip song from queue",
			"!skip <none>",
			false,
		},
		{
			"[Commands] Remove songs from queue",
			"!clear <none>",
			false,
		},
		{
			"[Commands] Stop music player",
			"!stop <none>",
			false,
		},
	}
	embed.Fields = fields
	embed.Color = 0x3dea1a
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    "Dev: serje322#4196; Github: serje3",
		IconURL: "https://cdn.discordapp.com/avatars/263430624080035841/99b51ce89e05651f82910e13bec8e2b0.png",
	}
	return embed
}

func songEmbed(video YoutubeVideoDetails) *discordgo.MessageEmbed {
	var embed *discordgo.MessageEmbed

	embed = &discordgo.MessageEmbed{}
	embed.Title = video.Name
	embed.URL = youtubeVideoUrlPattern + video.ID
	embed.Description = video.Description
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{}
	embed.Thumbnail.URL = video.Thumbnail
	embed.Color = getRandomColor()
	return embed
}

func getRandomColor() int {
	randomColor := randomcolor.GetRandomColorInRgb()
	r := randomColor.Red
	g := randomColor.Green
	b := randomColor.Blue
	rgb := b | (g << 8) | (r << 16)
	return (0x1000000 | rgb) / 2
}
