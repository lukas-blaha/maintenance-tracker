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
	listMachinesModel
	addMachineForm
	// updateMachineModel
)

var models []tea.Model

type App struct {
	Machines []Machine
}

func (a App) ListMachines() string {
	var text string
	for _, m := range a.Machines {
		if len(text) == 0 {
			text = m.title
		} else {
			text += " " + m.title
		}
	}

	return text
}

func main() {
	app := &App{}
	menu := NewMenu(app)
	menu.InitMenu(20, 10)

	models = []tea.Model{menu, NewMachineForm(app), NewMachinesList(app)}
	m := models[menuModel]
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
