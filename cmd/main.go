package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type item int

const (
	menuModel item = iota
	// listEntriesModel
	// addEntryModel
	// updateEntryModel
	// listMachinesModel
	addMachineForm
	// updateMachineModel
)

var models []tea.Model

func main() {
	menu := NewMenu()
	menu.InitMenu(20, 10)

	models = []tea.Model{menu, NewMachineForm()}
	m := models[menuModel]
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
