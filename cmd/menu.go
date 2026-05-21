package main

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const menuItemWidth = 40

var normalMenuStyle = lipgloss.NewStyle().
	Width(menuItemWidth).
	Bold(true).
	Padding(0, 5).
	Border(lipgloss.HiddenBorder()).
	Foreground(lipgloss.Color("#808080"))

var selectedMenuStyle = lipgloss.NewStyle().
	Padding(1, 5).
	Bold(true).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62"))

const (
	listEntries item = iota
	listMachines
	addEntry
	updateEntry
	addMachine
	updateMachine
)

type MenuItem struct {
	item        item
	title       string
	description string
}

func (i MenuItem) FilterValue() string {
	return i.title
}

func (i MenuItem) Title() string {
	return i.title
}

func (i MenuItem) Description() string {
	return i.description
}

type MenuModel struct {
	app      *App
	focused  item
	items    list.Model
	err      error
	loaded   bool
	quitting bool
	width    int
	height   int
}

func NewMenu(app *App) *MenuModel {
	return &MenuModel{app: app}
}

func (m *MenuModel) InitMenu(width, height int) {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = selectedMenuStyle
	delegate.Styles.NormalTitle = normalMenuStyle
	menuList := list.New([]list.Item{}, delegate, width, height)
	menuList.SetShowHelp(false)
	menuList.SetShowTitle(false)
	menuList.SetShowStatusBar(false)

	m.items = menuList

	m.items.SetItems([]list.Item{
		MenuItem{listMachines, "Show saved entries", ""},
		MenuItem{addEntry, "Add new maintenance entry", ""},
		MenuItem{updateEntry, "Update existig entry", ""},
		MenuItem{addMachine, "Add new machine", ""},
		MenuItem{updateMachine, "Edit machine data", ""},
	})

	m.loaded = true
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// m.items.SetSize(msg.Width, msg.Height)
		m.items.SetSize(50, 20)

		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			selected, ok := m.items.SelectedItem().(MenuItem)
			if !ok {
				return m, nil
			}

			m.focused = selected.item

			switch selected.item {
			case listMachines:
				models[menuModel] = m
				machineList := NewMachinesList(m.app)
				machineList.UpdateMachinesList(100, 20)
				models[listMachinesModel] = machineList
				return models[listMachinesModel].Update(nil)
			case addEntry:
				// TODO: switch to add entry screen/model
			case updateEntry:
				// TODO: switch to update entry screen/model
			case addMachine:
				models[menuModel] = m
				models[addMachineForm] = NewMachineForm(m.app)
				return models[addMachineForm].Update(nil)
			case updateMachine:
				// TODO: switch to update machine screen/model
			}

			return m, nil
		}
	}

	m.items, cmd = m.items.Update(msg)
	return m, cmd
}

func (m MenuModel) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	if m.loaded {
		v := tea.NewView(lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.items.View(),
		))
		v.AltScreen = true
		return v
	} else {
		return tea.NewView("loading...")
	}
}
