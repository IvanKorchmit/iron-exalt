package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IvanKorchmit/akevitt"
	"github.com/IvanKorchmit/akevitt/plugins"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	gob.Register(&Character{})
	gob.Register(map[string]any{})

	customRoom := akevitt.Room{Name: "My room"}
	engine := akevitt.NewEngine().
		UseRootUI(root).
		UseSpawnRoom(&customRoom).
		AddPlugin(plugins.NewBoltPlugin[*akevitt.Account]("database.db")).
		AddPlugin(plugins.NewDefaultPlugins()).
		UseOnJoin(InitBundle).
		UseRegisterCommand("look", func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, arguments string) error {
			message, err := akevitt.FetchPlugin[*plugins.MessagePlugin](engine)
			if err != nil {
				return err
			}

			bundle, ok := session.Data[IronExaltBundle].(Bundle)

			if !ok {
				return errors.New("could not cast to bundle")
			}

			output := fmt.Sprintf("You're in %s", bundle.Character.room.Name)

			return (*message).Message(engine, bundle.Character.Name, output, "", session)
		}).
		UseRegisterCommand("boop", func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, arguments string) error {
			heartbeat, err := akevitt.FetchPlugin[*plugins.HeartBeatsPlugin](engine)
			if err != nil {
				return err
			}
			messaging, err := akevitt.FetchPlugin[*plugins.MessagePlugin](engine)
			if err != nil {
				return err
			}
			counter := 0
			return (*heartbeat).SubscribeToHeartBeat(time.Second, func() error {
				if counter >= 5 {
					return errors.New("done")
				}
				counter += 1
				return (*messaging).Message(engine, "ooc", "Beep "+fmt.Sprint(counter), session.Account.Username, session)
			})
		}).
		UseRegisterCommand("say", func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, arguments string) error {
			message, err := akevitt.FetchPlugin[*plugins.MessagePlugin](engine)

			if err != nil {
				return err
			}

			character := session.Account.PersistentData[IronExaltBundle].(*Bundle).Character

			return (*message).Message(engine, character.room.Name, arguments, character.Name, session)
		}).
		UseBind(":2222").
		Finish()
	log.Fatal(engine.Run())
}

func root(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	modal := tview.NewModal().
		AddButtons([]string{"Login", "Register"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Register" {
				session.Application.SetRoot(akevitt.RegistrationScreen(engine, session, characterWizard), true)
			} else if buttonLabel == "Login" {
				session.Application.SetRoot(akevitt.LoginScreen(engine, session, func(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
					character := Fetch[*Bundle](session.Data, IronExaltBundle)

					bundle := session.Data[IronExaltBundle].(Bundle)

					bundle.Character = *character
					session.Account.PersistentData["character"] = &bundle.Character
					r, _ := engine.GetRoom(bundle.Character.RoomKey)
					bundle.Character.UpdateRoom(r)
					return gameScreen(engine, session)
				}), true)
			}
		})

	modal.SetTitle("Welcome to Iron Exalt!")

	return modal
}

func characterWizard(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	form := tview.NewForm()

	characterName := ""

	form.AddInputField("Character's name", "", 0, nil, func(text string) {
		characterName = text
	})
	form.AddButton("Create & Go", func() {
		if strings.TrimSpace(characterName) == "" {
			akevitt.ErrorBox("character name must not be empty", session.Application, form)
			return
		}

		bundle := session.Data[IronExaltBundle].(Bundle)
		bundle.Character.Name = characterName
		bundle.Character.UpdateRoom(engine.GetSpawnRoom())

		database, err := akevitt.FetchPlugin[akevitt.DatabasePlugin[*akevitt.Account]](engine)

		session.Account.PersistentData[IronExaltBundle] = bundle.Character
		channels := session.Data[plugins.MessagePluginData].([]string)
		channels = append(channels, characterName)
		session.Data[plugins.MessagePluginData] = channels

		if err != nil {
			panic(err)
		}
		if err := (*database).Save(session.Account); err != nil {
			panic(err)
		}

		session.Application.SetRoot(gameScreen(engine, session), true)
	})

	return form
}

func gameScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	playerMessage := ""

	defaults, _ := akevitt.FetchPlugin[*plugins.DefaultPlugins](engine)

	(*defaults).Messages.GetChatLog(session)

	inputField := tview.NewInputField().
		SetChangedFunc(func(text string) {
			playerMessage = text
		})

	game := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		AddItem(inputField, 2, 0, 1, 3, 0, 0, true).
		AddItem((*defaults).Messages.GetChatLog(session), 1, 1, 1, 2, 0, 0, false).
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

		akevitt.AppendText("\t>"+playerMessage, (*defaults).Messages.GetChatLog(session))
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
