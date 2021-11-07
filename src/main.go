package main

import (
	"errors"
	"flag"
	"github.com/bigkevmcd/go-configparser"
	"log"
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
}

func main() {
	if token != "" {
		bot.DiscordConnect()
	} else {
		log.Fatal("Discord api token is not provided")
	}

	defer CloseLogFile()
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
		_ = config.AddSection("Credentials")
		err = config.Set("Credentials", "token", "<INSERT HERE YOUR DISCORD BOT API TOKEN>")
		if err != nil {
			log.Fatal("Cannot write in file config.cfg default data: ", err)
		}
		log.Println("Enter token in config.cfg")
	}
	err = config.SaveWithDelimiter("config.cfg", "=")
	SimpleFatalErrorHandler(err)

}
