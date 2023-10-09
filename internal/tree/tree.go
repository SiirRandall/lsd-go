package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SiirRandall/lsd-go/internal/config"
	"github.com/SiirRandall/lsd-go/internal/style"

	"github.com/charmbracelet/lipgloss"
)

func Tree(config config.Config) {
	startPath := config.Dir

	// Get the base directory for output
	baseDir := filepath.Base(startPath)
	if startPath == "." || startPath == "./" {
		baseDir = "."
	}
	iconStyle, found := style.FileTypeIconMap[baseDir]
	if !found {
		iconStyle = style.FileTypeIcon{Icon: " ", Color: "#00FFFF"}
	}
	icon := lipgloss.NewStyle().Foreground(lipgloss.Color(iconStyle.Color)).Render(iconStyle.Icon)
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	coloredName := cyan.Render(baseDir)

	result := icon + coloredName + "\n" + traverseDir(startPath, 0, config.MaxDepth, config)
	fmt.Print(result)
}

func traverseDir(path string, depth int, maxDepth int, config config.Config) string {
	if maxDepth != -1 && depth > maxDepth {
		return ""
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", path, "-", err)
		return ""
	}

	var filteredEntries []os.DirEntry
	for _, entry := range entries {
		if !config.ShowDotFiles && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		filteredEntries = append(filteredEntries, entry)
	}

	// Sort entries
	sort.Slice(filteredEntries, func(i, j int) bool {
		return strings.ToLower(filteredEntries[i].Name()) < strings.ToLower(filteredEntries[j].Name())
	})

	var out strings.Builder
	indent := strings.Repeat("│  ", depth)
	prefix := "├── "

	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
	white := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	for i, entry := range filteredEntries {
		if i == len(filteredEntries)-1 {
			prefix = "└── "
		}

		if entry.IsDir() {
			iconStyle, found := style.FileTypeIconMap[entry.Name()]
			if !found {
				iconStyle = style.FileTypeIcon{Icon: " ", Color: "#00FFFF"}
			}
			icon := lipgloss.NewStyle().Foreground(lipgloss.Color(iconStyle.Color)).Render(iconStyle.Icon)
			coloredName := cyan.Render(entry.Name())
			out.WriteString(indent + prefix + icon + coloredName + "\n")
			out.WriteString(traverseDir(filepath.Join(path, entry.Name()), depth+1, maxDepth, config))
		} else {
			ext := filepath.Ext(entry.Name())
			iconStyle, found := style.ExtToFileTypeIconMap[ext]
			if !found {
				iconStyle = style.FileTypeIcon{Icon: " ", Color: "#FFFFFF"}
			}
			icon := lipgloss.NewStyle().Foreground(lipgloss.Color(iconStyle.Color)).Render(iconStyle.Icon)
			coloredName := white.Render(entry.Name())
			out.WriteString(indent + prefix + icon + coloredName + "\n")
		}
	}

	return out.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
