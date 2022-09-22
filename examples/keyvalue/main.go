package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/arriqaaq/flashdb"
	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tea "github.com/charmbracelet/bubbletea"
)

type kvMetadata struct {
	commandinput.CmdMetadata
	eval func(tx *flashdb.Tx, m completerModel) ([]string, error)
	name string
}

var baseSuggestions = []input.Suggestion[kvMetadata]{
	{Text: "commit"},
	{Text: "delete", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}}}},
	{Text: "exists", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}}}},
	{Text: "expire", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<duration>")}}}},
	{Text: "get", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}}}},
	{Text: "hash", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<subcommand>")}}}},
	{Text: "rollback"},
	{Text: "set"},
	{Text: "set-key", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<value>")}, HasFlags: true}, eval: func(tx *flashdb.Tx, m completerModel) ([]string, error) {
		parsed := m.textInput.ParsedValue()
		key := parsed.Args.Value[0].Value
		value := parsed.Args.Value[1].Value
		for _, flag := range parsed.Flags.Value {
			if flag.Name == "-t" || flag.Name == "--ttl" {
				intVal, _ := strconv.Atoi(flag.Value.Value)
				err := tx.SetEx(key, value, int64(intVal))
				return []string{}, err
			}
		}
		err := tx.Set(key, value)
		return []string{}, err
	}}},
	{Text: "ttl", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}}}},
	{Text: "zset"},
}

