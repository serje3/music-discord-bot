package main

import (
	"errors"
	"github.com/bigkevmcd/go-configparser"
	"os"
)

// config values
var (
	token                   string
	developerKey            string
	commandPrefix           string
	executeCommandMsgDelete = "false"
)

func (cfg *Config) init() {
	var err error
	cfg.createConfig()

	cfg.options = map[string]map[string]*string{
		"Credentials": {
			"TOKEN":         &token,
			"DEVELOPER_KEY": &developerKey,
		},
		"Commands": {
			"COMMAND_PREFIX":             &commandPrefix,
			"COMMAND_EXECUTE_MSG_DELETE": &executeCommandMsgDelete,
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
	var err error
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
	var err error
	// pass the error, section maybe already created.
	_ = cfg.parser.AddSection(sectionName)
	err = cfg.parser.Set(sectionName, optionName, optionValue)
	if err != nil {
		FatalError("Cannot write in file config.cfg default data: ", err)
	}
}

func (cfg *Config) setValuesToVariables() {
	var err error
	for sectionName, optionMap := range cfg.options {
		for optionName, varPointer := range optionMap {
			*varPointer, err = cfg.parser.Get(sectionName, optionName)
			SimpleFatalErrorHandler(err)
		}
	}
}
