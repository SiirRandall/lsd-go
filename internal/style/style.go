package style

type FileTypeIcon struct {
	Icon  string
	Color string
}

// Maps to hold icon and color definitions for different file types
var FileTypeIconMap = map[string]FileTypeIcon{
	"Music":     {" ", "#00FFFF"},
	"Documents": {" ", "#00FFFF"},
	"Downloads": {"󱑢 ", "#00FFFF"},
	"Pictures":  {" ", "#00FFFF"},
	"Videos":    {"󰕧 ", "#00FFFF"},
	"Desktop":   {" ", "#00FFFF"},
	"Public":    {" ", "#00FFFF"},
}

var ExtToFileTypeIconMap = map[string]FileTypeIcon{
	".pub":  {"󱕵 ", "#FFFF00"},
	".py":   {"󰌠 ", "#add8e6"},
	".go":   {" ", "#039eff"},
	".zip":  {" ", "#FF0000"},
	".pdf":  {" ", "#FF0000"},
	".png":  {" ", "#00FF00"},
	".rs":   {"󰇷 ", "#663300"},
	".json": {"󰘦 ", "#FFFFFF"},
	".dmg":  {"󱧘 ", "#9900ff"},
	".txt":  {" ", "#FFFFFF"},
	".list": {" ", "#FFFFFF"},
}