var hashSuggestions = []input.Suggestion[kvMetadata]{
	{Text: "clear", Metadata: kvMetadata{name: "HClear", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "delete", Metadata: kvMetadata{name: "HDel", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<values...>")}, Level: 1}}},
	{Text: "exists", Metadata: kvMetadata{name: "HExists", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("[field]")}, Level: 1}, eval: func(tx *flashdb.Tx, m completerModel) ([]string, error) {
		parsed := m.textInput.ParsedValue()
		args := parsed.Args.Value
		if len(parsed.Args.Value) == 2 {
			exists := tx.HKeyExists(args[1].Value)
			return []string{strconv.FormatBool(exists)}, nil
		} else if len(parsed.Args.Value) == 3 {
			exists := tx.HExists(args[1].Value, args[2].Value)
			return []string{strconv.FormatBool(exists)}, nil
		} else {
			return nil, fmt.Errorf("exists requires 2 or 3 args")
		}
	}}},
	{Text: "expire", Metadata: kvMetadata{name: "HExpire", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<duration>")}, Level: 1}}},
	{Text: "get", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("[field]")}, Level: 1}, eval: func(tx *flashdb.Tx, m completerModel) ([]string, error) {
		parsed := m.textInput.ParsedValue()
		key := parsed.Args.Value[1].Value

		for _, flag := range parsed.Flags.Value {
			if flag.Name == "-a" || flag.Name == "--all" {
				all := tx.HGetAll(key)
				return all, nil
			}
		}
		if len(parsed.Args.Value) < 3 {
			return []string{}, fmt.Errorf("get requires two args or the --all flag")
		}

		val := tx.HGet(key, parsed.Args.Value[2].Value)
		return []string{val}, nil
	}}},
	{Text: "keys", Metadata: kvMetadata{name: "HKeys", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "len", Metadata: kvMetadata{name: "HLen", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "set", Metadata: kvMetadata{name: "HSet", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<field>"), commandinput.NewPositionalArg("<value>")}, Level: 1}}},
	{Text: "ttl", Metadata: kvMetadata{name: "HTTL", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "values", Metadata: kvMetadata{name: "HVals", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
}

var setSuggestions = []input.Suggestion[kvMetadata]{
	{Text: "add", Metadata: kvMetadata{name: "SAdd", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<members...>")}, Level: 1}}},
	{Text: "card", Metadata: kvMetadata{name: "SCard", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "clear", Metadata: kvMetadata{name: "SClear", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "diff", Metadata: kvMetadata{name: "SDiff", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<keys...>")}, Level: 1}}},
	{Text: "exists", Metadata: kvMetadata{name: "SKeyExists", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "expire", Metadata: kvMetadata{name: "SExpire", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<duration>")}, Level: 1}}},
	{Text: "is-member", Metadata: kvMetadata{name: "SIsMember", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<member>")}, Level: 1}}},
	{Text: "members", Metadata: kvMetadata{name: "SMembers", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "move", Metadata: kvMetadata{name: "SMove", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<source>"), commandinput.NewPositionalArg("<destination>"), commandinput.NewPositionalArg("<member>")}, Level: 1}}},
	{Text: "random", Metadata: kvMetadata{name: "SRandMember", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<count>")}, Level: 1}}},
	{Text: "remove", Metadata: kvMetadata{name: "SRem", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<members...>")}, Level: 1}}},
	{Text: "ttl", Metadata: kvMetadata{name: "STTL", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "union", Metadata: kvMetadata{name: "SUnion", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<keys...>")}, Level: 1}}},
}

var zsetSuggestions = []input.Suggestion[kvMetadata]{
	{Text: "add", Metadata: kvMetadata{name: "ZAdd", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<score>"), commandinput.NewPositionalArg("<member>")}, Level: 1}}},
	{Text: "card", Metadata: kvMetadata{name: "ZCard", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "clear", Metadata: kvMetadata{name: "ZClear", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "exists", Metadata: kvMetadata{name: "ZKeyExists", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
	{Text: "expire", Metadata: kvMetadata{name: "ZExpire", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<duration>")}, Level: 1}}},
	{Text: "get-by-rank", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<rank>")}, Level: 1, HasFlags: true}, eval: func(tx *flashdb.Tx, m completerModel) ([]string, error) {
		parsed := m.textInput.ParsedValue()
		args := parsed.Args.Value
		for _, flag := range parsed.Flags.Value {
			if flag.Name == "-r" || flag.Name == "--reverse" {
				intVal, err := strconv.ParseInt(args[2].Value, 10, 32)
				if err != nil {
					return nil, err
				}
				all := tx.ZRevGetByRank(args[1].Value, int(intVal))
				retVals := []string{}
				for _, ret := range all {
					retVals = append(retVals, fmt.Sprintf("%v", ret))
				}
				return retVals, nil

			}
		}
		intVal, err := strconv.ParseInt(args[2].Value, 10, 32)
		if err != nil {
			return nil, err
		}
		all := tx.ZGetByRank(args[1].Value, int(intVal))
		retVals := []string{}
		for _, ret := range all {
			retVals = append(retVals, fmt.Sprintf("%v", ret))
		}
		return retVals, nil
	}}}, //-reverse
	{Text: "range", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<start>"), commandinput.NewPositionalArg("<stop>")}, Level: 1}, eval: func(tx *flashdb.Tx, m completerModel) ([]string, error) {
		return nil, nil
	}}}, // -scores -reverse
	{Text: "rank", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<member>")}, Level: 1}, eval: func(tx *flashdb.Tx, m completerModel) ([]string, error) {
		return nil, nil
	}}}, // -reverse
	{Text: "remove", Metadata: kvMetadata{name: "ZRem", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<member>")}, Level: 1}}},
	{Text: "score", Metadata: kvMetadata{name: "ZScore", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<member>")}, Level: 1}}},
	{Text: "score-range", Metadata: kvMetadata{CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>"), commandinput.NewPositionalArg("<min>"), commandinput.NewPositionalArg("<max>")}, Level: 1}, eval: func(tx *flashdb.Tx, m completerModel) ([]string, error) {
		return nil, nil
	}}}, //-rev
	{Text: "ttl", Metadata: kvMetadata{name: "ZTTL", CmdMetadata: commandinput.CmdMetadata{PositionalArgs: []commandinput.PositionalArg{commandinput.NewPositionalArg("<key>")}, Level: 1}}},
}

func (kv kvMetadata) Create(args []commandinput.PositionalArg, placeholder commandinput.Placeholder) commandinput.CmdMetadataAccessor {
	meta, _ := new(commandinput.CmdMetadata).Create(args, placeholder).(commandinput.CmdMetadata)
	return kvMetadata{
		CmdMetadata: meta,
	}
}

type model struct {
	prompt prompt.Model[kvMetadata]
}

type completerModel struct {
	db        *flashdb.FlashDB
	textInput *commandinput.Model[kvMetadata]
}

func (m model) Init() tea.Cmd {
	return m.prompt.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p, cmd := m.prompt.Update(msg)
	m.prompt = p
	return m, cmd
}

func (m model) View() string {
	return m.prompt.View()
}

func (m completerModel) completer(document prompt.Document, promptModel prompt.Model[kvMetadata]) ([]input.Suggestion[kvMetadata], error) {
	suggestions := []input.Suggestion[kvMetadata]{}
	args := m.textInput.CompletedArgsBeforeCursor()
	numArgs := len(args)
	filterSuggestions := true
	m.db.View(func(tx *flashdb.Tx) error {
		if !m.textInput.CommandCompleted() {
			suggestions = baseSuggestions
		} else {
			command := m.textInput.CommandBeforeCursor()
			switch command {
			case "set-key":
				if numArgs >= 2 {
					filterSuggestions = false
					suggestions = m.textInput.FlagSuggestions(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundUp), []commandinput.Flag{{
						Short:       "t",
						Long:        "ttl",
						Description: "Key TTL",
						Placeholder: "<ttl>",
						RequiresArg: true,
					}}, nil)
				}

			case "hash":
				if numArgs == 0 {
					suggestions = hashSuggestions
				} else {
					switch args[0] {
					case "get":
						if numArgs == 2 {
							filterSuggestions = false
							suggestions = append(suggestions, m.textInput.FlagSuggestions(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundUp), []commandinput.Flag{{
								Short:       "a",
								Long:        "all",
								Description: "Get all",
							}}, func(f commandinput.Flag) kvMetadata {
								return kvMetadata{name: "HGetAll", CmdMetadata: commandinput.CmdMetadata{FlagPlaceholder: commandinput.Placeholder{Text: f.Placeholder}, Level: 2, PreservePlaceholder: true}}
							})...)

							return nil
						}
					}
				}

			case "set":
				if numArgs == 0 {
					suggestions = setSuggestions
				}
			case "zset":
				if numArgs > 0 {
					switch args[0] {
					case "get-by-rank":
						if numArgs >= 3 {
							filterSuggestions = false
							suggestions = m.textInput.FlagSuggestions(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundUp), []commandinput.Flag{{
								Short:       "r",
								Long:        "reverse",
								Description: "Invert results",
							}}, nil)
						}
						return nil
					}
				} else {
					suggestions = zsetSuggestions
				}

			}
		}

		return nil
	})
	if filterSuggestions {
		return completers.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(commandinput.RoundDown), suggestions), nil
	}
	return suggestions, nil
}

func (m completerModel) executor(input string, selectedSuggestion *input.Suggestion[kvMetadata]) (tea.Model, error) {

	outStr := ""

	err := m.db.Update(func(tx *flashdb.Tx) error {
		parsed := m.textInput.ParsedValue()
		selectedCommand := m.textInput.SelectedCommand()
		if selectedCommand == nil {
			return nil
		}
		subcommand := ""

		if len(parsed.Args.Value) > 0 {
			subcommand = parsed.Args.Value[0].Value
		}
		switch selectedCommand.Text {
		case "hash":
			for _, suggestion := range hashSuggestions {
				if suggestion.Text == subcommand {
					text, err := m.execMethod(tx, &suggestion, m.textInput.AllValues()[2:])
					outStr = text
					return err
				}
			}
		case "set":
			for _, suggestion := range setSuggestions {
				if suggestion.Text == subcommand {
					text, err := m.execMethod(tx, &suggestion, m.textInput.AllValues()[2:])
					outStr = text
					return err
				}
			}
		case "zset":
			for _, suggestion := range zsetSuggestions {
				if suggestion.Text == subcommand {
					text, err := m.execMethod(tx, &suggestion, m.textInput.AllValues()[2:])
					outStr = text
					return err
				}
			}
		default:
			text, err := m.execMethod(tx, m.textInput.SelectedCommand(), m.textInput.AllValues()[1:])
			outStr = text
			return err
		}

		return nil
	})

	return executors.NewStringModel(outStr), err
}

func (m completerModel) execMethod(tx *flashdb.Tx, suggestion *input.Suggestion[kvMetadata], params []string) (string, error) {
	methodName := ""
	if suggestion != nil {

		if suggestion.Metadata.eval != nil {
			retVals, err := suggestion.Metadata.eval(tx, m)
			outStr := strings.Join(retVals, " ")
			return outStr, err
		}
		if len(suggestion.Metadata.name) > 0 {
			methodName = suggestion.Metadata.name
		} else {
			methodName = strings.ToUpper(string(suggestion.Text[0])) + suggestion.Text[1:]
		}
	} else {
		methodName = m.textInput.ParsedValue().Command.Value
	}

	method, found := reflect.TypeOf(tx).MethodByName(methodName)
	if !found {
		return "", fmt.Errorf("command not found")
	}

	expectedParams := method.Type.NumIn()
	isVariadic := method.Type.In(expectedParams-1).Kind() == reflect.Slice
	if (isVariadic && len(params) < expectedParams-1) || (!isVariadic && len(params) != expectedParams-1) {
		// Subtract one for the tx object
		return "", fmt.Errorf("expected %d params but got %d", expectedParams-1, len(params))
	}
	paramVals, err := getReflectParams(params, tx, method.Type)
	if err != nil {
		return "", err
	}

	out := method.Func.Call(paramVals)
	retVals := []string{}
	for _, outVal := range out {
		if outVal.CanInterface() {
			iface := outVal.Interface()
			if iface == nil {
				continue
			}
			switch ifaceVal := iface.(type) {
			case error:
				return "", fmt.Errorf(ifaceVal.Error())
			case []string:
				retVals = append(retVals, strings.Join(ifaceVal, ","))
			case string:
				retVals = append(retVals, ifaceVal)
			case bool:
				retVals = append(retVals, strconv.FormatBool(ifaceVal))
			case int64:
				retVals = append(retVals, strconv.FormatInt(ifaceVal, 10))
			case int:
				retVals = append(retVals, strconv.FormatInt(int64(ifaceVal), 10))
			case float64:
				retVals = append(retVals, strconv.FormatFloat(float64(ifaceVal), 'f', 3, 64))
			case float32:
				retVals = append(retVals, strconv.FormatFloat(float64(ifaceVal), 'f', 3, 32))

			}
		} else {
			retVals = append(retVals, outVal.String())
		}

	}
	return strings.Join(retVals, " "), nil
}

func getReflectParams(params []string, tx *flashdb.Tx, methodType reflect.Type) ([]reflect.Value, error) {
	paramVals := []reflect.Value{reflect.ValueOf(tx)}
	for i, p := range params {
		var reflectVal any
		var err error
		methodParam := methodType.In(i + 1)
		switch methodParam.Kind() {
		case reflect.Int:
			var intVal int64
			intVal, err = strconv.ParseInt(p, 10, 32)
			reflectVal = int(intVal)
		case reflect.Int64:
			reflectVal, err = strconv.ParseInt(p, 10, 64)
		case reflect.Float32:
			reflectVal, err = strconv.ParseFloat(p, 32)
		case reflect.Float64:
			reflectVal, err = strconv.ParseFloat(p, 64)
		case reflect.String:
			reflectVal = p
		case reflect.Slice:
			for j := i; j < len(params); j++ {
				paramVals = append(paramVals, reflect.ValueOf(params[j]))
			}
			return paramVals, nil
		}
		if err != nil {
			return nil, err
		}
		paramVals = append(paramVals, reflect.ValueOf(reflectVal))

	}
	return paramVals, nil
}

func main() {
	config := &flashdb.Config{}
	db, _ := flashdb.New(config)

	var textInput input.Input[kvMetadata] = commandinput.New[kvMetadata]()
	completerModel := completerModel{db: db, textInput: textInput.(*commandinput.Model[kvMetadata])}

	m := model{prompt: prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
	)}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
