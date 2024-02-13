package main

import (
	"log"
	"strings"

	"github.com/IvanKorchmit/akevitt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

func main() {
	customRoom := akevitt.Room{Name: "My room"}
	engine := akevitt.NewEngine().
		UseRootUI(root).
		UseSpawnRoom(&customRoom).
		UseRegisterCommand("say", akevitt.OocCmd).
		UseDBPath("database.db").
		UseOnMessage(onMessage).
		UseBind(":2222").
		Finish()
	log.Fatal(engine.Run())
}

func root(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	modal := tview.NewModal().
		AddButtons([]string{"Login", "Register"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Register" {
				session.Application.SetRoot(akevitt.RegistrationScreen(engine, session, gameRoom), true)
			} else if buttonLabel == "Login" {
				session.Application.SetRoot(akevitt.LoginScreen(engine, session, gameRoom), true)
			}
		})

	modal.SetTitle("Welcome to Iron Exalt!")

	return modal
}

func gameRoom(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	playerMessage := ""

	// Preparing session by initializing UI primitives, channels and collections.
	chatlog := logview.NewLogView()
	chatlog.SetLevelHighlighting(true)
	session.Data["chat"] = chatlog

	inputField := tview.NewInputField().
		SetChangedFunc(func(text string) {
			playerMessage = text
		})

	game := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		AddItem(inputField, 2, 0, 1, 3, 0, 0, true).
		AddItem(chatlog, 1, 1, 1, 2, 0, 0, false).
		SetBorders(true)
	inputField.SetFinishedFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}
		playerMessage = strings.TrimSpace(playerMessage)
		if playerMessage == "" {
			inputField.SetText("")
			return
		}

		akevitt.AppendText("\t>"+playerMessage, chatlog)
		err := engine.ExecuteCommand(playerMessage, session)
		if err != nil {
			akevitt.ErrorBox(err.Error(), session.Application, game)
			inputField.SetText("")
			return
		}
		playerMessage = ""
		inputField.SetText("")
	})
	return game
}
