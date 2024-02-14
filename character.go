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
	Name string
}

func InitBundle(session *akevitt.ActiveSession) {
	session.Data[IronExaltBundle] = Bundle{
		Character: Character{},
		Familiars: make([]string, 0),
	}
}
