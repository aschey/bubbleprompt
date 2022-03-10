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

func CleanFilePath(path string) (dir, base string, err error) {
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
	tokens := strings.Split(sub, " ")
	filteredCompletions := FilterHasPrefix(sub, completions)
	if len(tokens) > 1 {
		allExceptLast := strings.Join(tokens[0:len(tokens)-1], " ")
		newCompletions := []Suggestion{}
		for _, completion := range filteredCompletions {
			completion.Name = completion.Name[len(allExceptLast)+1:]
			newCompletions = append(newCompletions, completion)
		}

		return newCompletions
	}
	return filteredCompletions
}

func (c *FilePathCompleter) Complete(path string) []Suggestion {
	if c.fileListCache == nil {
		c.fileListCache = make(map[string][]Suggestion, 4)
	}

	dir, base, err := CleanFilePath(path)
	if err != nil {
		fmt.Println("completer: cannot get current user:" + err.Error())
		return nil
	}

	if cached, ok := c.fileListCache[dir]; ok {
		return c.adjustCompletions(cached, base)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
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
		suggests = append(suggests, Suggestion{Name: f.Name(), Metadata: filepath.Join(filePath, f.Name())})
	}
	c.fileListCache[dir] = suggests
	return c.adjustCompletions(suggests, base)
}
