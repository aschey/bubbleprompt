package prompt

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type FilePathCompleter struct {
	Filter        func(de fs.DirEntry) bool
	IgnoreCase    bool
	fileListCache map[string][]Suggestion
}

func cleanFilePath(path string) (dir string, base string, err error) {
	if path == "" {
		return ".", "", nil
	}

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, "", nil
	}

	var endsWithSeparator bool
	if len(path) >= 1 && equalsSeparator(path[len(path)-1]) {
		endsWithSeparator = true
	}

	if len(path) >= 2 && path[0:1] == "~" && equalsSeparator(path[1]) {
		me, err := user.Current()
		if err != nil {
			return "", "", err
		}
		path = filepath.Join(me.HomeDir, path[1:])
	}
	path = filepath.Clean(os.ExpandEnv(path))
	dir = filepath.Dir(path)
	base = filepath.Base(path)

	if endsWithSeparator {
		dir = path + string(os.PathSeparator) // Append slash(in POSIX) if path ends with slash.
		base = ""                             // Set empty string if path ends with separator.
	}
	return dir, base, nil
}

func equalsSeparator(check byte) bool {
	return strings.ContainsAny(string(check), "/\\")
}

func (c *FilePathCompleter) adjustCompletions(completions []Suggestion, sub string) []Suggestion {
	//tokens := strings.Split(sub, " ")
	filteredCompletions := FilterCompletionTextHasPrefix(sub, completions)
	// if len(tokens) > 1 {
	// 	allExceptLast := strings.Join(tokens[0:len(tokens)-1], " ")
	// 	newCompletions := []Suggestion{}
	// 	for _, completion := range filteredCompletions {
	// 		completion.Text = completion.Text[len(allExceptLast)+1:]
	// 		newCompletions = append(newCompletions, completion)
	// 	}

	// 	return newCompletions
	// }
	return filteredCompletions
}

func (c *FilePathCompleter) Complete(path string) []Suggestion {
	if c.fileListCache == nil {
		c.fileListCache = make(map[string][]Suggestion, 4)
	}

	dir, base, err := cleanFilePath(path)
	if err != nil {
		fmt.Println("completer: cannot get current user:" + err.Error())
		return nil
	}

	if cached, ok := c.fileListCache[dir]; ok {
		return c.adjustCompletions(cached, base)
	}
	isAbs := filepath.IsAbs(path) || strings.HasPrefix(path, "~")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting cwd", err)
		return nil
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error getting path", err)
		return nil
	}
	filePath, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println("Error getting path", err)
		return nil
	}

	suggests := make([]Suggestion, 0, len(files))
	for _, f := range files {
		if c.Filter != nil && !c.Filter(f) {
			continue
		}
		full := filepath.Join(filePath, f.Name())
		if !isAbs {
			full, err = filepath.Rel(cwd, full)
			if err != nil {
				fmt.Println("Error getting rel path", err)
			}
		}
		if strings.Contains(full, " ") {
			full = fmt.Sprintf("\"%s\"", full)
		}
		suggests = append(suggests, Suggestion{
			Text:           full,
			CompletionText: f.Name(),
		})
	}
	c.fileListCache[dir] = suggests
	return c.adjustCompletions(suggests, base)
}
