package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
)

const (
	// Define the width of the columns
	fileTypeColWidth = 26
	linesColWidth    = 26
)

var fileTypeIcons = map[string]string{
	"go":            "🐹",
	"py":            "🐍",
	"js":            "✨",
	"ts":            "🟦",
	"html":          "🌐",
	"css":           "🎨",
	"java":          "☕",
	"c":             "🔵",
	"cpp":           "🔷",
	"rb":            "💎",
	"php":           "🐘",
	"rs":            "🦀",
	"sh":            "🐚",
	"yaml":          "📄",
	"json":          "📦",
	"md":            "📝",
	"lua":           "🌙",
	"tf":            "🌍",
	"tpl":           "🔧",
	"helm":          "⛵",
	"dockerfile":    "🐳",
	"dockercompose": "📦",
	"other":         "❓",
}

var codeFileExtensions = map[string]bool{
	"go":                 true,
	"py":                 true,
	"js":                 true,
	"ts":                 true,
	"html":               true,
	"css":                true,
	"java":               true,
	"c":                  true,
	"cpp":                true,
	"rb":                 true,
	"php":                true,
	"rs":                 true,
	"sh":                 true,
	"yaml":               true,
	"yml":                true,
	"json":               true,
	"md":                 true,
	"lua":                true,
	"tf":                 true,
	"tpl":                true,
	"helm":               true,
	"Dockerfile":         true,
	"docker-compose.yml": true,
}

func countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var count int
	buf := make([]byte, 32*1024)
	lineSep := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		count += strings.Count(string(buf[:c]), string(lineSep))

		switch err {
		case nil:
			continue
		case io.EOF:
			return count, nil
		default:
			return count, err
		}
	}
}

func getFileType(ext string, name string) string {
	ext = strings.TrimPrefix(ext, ".")
	if name == "Dockerfile" {
		return fmt.Sprintf("%s %s", fileTypeIcons["dockerfile"], "Dockerfile")
	}
	if name == "docker-compose.yml" {
		return fmt.Sprintf("%s %s", fileTypeIcons["dockercompose"], "docker-compose")
	}
	if icon, exists := fileTypeIcons[ext]; exists {
		return fmt.Sprintf("%s %s", icon, ext)
	}
	return fmt.Sprintf("%s other", fileTypeIcons["other"])
}

func processFile(path string, fileCounts map[string]int, totalLines *int, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	lines, err := countLines(path)
	if err != nil {
		return // Skip files that cannot be read
	}

	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	name := filepath.Base(path)
	var fileType string
	if codeFileExtensions[ext] || name == "Dockerfile" || name == "docker-compose.yml" {
		fileType = getFileType(ext, name)
		mutex.Lock()
		*totalLines += lines
		fileCounts[fileType] += lines
		mutex.Unlock()
	} else {
		fileType = getFileType("other", name)
		mutex.Lock()
		fileCounts[fileType] += lines
		mutex.Unlock()
	}
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}

	var totalLines int
	fileCounts := make(map[string]int)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip directories and inaccessible files
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// Skip hidden files and directories
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			wg.Add(1)
			go processFile(path, fileCounts, &totalLines, &mutex, &wg)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", dir, err)
		os.Exit(1)
	}

	wg.Wait()

	// Print the results in a pretty ASCII table
	color.Set(color.FgYellow)
	color.Unset()

	// Top border of the table
	fmt.Printf("+-%s-+-%s-+\n", strings.Repeat("-", fileTypeColWidth), strings.Repeat("-", linesColWidth))

	// Table headers
	fmt.Printf("| %-*s | %-*s |\n", fileTypeColWidth, "File Type", linesColWidth, "Lines of Code")

	// Header separator
	fmt.Printf("+-%s-+-%s-+\n", strings.Repeat("-", fileTypeColWidth), strings.Repeat("-", linesColWidth))

	for fileType, lines := range fileCounts {
		if fileType == "❓ other" {
			continue
		}
		color.Set(color.FgGreen)
		fmt.Printf("|%-*s | ", fileTypeColWidth, fileType)
		color.Set(color.FgCyan)
		fmt.Printf("%-*d |\n", linesColWidth, lines)
		color.Unset()
	}

	// Bottom border of the table
	fmt.Printf("+-%s-+-%s-+\n", strings.Repeat("-", fileTypeColWidth), strings.Repeat("-", linesColWidth))
	color.Set(color.FgMagenta)
	fmt.Printf("| %-*s | %-*d |\n", fileTypeColWidth, "Total", linesColWidth, totalLines)
	color.Unset()
	fmt.Printf("+-%s-+-%s-+\n", strings.Repeat("-", fileTypeColWidth), strings.Repeat("-", linesColWidth))
}
