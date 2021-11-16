package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/serje3/dgvoice"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const (
	channels  int = 2     // 1 for mono, 2 for stereo
	frameRate int = 48000 // audio sampling rate
	frameSize int = 960   // uint16 size of each audio frame
)

func PlayAudioFile(v *discordgo.VoiceConnection, filename string, guild *GuildVars) {

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := v.Speaking(false)
		fmt.Println("defer PlayAudioFile: speaking: ", *guild.speaking)
		if err != nil {
			fmt.Println("dgvoice: "+"Couldn't stop speaking", err)
		}
		guild.skipSong <- false
	}()

	// Create a shell command "object" to run.
	run := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegout, err := run.StdoutPipe()
	if err != nil {
		fmt.Println("dgvoice: "+"StdoutPipe Error", err)
		return
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	// Starts the ffmpeg command
	err = run.Start()
	if err != nil {
		fmt.Println("dgvoice: "+"RunStart Error", err)
		return
	}

	// prevent memory leak from residual ffmpeg streams
	defer func(Process *os.Process) {
		err := Process.Kill()
		if err != nil {
			log.Println(err)
		}
	}(run.Process)

	//when stop is sent, kill ffmpeg
	go func() {
		<-guild.stopMusic
		err = run.Process.Kill()
	}()

	// Send "speaking" packet over the voice websocket
	err = v.Speaking(true)
	fmt.Println("PlayAudioFile: speaking: ", *guild.speaking)
	if err != nil {
		fmt.Println("dgvoice: "+"Couldn't set speaking", err)
	}

	send := make(chan []int16, 2)
	defer close(send)

	closed := make(chan bool)
	go func() {
		dgvoice.SendPCM(v, send)
		fmt.Println("closing")
		closed <- true
		fmt.Println("closed 1")
	}()

	for {
		// read data from ffmpeg stdout
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			fmt.Println("dgvoice: "+"error reading from ffmpeg stdout", err)
			return
		}
		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-closed:
			fmt.Println("closed 2")
			return
		}
	}
}

func PlayQueue(voiceConnection *discordgo.VoiceConnection, guildInfo *GuildVars) {
	log.Println("PlayQueue: start")
	*guildInfo.speaking = true
	for length := guildInfo.queue.Len(); length != 0; length = guildInfo.queue.Len() {
		song := guildInfo.queue.Pop()
		fmt.Println("Len(): ", length, "Name: ", song.details.Name)
		go PlayAudioFile(voiceConnection, song.stream, guildInfo)

		isSkipFromUser := <-guildInfo.skipSong
		fmt.Println("isFromUser: ", isSkipFromUser)
		if isSkipFromUser {
			guildInfo.stopMusic <- true
			// receiving garbage from defer method in PlayAudioFile
			<-guildInfo.skipSong
		}
	}
	*guildInfo.speaking = false
	log.Println("PlayQueue: end")
}
