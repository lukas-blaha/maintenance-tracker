package main

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	labelWidth = 18
	inputWidth = 40
)

var normalMachineStyle = lipgloss.NewStyle().
	Bold(true).
	Padding(0, 5).
	Border(lipgloss.HiddenBorder()).
	Foreground(lipgloss.Color("#808080"))

var selectedMachineStyle = lipgloss.NewStyle().
	Padding(1, 5).
	Bold(true).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62"))

var (
	machineTitleStyle = lipgloss.NewStyle().
				Width(20).
				PaddingRight(2)

	machineDescriptionStyle = lipgloss.NewStyle().
				Width(30).
				PaddingRight(2)

	machineHoursStyle = lipgloss.NewStyle().
				Width(20)
)

var (
	labelStyle = lipgloss.NewStyle().
			Width(labelWidth).
			Align(lipgloss.Left).
			PaddingRight(2)

	inputStyle = lipgloss.NewStyle().
			Width(inputWidth)
)

var machineHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("62")).
	MarginLeft(5)

func machineHeader() string {
	return machineHeaderStyle.Render(lipgloss.JoinHorizontal(
		lipgloss.Left,
		machineTitleStyle.Render("Title"),
		machineDescriptionStyle.Render(" Description"),
		machineHoursStyle.Render("Hours"),
	))
}

type Machine struct {
	title       string
	description string
	hours       string
}

func NewMachine(title, description, hours string) Machine {
	return Machine{title: title, description: description, hours: hours}
}

func (m Machine) FilterValue() string {
	return m.title
}

func (m Machine) Title() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		machineTitleStyle.Render(m.title),
		machineDescriptionStyle.Render(m.description),
		machineHoursStyle.Render(m.hours),
	)
}

func (m Machine) Description() string {
	return m.description
}

type AddMachineForm struct {
	app         *App
	title       textinput.Model
	description textinput.Model
	hours       textinput.Model
	width       int
	height      int
}

func NewMachineForm(app *App) *AddMachineForm {
	form := &AddMachineForm{}
	form.app = app
	form.title = textinput.New()
	form.title.Focus()
	form.description = textinput.New()
	form.hours = textinput.New()
	return form
}

func (f AddMachineForm) CreateMachine() tea.Msg {
	machine := NewMachine(f.title.Value(), f.description.Value(), f.hours.Value())
	f.app.Machines = append(f.app.Machines, machine)
	return machine
}

func (f AddMachineForm) Init() tea.Cmd {
	return nil
}

func (f AddMachineForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height

		return f, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return f, tea.Quit
		case "enter":
			if f.title.Focused() {
				f.title.Blur()
				f.hours.Blur()
				f.description.Focus()
				return f, textinput.Blink
			} else if f.description.Focused() {
				f.title.Blur()
				f.description.Blur()
				f.hours.Focus()
				return f, textinput.Blink
			} else {
				models[addMachineForm] = f
				return models[menuModel], f.CreateMachine
			}
		}
	}
	if f.title.Focused() {
		f.title, cmd = f.title.Update(msg)
		return f, cmd
	} else if f.description.Focused() {
		f.description, cmd = f.description.Update(msg)
		return f, cmd
	} else {
		f.hours, cmd = f.hours.Update(msg)
		return f, cmd
	}
}

func (f AddMachineForm) View() tea.View {
	form := lipgloss.JoinVertical(
		lipgloss.Center,
		renderField("Title:", f.title.View()),
		renderField("Description:", f.description.View()),
		renderField("Hours:", f.hours.View()),
	)

	v := tea.NewView(form)
	v.AltScreen = true

	return v
}

func renderField(label string, input string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		labelStyle.Render(label),
		inputStyle.Render(input),
	)
}

type ListMachinesModel struct {
	app      *App
	machines list.Model
	loaded   bool
	width    int
	height   int
}

func NewMachinesList(app *App) *ListMachinesModel {
	return &ListMachinesModel{app: app}
}

func (m *ListMachinesModel) UpdateMachinesList(width, height int) {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = selectedMachineStyle
	delegate.Styles.NormalTitle = normalMachineStyle
	machineList := list.New([]list.Item{}, delegate, width, height)
	machineList.SetShowHelp(false)
	machineList.SetShowTitle(false)
	machineList.SetShowStatusBar(false)

	for _, machine := range m.app.Machines {
		machineList.InsertItem(len(machineList.Items()), machine)
	}

	m.machines = machineList
	m.loaded = true
}

func (m ListMachinesModel) Init() tea.Cmd {
	return nil
}

func (m ListMachinesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.machines.SetSize(100, 20)

		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc":
			models[listMachinesModel] = m
			return models[menuModel], nil
		case "enter":
			_, ok := m.machines.SelectedItem().(Machine)
			if !ok {
				return m, nil
			}

			return m, nil
		}
	}

	m.machines, cmd = m.machines.Update(msg)
	return m, cmd
}

func (m ListMachinesModel) View() tea.View {
	if !m.loaded {
		return tea.NewView("loading...")
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		machineHeader(),
		m.machines.View(),
	)

	v := tea.NewView(lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	))

	v.AltScreen = true
	return v
}
