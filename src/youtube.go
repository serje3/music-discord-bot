package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jonas747/dca"
	YT "github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
)

type Downloader downloader.Downloader

type YoutubeAPI struct {
	ctx     context.Context
	service *youtube.Service
	client  YT.Client
}

type YoutubeVideoDetails struct {
	Name string
	ID   string
}

const defaultExtension = ".mp3"

//.aac	audio/aac
//.mp3	audio/mpeg
//.oga	audio/ogg
//.opus	audio/opus
//.wav	audio/wav
//.weba	audio/webm
var canonicals = map[string]string{
	"audio/mp4":  ".m4a",
	"audio/aac":  ".aac",
	"audio/mpeg": ".mp3",
	"audio/ogg":  ".oga",
	"audio/opus": ".wav",
	"audio/webm": ".weba",
}

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

func (yt *YoutubeAPI) DownloadAudio(url string) (string, error) {
	v, err := yt.GetVideo(url)

	audioFormat := findAudioFormat(v.Formats)
	if audioFormat == nil {
		return "", errors.New("audio format not found")
	}

	stream, err := yt.client.GetStreamURL(v, audioFormat)
	log.Printf("Stream url: %s", stream)

	log.Printf("Title '%s' - Audio Codec '%s'", v.Title, audioFormat.MimeType)

	destFile, err := yt.getOutputFile(v, audioFormat)
	if err != nil {
		return "", err
	}

	// Create audio file
	audioFile, err := os.Create(destFile)
	if err != nil {
		return "", err
	}

	log.Printf("Downloading audio file...")
	err = yt.videoDLWorker(audioFile, v, audioFormat)
	if err != nil {
		return "", err
	}

	return destFile, err
}

func (yt *YoutubeAPI) getOutputFile(v *YT.Video, format *YT.Format) (string, error) {
	outputFile := downloader.SanitizeFilename(v.ID)
	outputFile += pickIdealFileExtension(format.MimeType)

	if AUDIO_FOLDER != "" {
		if err := os.MkdirAll(AUDIO_FOLDER, 0o755); err != nil {
			return "", err
		}
		outputFile = filepath.Join(AUDIO_FOLDER, outputFile)
	}

	return outputFile, nil
}

func pickIdealFileExtension(mediaType string) string {
	mediaType, _, err = mime.ParseMediaType(mediaType)
	if err != nil {
		return defaultExtension
	}

	if extension, ok := canonicals[mediaType]; ok {
		return extension
	}

	// Our last resort is to ask the operating system, but these give multiple results and are rarely canonical.
	extensions, err := mime.ExtensionsByType(mediaType)
	if err != nil || extensions == nil {
		return defaultExtension
	}

	return extensions[0]
}

func (yt *YoutubeAPI) videoDLWorker(out *os.File, video *YT.Video, format *YT.Format) error {
	stream, _, err := yt.client.GetStreamContext(yt.ctx, video, format)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, stream)
	if err != nil {
		return err
	}

	return nil
}

func (yt *YoutubeAPI) StreamAudioCreate(link string) (string, *dca.EncodeOptions, error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	video, err := yt.GetVideo(link)
	if err != nil {
		// Handle the error
		return "", options, err
	}

	format := findAudioFormat(video.Formats)
	streamURL, err := yt.client.GetStreamURL(video, format)
	if err != nil {
		fmt.Println("no stream url received")
		return streamURL, options, err
	}

	return streamURL, options, err
}
