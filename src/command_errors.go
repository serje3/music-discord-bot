package main

import "log"

func (cError CommandErrors) SimpleCommandErrorCheck(cID string, msg string, err error) bool {
	if err != nil {
		log.Println(err.Error())
		bot.actions.sendChannelMessage(cID, msg)
		return true
	}
	return false
}
