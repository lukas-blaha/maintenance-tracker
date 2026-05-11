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

var (
	labelStyle = lipgloss.NewStyle().
			Width(labelWidth).
			Align(lipgloss.Left).
			PaddingRight(2)

	inputStyle = lipgloss.NewStyle().
			Width(inputWidth)
)

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
	return m.title
}

func (m Machine) Description() string {
	return m.description
}

type AddMachineForm struct {
	title       textinput.Model
	description textinput.Model
	hours       textinput.Model
	width       int
	height      int
}

func NewMachineForm() *AddMachineForm {
	form := &AddMachineForm{}
	form.title = textinput.New()
	form.title.Focus()
	form.description = textinput.New()
	form.hours = textinput.New()
	return form
}

func (f AddMachineForm) CreateMachine() tea.Msg {
	machine := NewMachine(f.title.Value(), f.description.Value(), f.hours.Value())
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

type MachinesModel struct {
	machines []list.Model
}
