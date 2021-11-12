package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"os"
)

var token = ""

var logfile *os.File

func init() {
	// only in production
	SetLogOutputToFile()
	setTokenFromConfig()

	// only in development
	// setTokenFromFlag()

	youtubeClient.init()

	err := os.Mkdir(AUDIO_FOLDER, os.ModePerm)
	if !os.IsExist(err) {
		SimpleFatalErrorHandler(err)
	}
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

func setTokenFromFlag() {
	flag.StringVar(&token, "t", "", "Insert your Discord bot token here")
	flag.Parse()
}

func setTokenFromConfig() {
	if _, err := os.Stat("config.cfg"); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create("config.cfg")
		SimpleFatalErrorHandler(err)

		err = file.Close()
		SimpleFatalErrorHandler(err)
	}

	config, err := configparser.Parse("config.cfg")
	SimpleFatalErrorHandler(err)

	token, err = config.Get("Credentials", "token")
	if err != nil {
		setTokenFromStdin()
		// pass the error, section maybe already created.
		_ = config.AddSection("Credentials")
		err = config.Set("Credentials", "token", token)
		if err != nil {
			FatalError("Cannot write in file config.cfg default data: ", err)
		}

	}
	err = config.SaveWithDelimiter("config.cfg", "=")
	SimpleFatalErrorHandler(err)

}

func setTokenFromStdin() {
	fmt.Print("Enter token: ")
	_, err = fmt.Scanf("%s\n", &token)
}
