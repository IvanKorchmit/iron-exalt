package main

import (
	"github.com/IvanKorchmit/akevitt"
)

const IronExaltBundle = "Iron Exalt Bundle"

type Bundle struct {
	Character Character
	Familiars []string
}

type Character struct {
	Name    string
	room    *akevitt.Room
	RoomKey uint64
}

func (character *Character) UpdateRoom(room *akevitt.Room) {
	character.room = room
	character.RoomKey = room.GetKey()
}

func (character *Character) GetName() string {
	return character.Name
}

func InitBundle(session *akevitt.ActiveSession) {
	session.Data[IronExaltBundle] = Bundle{
		Character: Character{},
		Familiars: make([]string, 0),
	}
}

func Fetch[T any](account *akevitt.Account, key string) *T {
	result, ok := account.PersistentData[key].(T)

	if !ok {
		return nil
	}

	return &result
}
