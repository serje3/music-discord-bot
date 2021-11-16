package main

import (
	"context"
	"errors"
	"fmt"
	YT "github.com/kkdai/youtube/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"strings"
)

const youtubeVideoUrlPattern = "https://www.youtube.com/watch?v="

func (yt *YoutubeAPI) init() {
	var err error
	yt.ctx = context.Background()
	yt.service, err = youtube.NewService(yt.ctx, option.WithAPIKey(developerKey))
	if err != nil {
		return
	}
	SimpleFatalErrorHandler(err)
	yt.client = YT.Client{}

}

func (yt YoutubeAPI) searchVideo(query string) (YoutubeVideoDetails, error) {
	// refactor later
	query = strings.TrimSpace(query)
	if strings.HasPrefix(youtubeVideoUrlPattern, query) {
		fmt.Println("legit")
		call := yt.service.Videos.List([]string{"id", "snippet"}).
			Id(query[len(youtubeVideoUrlPattern)-1:])
		response, err := call.Do()
		if err != nil {
			fmt.Println(err)
			return YoutubeVideoDetails{}, err
		}
		video := YoutubeVideoDetails{
			Name:        response.Items[0].Snippet.Title,
			ID:          response.Items[0].Id,
			Thumbnail:   response.Items[0].Snippet.Thumbnails.High.Url,
			Description: response.Items[0].Snippet.Description,
		}
		return video, err
	} else {
		call := yt.service.Search.List([]string{"id", "snippet"}).
			Q(query).
			MaxResults(5)
		response, err := call.Do()
		if err != nil {
			fmt.Println(err)
			return YoutubeVideoDetails{}, err
		}
		responseItem := yt.chooseVideoFromSearchList(response.Items)

		video := YoutubeVideoDetails{
			Name:        responseItem.Snippet.Title,
			ID:          responseItem.Id.VideoId,
			Thumbnail:   responseItem.Snippet.Thumbnails.High.Url,
			Description: responseItem.Snippet.Description,
		}

		return video, err
	}

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

func (yt *YoutubeAPI) RequestAudioPath(videoDetails YoutubeVideoDetails) (url string, err error) {
	url = youtubeVideoUrlPattern + videoDetails.ID

	url, err = youtubeClient.GetStreamURL(url)

	return url, err
}

func (yt *YoutubeAPI) chooseVideoFromSearchList(items []*youtube.SearchResult) *youtube.SearchResult {
	for _, item := range items {
		if item.Id.Kind == "youtube#video" {
			return item
		}
	}
	return items[0]
}
