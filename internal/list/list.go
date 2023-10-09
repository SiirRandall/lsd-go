package list

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/SiirRandall/lsd-go/internal/osfiles"
	"github.com/SiirRandall/lsd-go/internal/style"

	"github.com/charmbracelet/lipgloss"
)

const (
	linkcolor  = "#BE67F5"
	executable = "#ff0303" //"#b00000"
	whiteColor = "#FFFFFF"
	center     = lipgloss.Center
)

type maxLen struct {
	userLen    int
	groupLen   int
	sizeNumLen int
	inodeLen   int
}

var noColor *bool

func ListFiles(args []string, showDotFiles bool, showInodes bool, headers bool, noColor bool) {
	files, dir := osfiles.GetFiles(args, showDotFiles)
	max := maxLen{}
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving file info: %v\n", err)
			continue // skip to the next iteration
		}
		user, group := getUserAndGroup(fileInfo, dir)
		user = strings.ReplaceAll(user, " ", "")
		group = strings.ReplaceAll(group, " ", "")
		if len(user) > max.userLen {
			max.userLen = len(user)
		}
		if len(group) > max.groupLen {
			max.groupLen = len(group)
		}
		sizeNum, _ := formatSize(fileInfo.Size())
		if len(sizeNum) > max.sizeNumLen {
			max.sizeNumLen = len(sizeNum)
		}
		sys, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			fmt.Fprintf(os.Stderr, "error retrieving file info: \n")
			continue // skip to the next iteration
		}
		inode := sys.Ino
		inodeLen := len(strconv.Itoa(int(inode)))
		if inodeLen > max.inodeLen {
			max.inodeLen = inodeLen
		}
	}
	if headers {
		if max.inodeLen < 6 {
			max.inodeLen = 6
		}

		inodeHeaderStyle := createHeaderStyle(whiteColor, max.inodeLen, center, "Inodes")
		permHeaderStyle := createHeaderStyle(whiteColor, 11, center, "Permissions")
		userHeaderStyle := createHeaderStyle(whiteColor, max.userLen, center, "User")
		groupHeaderStyle := createHeaderStyle(whiteColor, max.groupLen, center, "Group")
		sizeHeaderStyle := createHeaderStyle(whiteColor, max.sizeNumLen+2, center, "Size")
		timeHeaderStyle := createHeaderStyle(whiteColor, 24, center, "Last Modified")
		nameStyle := createHeaderStyle(whiteColor, 0, center, "Name") // Adjust width as necessary

		if headers && showInodes {
			fmt.Printf("%-s  %-s %-s %-s %s %-s %s %s\n",
				inodeHeaderStyle, permHeaderStyle, userHeaderStyle,
				groupHeaderStyle, sizeHeaderStyle, "", timeHeaderStyle, "Name")
		} else if headers && !showInodes {
			fmt.Printf("%-s %-s %-s %s %-s %s   %s\n",
				permHeaderStyle, userHeaderStyle, groupHeaderStyle,
				sizeHeaderStyle, "", timeHeaderStyle, nameStyle)
		}
	}

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving file info: %v\n", err)
			continue // skip to the next iteration
		}
		printFileDetails(dir, fileInfo, max, showInodes, headers, noColor)
	}
}

func printFileDetails(dir string, file os.FileInfo, max maxLen, showInodes bool, headers bool, noColor bool) {
	permStyledString := getPermissionStyle(file, noColor)
	user, group := getUserAndGroup(file, dir)
	sizeNum, sizeUnit := formatSize(file.Size())
	sizeStyle, color := getSizeStyleAndColor(file.Size())

	userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#fcfbd2")).Width(max.userLen)
	groupStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#d1d0ab")).Width(max.groupLen)
	fileNameStyle, nerdFontSymbol := getFileNameStyle(file)
	styledFileName := fileNameStyle.Render(nerdFontSymbol + file.Name())

	if file.Mode()&os.ModeSymlink != 0 {
		targetPath, err := os.Readlink(filepath.Join(dir, file.Name()))
		if err != nil {
			fmt.Println("Error reading symlink target:", err)
			return
		}
		styledFileName += " ⇒ " + lipgloss.NewStyle().Foreground(lipgloss.Color(linkcolor)).Render(targetPath)
	}

	var inodeStyle lipgloss.Style
	var styledString string
	if showInodes {
		sys, ok := file.Sys().(*syscall.Stat_t)
		if !ok {
			fmt.Fprintf(os.Stderr, "error retrieving file info: \n")
			return
		}
		inode := sys.Ino
		color = "#FFFFFF"
		inodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Width(max.inodeLen)
		styledString = inodeStyle.Render(strconv.Itoa(int(inode)))
		// lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Width(max.inodeLen).Render(strconv.Itoa(int(inode)))
	}
	if headers {
		permStyledString += " "
	}

	format := []interface{}{
		permStyledString,
		userStyle.Render(user),
		groupStyle.Render(group),
		sizeStyle.Align(lipgloss.Right).Width(max.sizeNumLen).Render(sizeNum),
		lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Width(2).Render(sizeUnit),
		getTimeStyle(file.ModTime()).Width(24).Render(file.ModTime().Format("Mon Jan 02 15:04:05 2006")),
		styledFileName,
	}

	if showInodes {
		if headers {
			format = append([]interface{}{styledString, ""}, format...)
		} else {
			format = append([]interface{}{styledString}, format...)
		}
	}

	fmt.Printf(strings.Repeat("%-s ", len(format)-1)+"%s\n", format...)
}

func getGroupName(gid string) (string, error) {
	file, err := os.Open("/etc/group")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) >= 3 && parts[2] == gid {
			return parts[0], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("no matching group found")
}

