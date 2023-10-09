package main

import (
	"flag"

	"github.com/SiirRandall/lsd-go/internal/config"
	"github.com/SiirRandall/lsd-go/internal/list"
	"github.com/SiirRandall/lsd-go/internal/starndardls"
	"github.com/SiirRandall/lsd-go/internal/tree"
)

var (
	showDotFiles     = flag.Bool("a", false, "Show dotfiles")
	noColor          = flag.Bool("no-color", false, "Disable colored output")
	showInodes       = flag.Bool("inodes", false, "Show inodes")
	headers          = flag.Bool("headers", false, "Show headers")
	listDetails      = flag.Bool("l", false, "List")
	sortAlphabetical = flag.Bool("alpha", false, "Sort files alphabetically")
	sortReverse      = flag.Bool("reverse", false, "Sort files in reverse order")
	dirsFirst        = flag.Bool("dirsfirst", true, "Sort directories first and then files alphabetically")
	maxDepth         = flag.Int("depth", 1, "Maximum depth for directory traversal. -1 means no limit.")
	treeview         = flag.Bool("tree", false, "Show tree view")
)

func main() {
	flag.Parse()

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(flag.NArg() - 1)
	}
	config := config.Config{
		SortAlphabetical: *sortAlphabetical,
		SortReverse:      *sortReverse,
		DirsFirst:        *dirsFirst,
		ShowDotFiles:     *showDotFiles,
		Dir:              dir,
		MaxDepth:         *maxDepth,
	}
	args := flag.Args() // Get the non-flag command-line arguments
	if *listDetails {
		list.ListFiles(args, *showDotFiles, *showInodes, *headers, *noColor)
	} else if *treeview {
		tree.Tree(config)
	} else {
		starndardls.StdLS(config)
	}
}
