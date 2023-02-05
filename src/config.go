package main

import (
	"errors"
	"github.com/bigkevmcd/go-configparser"
	"os"
)

type Config struct {
	parser  *configparser.ConfigParser
	options map[string]map[string]*string
}

var config Config

// config values
var (
	token                   string
	developerKey            string
	openAIToken             string
	openAIMaxLength         string
	openAIModel             string
	openAITemp              string
	commandPrefix           string
	executeCommandMsgDelete = "false"
)

func (cfg *Config) init() {
	cfg.createConfig()

	cfg.options = map[string]map[string]*string{
		"Credentials": {
			"TOKEN":         &token,
			"DEVELOPER_KEY": &developerKey,
			"CHATGPT_TOKEN": &openAIToken,
		},
		"Commands": {
			"COMMAND_PREFIX":             &commandPrefix,
			"COMMAND_EXECUTE_MSG_DELETE": &executeCommandMsgDelete,
			"CHATGPT_MAX_TOKENS":         &openAIMaxLength,
			"CHATGPT_MODEL":              &openAIModel,
			"CHATGPT_TEMPERATURE":        &openAITemp,
		},
	}

	cfg.parser, err = configparser.Parse("config.cfg")
	SimpleFatalErrorHandler(err)
	cfg.createDefaultOptionsIfNotExist()
	cfg.setValuesToVariables()
}

func (cfg *Config) createConfig() {
	if _, err := os.Stat("config.cfg"); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create("config.cfg")
		SimpleFatalErrorHandler(err)

		err = file.Close()
		SimpleFatalErrorHandler(err)
	}
}

func (cfg *Config) createDefaultOptionsIfNotExist() {
	changed := false
	for section, optionMap := range cfg.options {
		for option := range optionMap {
			_, err = cfg.parser.Get(section, option)
			if err != nil {
				changed = true
				cfg.createOption(section, option, "<INSERT VALUE>")
			}
		}
	}
	if changed {
		err = cfg.parser.SaveWithDelimiter("config.cfg", "=")
		EndProgramWithMessage("missing options created in config.cfg")
	}
}

func (cfg *Config) createOption(sectionName, optionName, optionValue string) {
	// pass the error, section maybe already created.
	_ = cfg.parser.AddSection(sectionName)
	err = cfg.parser.Set(sectionName, optionName, optionValue)
	if err != nil {
		FatalError("Cannot write in file config.cfg default data: ", err)
	}
}

func (cfg *Config) setValuesToVariables() {
	for sectionName, optionMap := range cfg.options {
		for optionName, varPointer := range optionMap {
			*varPointer, err = cfg.parser.Get(sectionName, optionName)
			SimpleFatalErrorHandler(err)
		}
	}
}
