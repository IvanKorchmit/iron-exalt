package main

import (
	"log"
	"strings"
	"unicode"

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
				session.Application.SetRoot(register(engine, session), true)
			} else if buttonLabel == "Login" {
				session.Application.SetRoot(login(engine, session), true)
			}
		})

	modal.SetTitle("Welcome to Iron Exalt!")

	return modal
}

func register(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	username := ""
	password := ""
	repeatPassword := ""

	form := tview.NewForm()

	form.AddInputField("Username", "", 0, func(textToCheck string, lastChar rune) bool {
		if !unicode.IsLetter(lastChar) && !unicode.IsDigit(lastChar) || lastChar > unicode.MaxASCII {
			return false
		}

		username = textToCheck
		return true
	}, nil).
		AddPasswordField("Repeat password", "", 0, '*', func(text string) {
			password = text
		}).
		AddPasswordField("Repeat password", "", 0, '*', func(text string) {
			repeatPassword = text
		}).
		AddButton("Register", func() {
			err := engine.Register(username, password, repeatPassword, session)

			if err != nil {
				akevitt.ErrorBox(err.Error(), session.Application, form)
				return
			}

			session.Application.SetRoot(gameRoom(engine, session), true)
		})

	form.SetTitle("Registration")

	return form
}

func login(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	username := ""
	password := ""

	form := tview.NewForm()

	form.AddInputField("Username", "", 0, func(textToCheck string, lastChar rune) bool {
		if !unicode.IsLetter(lastChar) && !unicode.IsDigit(lastChar) || lastChar > unicode.MaxASCII {
			return false
		}

		username = textToCheck
		return true
	}, nil).
		AddPasswordField("Repeat password", "", 0, '*', func(text string) {
			password = text
		}).
		AddButton("Register", func() {
			err := engine.Login(username, password, session)

			if err != nil {
				akevitt.ErrorBox(err.Error(), session.Application, form)
				return
			}

			session.Application.SetRoot(gameRoom(engine, session), true)
		})

	form.SetTitle("Registration")

	return form
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
