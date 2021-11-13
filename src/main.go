package main

import (
	"fmt"
	"os"
)

var logfile *os.File

func init() {
	// only in production
	SetLogOutputToFile()
	config.init()
	// only in development
	// setTokenFromFlag()

	youtubeClient.init()
}

func main() {
	defer CloseLogFile()
	if token == "" {
		errorMsg := "Discord api token is not provided. Check config.cfg"
		// stdout
		fmt.Println(errorMsg)
		FatalError(errorMsg)
	} else {
		bot.DiscordConnect()
	}

}
