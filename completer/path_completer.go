package completer

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/aschey/bubbleprompt/input"
)

type PathCompleter[T any] struct {
	Filter        func(de fs.DirEntry) bool
	IgnoreCase    bool
	fileListCache map[string][]input.Suggestion[T]
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
	stripLast := false
	if strings.HasSuffix(path, ".") {
		// filepath functions treat "." as current directory
		// Need to add something to the end to treat the "." as a normal character
		path += "!"
		stripLast = true
	}

	path = filepath.Clean(os.ExpandEnv(path))
	dir = filepath.Dir(path)
	base = filepath.Base(path)
	if stripLast {
		base = base[:len(base)-1]
		path = path[:len(path)-1]
	}

	if endsWithSeparator {
		dir = path + string(os.PathSeparator) // Append slash(in POSIX) if path ends with slash.
		base = ""                             // Set empty string if path ends with separator.
	}
	return dir, base, nil
}

func equalsSeparator(check byte) bool {
	return strings.ContainsAny(string(check), "/\\")
}

func (c *PathCompleter[T]) adjustCompletions(completions []input.Suggestion[T], sub string) []input.Suggestion[T] {
	filteredCompletions := FilterHasPrefix(sub, completions)

	return filteredCompletions
}

func (c *PathCompleter[T]) Complete(path string) []input.Suggestion[T] {
	path = strings.ReplaceAll(path, "\"", "")
	if c.fileListCache == nil {
		c.fileListCache = make(map[string][]input.Suggestion[T], 4)
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
		return nil
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	filePath, err := filepath.Abs(dir)
	if err != nil {
		return nil
	}

	suggests := make([]input.Suggestion[T], 0, len(files))
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
		cursorOffset := 0
		if strings.Contains(full, " ") {
			full = fmt.Sprintf("\"%s\"", full)
			cursorOffset = 1
		}
		suggests = append(suggests, input.Suggestion[T]{
			Text:           full,
			CompletionText: f.Name(),
			CursorOffset:   cursorOffset,
		})
	}
	c.fileListCache[dir] = suggests
	return c.adjustCompletions(suggests, base)
}
