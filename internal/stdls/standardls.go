package stdls

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SiirRandall/lsd-go/internal/config"
	"github.com/SiirRandall/lsd-go/internal/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	terminal "golang.org/x/term"
)

type fileEntry struct {
	original string
	styled   string
}

type model struct {
	grid [][]string
	quit bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.quit = true
	}
	return m, nil
}

func (m model) View() string {
	var output []string
	columns := len(m.grid)
	for row := 0; row < len(m.grid[0]); row++ {
		var rowStrs []string
		for col := 0; col < columns; col++ {
			if row < len(m.grid[col]) {
				rowStrs = append(rowStrs, m.grid[col][row])
			}
		}
		output = append(output, strings.Join(rowStrs, "  "))
	}
	return strings.Join(output, "\n")
}

func StdLS(config config.Config) {
	dir := config.Dir
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}
	// Filter out dot files if the showdot flag is not set.
	filteredFiles := []fs.DirEntry{}
	for _, file := range files {
		if config.ShowDotFiles || !strings.HasPrefix(file.Name(), ".") {
			filteredFiles = append(filteredFiles, file)
		}
	}
	files = filteredFiles

	sort.Slice(files, sorter(files, config))

	width, _, _ := terminal.GetSize(int(os.Stdout.Fd()))

	maxFilenameLength := 0
	for _, file := range files {
		icon, _ := getIconAndColorForFileOrDir(dir, file.Name())
		visualLength := visualWidth(file.Name()) + visualWidth(icon)
		if visualLength > maxFilenameLength {
			maxFilenameLength = visualLength
		}
	}

	columnSpacing := 2 // space between columns
	initialColumnWidth := maxFilenameLength + columnSpacing
	numColumns := width / initialColumnWidth

	var grid [][]string

	for _, file := range files {
		if len(grid) == 0 || len(grid[len(grid)-1]) >= (len(files)+numColumns-1)/numColumns {
			grid = append(grid, []string{})
		}
		grid[len(grid)-1] = append(grid[len(grid)-1], file.Name())
	}

	for col, column := range grid {
		maxColWidth := 0
		for _, filename := range column {
			icon, _ := getIconAndColorForFileOrDir(dir, filename)
			visualLength := visualWidth(filename) + visualWidth(icon)
			if visualLength > maxColWidth {
				maxColWidth = visualLength
			}
		}

		for idx, filename := range column {
			icon, style := getIconAndColorForFileOrDir(dir, filename)
			paddedName := fmt.Sprintf("%-*s", maxColWidth, icon+filename)
			grid[col][idx] = style.Render(paddedName)
		}
	}

	m := model{grid: grid}
	fmt.Println(m.View())
}

func sorter(files []fs.DirEntry, config config.Config) func(i, j int) bool {
	return func(i, j int) bool {
		file1 := files[i]
		file2 := files[j]

		if config.DirsFirst {
			// If dirsFirst is enabled and one is a directory while the other isn't
			if file1.IsDir() && !file2.IsDir() {
				return true
			} else if !file1.IsDir() && file2.IsDir() {
				return false
			}
		}

		// If sortReverse is enabled
		if config.SortReverse {
			return strings.ToLower(file2.Name()) < strings.ToLower(file1.Name())
		}
		// Default to alphabetical order
		return strings.ToLower(file1.Name()) < strings.ToLower(file2.Name())
	}
}

func isDir(baseDir, filename string) bool {
	info, err := os.Stat(filepath.Join(baseDir, filename))
	if err != nil {
		return false
	}
	return info.IsDir()
}

func getIconAndColorForFileOrDir(baseDir, filename string) (string, lipgloss.Style) {
	// Check for directories
	if isDir(baseDir, filename) {
		if icon, ok := style.FileTypeIconMap[filename]; ok {
			return icon.Icon, lipgloss.NewStyle().Foreground(lipgloss.Color(icon.Color))
		}
		return " ", lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")) // Default directory icon
	}

	// Check for files by extension
	ext := strings.ToLower(filepath.Ext(filename))
	if icon, ok := style.ExtToFileTypeIconMap[ext]; ok {
		return icon.Icon, lipgloss.NewStyle().Foreground(lipgloss.Color(icon.Color))
	}
	return " ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")) // Default file icon
}

func visualWidth(s string) int {
	width := 0
	for _, r := range s {
		if runewidth.RuneWidth(r) == 2 {
			width += 2
		} else {
			width++
		}
	}
	return width
}