func getUserAndGroup(file fs.FileInfo, dir string) (string, string) {
	// Attempt to get the actual user and group name
	sys, ok := file.Sys().(*syscall.Stat_t)
	if !ok {
		return "user", "group"
	}

	uid := strconv.Itoa(int(sys.Uid))
	gid := strconv.Itoa(int(sys.Gid))

	user, err := user.LookupId(uid)
	if err != nil {
		return uid, gid
	}

	group, err := getGroupName(gid)
	if err != nil {
		return user.Username, gid
	}
	// fmt.Println(len(user.Username))
	return user.Username, group
}

func formatSize(size int64) (string, string) {
	const KB = 1024
	const MB = 1024 * KB
	const GB = 1024 * MB

	switch {
	case size < KB:
		return strconv.FormatInt(size, 10), "B"
	case size < MB:
		if size < 10*KB {
			return fmt.Sprintf("%.1f", float64(size)/float64(KB)), "KB"
		}
		return fmt.Sprintf("%.0f", float64(size)/float64(KB)), "KB"
	case size < GB:
		if size < 10*MB {
			return fmt.Sprintf("%.1f", float64(size)/float64(MB)), "MB"
		}
		return fmt.Sprintf("%.0f", float64(size)/float64(MB)), "MB"
	default:
		if size < 10*GB {
			return fmt.Sprintf("%.1f", float64(size)/float64(GB)), "GB"
		}
		return fmt.Sprintf("%.0f", float64(size)/float64(GB)), "GB"
	}
}

func getPermissionStyle(fileInfo os.FileInfo, noColor bool) string {
	perm := fileInfo.Mode()
	var b strings.Builder

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		b.WriteString(colorize("l", linkcolor, noColor)) // Cyan for symlink
	} else if fileInfo.IsDir() {
		b.WriteString(colorize("d", "#00FFFF", noColor)) // Cyan for directory
	} else {
		b.WriteString("-")
	}

	permStr := formatPermissions(perm)
	for _, c := range permStr {
		switch c {
		case 'r':
			b.WriteString(colorize(string(c), "#00FF00", noColor)) // Green
		case 'w':
			b.WriteString(colorize(string(c), "#FFA500", noColor)) // Orange
		case 'x':
			b.WriteString(colorize(string(c), "#FF0000", noColor)) // Red
		case 's':
			b.WriteString(colorize(string(c), "#FFD700", noColor)) // Gold
		case 't':
			b.WriteString(colorize(string(c), "#FFC0CB", noColor)) // Pink
		case '-':
			b.WriteString(colorize(string(c), "", noColor)) // No color
		}
	}
	return b.String()
}

func colorize(input string, color string, noColor bool) string {
	if noColor {
		return input
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(input)
}

// Adjusting how getPermissionStyle is used in printFileDetails

func getSizeStyleAndColor(size int64) (lipgloss.Style, string) {
	if size < 1024*1024 { // KB
		color := "#FFFFFF" // replace with the exact white you desire
		return lipgloss.NewStyle().Foreground(lipgloss.Color(color)), color
	} else { // MB
		color := "#FFA500" // replace with the exact orange you desire
		return lipgloss.NewStyle().Foreground(lipgloss.Color(color)), color
	}
}

func getTimeStyle(t time.Time) lipgloss.Style {
	if t.Day() == time.Now().Day() {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#02e05f")) // replace with the exact bright green you desire
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#02bd91")) // replace with the exact dark green you desire
}

func formatPermissions(perm os.FileMode) string {
	var b strings.Builder

	// Check for special permissions
	setuid := perm&os.ModeSetuid != 0
	setgid := perm&os.ModeSetgid != 0
	sticky := perm&os.ModeSticky != 0

	for i := 0; i < 3; i++ {
		p := (perm >> (6 - 3*i)) & 7
		executeChar := 'x'
		if i == 0 && setuid { // User execute position
			executeChar = 's'
		} else if i == 1 && setgid { // Group execute position
			executeChar = 's'
		} else if i == 2 && sticky { // Other execute position
			executeChar = 't'
		}
		b.WriteString(fmt.Sprintf("%c%c%c",
			ifThenElse((p&4) != 0, 'r', '-'),
			ifThenElse((p&2) != 0, 'w', '-'),
			ifThenElse((p&1) != 0, executeChar, '-'),
		))
	}
	return b.String()
}

func ifThenElse(condition bool, a, b rune) rune {
	if condition {
		return a
	}
	return b
}

func getFileNameStyle(file os.FileInfo) (lipgloss.Style, string) {
	name := file.Name()
	var iconAndColor style.FileTypeIcon
	var ok bool

	if file.Mode()&os.ModeSymlink != 0 {
		iconAndColor = style.FileTypeIcon{Icon: " ", Color: linkcolor}
	} else if file.IsDir() {
		iconAndColor, ok = style.FileTypeIconMap[name]
		if !ok {
			iconAndColor = style.FileTypeIcon{Icon: "\uf115 ", Color: "#00FFFF"}
		}
	} else if isBinary(file) {
		iconAndColor = style.FileTypeIcon{Icon: "\uf489 ", Color: executable} // Green for binary files with bash icon
	} else {
		ext := strings.ToLower(filepath.Ext(name))
		iconAndColor, ok = style.ExtToFileTypeIconMap[ext]
		if !ok {
			iconAndColor = style.FileTypeIcon{Icon: "\uf15b ", Color: "#FFFFFF"}
		}
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(iconAndColor.Color))

	return style, iconAndColor.Icon
}

func isBinary(file os.FileInfo) bool {
	return file.Mode().Perm()&0111 != 0 // checks if any of the executable bits are set
}

func createHeaderStyle(color string, width int, align lipgloss.Position, text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Width(width).
		Underline(true).
		Align(align).
		Render(text)
}
