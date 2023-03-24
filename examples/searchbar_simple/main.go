package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/searchbar"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type url string

func getSections() []suggestion.Suggestion[url] {
	res, _ := http.Get("https://swapi.dev/api/")
	body, _ := io.ReadAll(res.Body)
	sections := make(map[string]string)
	_ = json.Unmarshal(body, &sections)
	suggestions := []suggestion.Suggestion[url]{}
	for section, sectionUrl := range sections {
		suggestions = append(suggestions, suggestion.Suggestion[url]{
			Text:     section,
			Metadata: url(sectionUrl),
		})
	}
	return suggestions
}

type response[T any] struct {
	Results []T
}

type person struct {
	Name      string
	BirthYear string `json:"birth_year"`
}

type planet struct {
	Name       string
	Population string
}

type species struct {
	Name     string
	Language string
}

type starship struct {
	Name  string
	Model string
}

type vehicle struct {
	Name       string
	Passengers string
}

type film struct {
	Title       string
	ReleaseDate string `json:"release_date"`
}

func getResource[T any](url url) []T {
	res, _ := http.Get(string(url))
	body, _ := io.ReadAll(res.Body)
	resBody := response[T]{}
	_ = json.Unmarshal(body, &resBody)
	return resBody.Results
}

func newModel() searchbar.Model[url] {
	textInput := simpleinput.New[url]()

	pmodel := promptModel{
		suggestions: []suggestion.Suggestion[url]{},
		textInput:   textInput,
		filterer:    completer.NewPrefixFilter[url](),
	}

	return searchbar.New[url](pmodel, textInput, newListModel())
}

type promptModel struct {
	suggestions []suggestion.Suggestion[url]
	textInput   *simpleinput.Model[url]
	filterer    completer.Filterer[url]
}

func (m promptModel) Complete(promptModel prompt.Model[url]) ([]suggestion.Suggestion[url], error) {
	if len(m.textInput.Tokens()) > 1 {
		return nil, nil
	}

	return m.filterer.Filter(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}

func (m promptModel) Execute(input string, promptModel *prompt.Model[url]) (tea.Model, error) {
	selected := promptModel.SuggestionManager().SelectedSuggestion()
	getItems := func() tea.Msg {
		return itemsMsg(m.getItems(selected.Text, selected.Metadata))
	}
	getTitle := func() tea.Msg {
		return titleMsg(strings.ToUpper(input[0:1]) + input[1:])
	}
	return executor.NewCmdModel("", tea.Batch(getItems, getTitle)), nil
}

type titleMsg string

type itemsMsg []list.Item

func (m promptModel) getItems(input string, url url) []list.Item {
	switch input {
	case "people":
		people := getResource[person](url)

		items := make([]list.Item, len(people))
		for i, p := range people {
			items[i] = listItem{title: p.Name, description: "Birth Year: " + p.BirthYear}
		}
		return items
	case "planets":
		planets := getResource[planet](url)

		items := make([]list.Item, len(planets))
		for i, p := range planets {
			items[i] = listItem{title: p.Name, description: "Population: " + p.Population}
		}
		return items
	case "species":
		species := getResource[species](url)

		items := make([]list.Item, len(species))
		for i, p := range species {
			items[i] = listItem{title: p.Name, description: "Speaks: " + p.Language}
		}
		return items
	case "starships":
		starships := getResource[starship](url)

		items := make([]list.Item, len(starships))
		for i, p := range starships {
			items[i] = listItem{title: p.Name, description: "Model: " + p.Model}
		}
		return items
	case "vehicles":
		vehicles := getResource[vehicle](url)

		items := make([]list.Item, len(vehicles))
		for i, p := range vehicles {
			items[i] = listItem{title: p.Name, description: "Passengers: " + p.Passengers}
		}
		return items
	case "films":
		vehicles := getResource[film](url)

		items := make([]list.Item, len(vehicles))
		for i, p := range vehicles {
			items[i] = listItem{title: p.Title, description: "Release Date: " + p.ReleaseDate}
		}
		return items
	}
	return []list.Item{}
}

func (m promptModel) Init() tea.Cmd {
	return suggestion.RefreshSuggestions(getSections)
}

func (m promptModel) Update(msg tea.Msg) (prompt.InputHandler[url], tea.Cmd) {
	if msg, ok := msg.(suggestion.RefreshSuggestionsMessage[url]); ok {
		m.suggestions = msg
	}
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
	list *list.Model
}

func newListModel() listModel {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	return listModel{list: &list}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case titleMsg:
		m.list.Title = string(msg)
	case itemsMsg:
		m.list.StopSpinner()
		cmds = append(cmds, m.list.SetItems(msg))
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-1)
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			m.list.SetShowTitle(true)
			cmds = append(cmds, m.list.StartSpinner())
		}
	}
	list, cmd := m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.list = &list
	return m, tea.Batch(cmds...)
}

func (m listModel) View() string {
	return "\n" + m.list.View()
}

func main() {
	if _, err := tea.NewProgram(newModel(), tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
