package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/aschey/bubbleprompt/searchbar"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/peterhellberg/swapi"
)

func newModel() customSearchbar {
	textInput := simpleinput.New[any]()
	suggestions := []suggestion.Suggestion[any]{
		{Text: "people"},
		{Text: "planets"},
		{Text: "species"},
		{Text: "starships"},
		{Text: "vehicles"},
	}

	pmodel := promptModel{
		swapiClient: swapi.DefaultClient,
		suggestions: suggestions,
		textInput:   textInput,
		filterer:    completer.NewPrefixFilter[any](),
	}

	suggestionStyle := suggestion.DefaultFormatters().Minimal()
	suggestionStyle.Suggestions.Border(lipgloss.RoundedBorder(), false, true, true)

	searchModel := searchbar.New[any](pmodel, textInput, newListModel(),
		searchbar.WithSearchbarStyle[any](lipgloss.NewStyle().Border(lipgloss.RoundedBorder())),
		searchbar.WithPromptOptions(prompt.WithFormatters[any](suggestionStyle)))
	return customSearchbar{model: searchModel}
}

type customSearchbar struct {
	model searchbar.Model[any]
}

func (s customSearchbar) Init() tea.Cmd {
	return s.model.Init()
}

func (s customSearchbar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := s.model.Update(msg)
	s.model = model.(searchbar.Model[any])
	return s, cmd
}

func (s customSearchbar) View() string {
	promptModel := s.model.PromptModel()
	suggestionManager := promptModel.SuggestionManager()
	if len(suggestionManager.Suggestions()) == 0 {
		return s.model.View()
	}

	maxNameLen, _ := suggestionManager.MaxSuggestionWidth()
	offset := promptModel.SuggestionOffset()

	view := s.model.BaseView()

	borderWidth := 2

	topView := lipgloss.NewStyle().
		PaddingLeft(offset - borderWidth).
		Render("╮" + strings.Repeat(" ", maxNameLen+borderWidth) + "╭")

	overlayView := lipgloss.JoinVertical(lipgloss.Left, s.model.Input(), topView, s.model.Body())
	return renderer.PlaceOverlay(s.model.OverlayX(), s.model.OverlayY(), overlayView, view)
}

type promptModel struct {
	swapiClient *swapi.Client
	suggestions []suggestion.Suggestion[any]
	textInput   *simpleinput.Model[any]
	filterer    completer.Filterer[any]
}

func (m promptModel) Complete(promptModel prompt.Model[any]) ([]suggestion.Suggestion[any], error) {
	if len(m.textInput.Tokens()) > 1 {
		return nil, nil
	}

	return m.filterer.Filter(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}

func (m promptModel) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	return executor.NewCmdModel("", func() tea.Msg { return m.getItems(input) }), nil
}

func (m promptModel) getItems(input string) []list.Item {
	switch input {
	case "people":
		people, err := m.swapiClient.AllPeople(context.Background())
		if err != nil {
			return nil
		}
		items := []list.Item{}
		for _, p := range people {
			items = append(items, listItem{title: p.Name, description: "From: " + p.Homeworld})
		}
		return items
	case "planets":
		planets, err := m.swapiClient.AllPlanets(context.Background())
		if err != nil {
			return nil
		}
		items := make([]list.Item, len(planets))

		for _, p := range planets {
			items = append(items, listItem{title: p.Name, description: "Population: " + p.Population})
		}
		return items
	case "species":
		species, err := m.swapiClient.AllSpecies(context.Background())
		if err != nil {
			return nil
		}
		items := make([]list.Item, len(species))
		for _, p := range species {
			items = append(items, listItem{title: p.Name, description: "Speaks: " + p.Language})
		}
		return items
	case "starships":
		starships, err := m.swapiClient.AllStarships(context.Background())
		if err != nil {
			return nil
		}
		items := make([]list.Item, len(starships))
		for _, p := range starships {
			items = append(items, listItem{title: p.Name, description: "Model: " + p.Model})
		}
		return items
	case "vehicles":
		vehicles, err := m.swapiClient.AllVehicles(context.Background())
		if err != nil {
			return nil
		}
		items := make([]list.Item, len(vehicles))
		for _, p := range vehicles {
			items = append(items, listItem{title: p.Name, description: "Class: " + p.VehicleClass})
		}
		return items
	}
	return []list.Item{}
}

func (m promptModel) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	return m, nil
}

type listItem struct {
	title       string
	description string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.description }
func (i listItem) FilterValue() string { return i.title }

type listModel struct {
	list.Model
}

func newListModel() listModel {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	return listModel{list}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []list.Item:
		cmd := m.Model.SetItems(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.Model.SetSize(msg.Width, msg.Height)
		return m, nil
	}
	list, cmd := m.Model.Update(msg)
	m.Model = list
	return m, cmd
}

func main() {
	if _, err := tea.NewProgram(newModel(), tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
