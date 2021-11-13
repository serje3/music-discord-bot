package main

import (
	"context"
	"errors"
	"fmt"
	YT "github.com/kkdai/youtube/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeAPI struct {
	ctx     context.Context
	service *youtube.Service
	client  YT.Client
}

type YoutubeVideoDetails struct {
	Name string
	ID   string
}

const youtubeVideoUrlPattern = "https://www.youtube.com/watch?v="

var youtubeClient YoutubeAPI

func (yt *YoutubeAPI) init() {
	yt.ctx = context.Background()
	yt.service, _ = youtube.NewService(yt.ctx, option.WithAPIKey(DEVELOPER_KEY))
	SimpleFatalErrorHandler(err)
	yt.client = YT.Client{}

}

func (yt YoutubeAPI) searchVideo(query string) (video YoutubeVideoDetails, err error) {
	call := yt.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(2)
	response, err := call.Do()
	if err != nil {
		return
	}

	video = YoutubeVideoDetails{
		Name: response.Items[0].Snippet.Title,
		ID:   response.Items[0].Id.VideoId,
	}
	return
}

func (yt *YoutubeAPI) GetVideo(url string) (*YT.Video, error) {
	video, err := yt.client.GetVideo(url)
	return video, err
}

func (yt *YoutubeAPI) GetStreamURL(link string) (string, error) {
	video, err := yt.GetVideo(link)
	if err != nil {
		// Handle the error
		return "", err
	}

	format := utils.findAudioFormat(video.Formats)
	streamURL, err := yt.client.GetStreamURL(video, format)
	if err != nil {
		fmt.Println("no stream url received")
		return streamURL, err
	}

	return streamURL, err
}

func (yt *YoutubeAPI) GetVideoDetails(query string) (YoutubeVideoDetails, error) {
	videoDetails, err := youtubeClient.searchVideo(query)
	if err != nil {
		fmt.Println("Fails here 1")
	} else if videoDetails.ID == "" {
		err = errors.New("got zero search results")
	}
	return videoDetails, err
}

func (videoDetails YoutubeVideoDetails) GetAudioPath() (url string, err error) {
	url = youtubeVideoUrlPattern + videoDetails.ID

	url, err = youtubeClient.GetStreamURL(url)

	return url, err
}
