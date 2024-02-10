package main

import (
	"log"

	"github.com/IvanKorchmit/akevitt"
	"github.com/rivo/tview"
)

func main() {
	customRoom := akevitt.Room{Name: "My room"}
	engine := akevitt.NewEngine().
		UseRootUI(root).
		UseSpawnRoom(&customRoom).
		UseDBPath("database.db").
		UseBind(":2222").
		Finish()
	log.Fatal(engine.Run())
}

func root(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	return tview.NewModal().SetText("Hello from Akevitt!")
}
