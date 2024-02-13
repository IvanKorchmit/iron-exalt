package main

import (
	"errors"
	"fmt"

	"github.com/IvanKorchmit/akevitt"
	"github.com/uaraven/logview"
)

func onMessage(engine *akevitt.Akevitt, session *akevitt.ActiveSession, channel, message, username string) error {
	if session == nil {
		return errors.New("session is nil. Probably the dead one")
	}

	chat, ok := session.Data["chat"].(*logview.LogView)

	if !ok {
		return errors.New("chatlog is not logview.LogView")
	}

	st := fmt.Sprintf("%s (%s): %s", username, channel, message)

	akevitt.AppendText(st, chat)

	return nil
}
