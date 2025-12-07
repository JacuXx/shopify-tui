package main

import (
	"os"
	"strings"
)

type IconSet struct {
	Success string
	Error   string
	Warning string
	Info    string

	Login  string
	Add    string
	Delete string
	Exit   string

	Store     string
	Server    string
	ServerOn  string
	ServerOff string
	Logs      string
	Stop      string

	Download string
	Upload   string
	Git      string
	Folder   string

	Terminal string
	Editor   string
	Code     string

	Arrow  string
	Dot    string
	Play   string
	Rocket string

	App string
}

var NerdIcons = IconSet{

	Success: "",
	Error:   "",
	Warning: "",
	Info:    "",

	Login:  "",
	Add:    "",
	Delete: "",
	Exit:   "",

	Store:     "",
	Server:    "",
	ServerOn:  "",
	ServerOff: "",
	Logs:      "",
	Stop:      "",

	Download: "",
	Upload:   "",
	Git:      "",
	Folder:   "",

	Terminal: "",
	Editor:   "",
	Code:     "",

	Arrow:  "",
	Dot:    "",
	Play:   "",
	Rocket: "",

	App: "",
}

var ASCIIIcons = IconSet{

	Success: "[OK]",
	Error:   "[X]",
	Warning: "[!]",
	Info:    "[i]",

	Login:  "[>]",
	Add:    "[+]",
	Delete: "[-]",
	Exit:   "[Q]",

	Store:     "[S]",
	Server:    "[#]",
	ServerOn:  "(*)  ",
	ServerOff: "( )",
	Logs:      "[=]",
	Stop:      "[X]",

	Download: "[v]",
	Upload:   "[^]",
	Git:      "[G]",
	Folder:   "[D]",

	Terminal: "[$]",
	Editor:   "[E]",
	Code:     "[<>]",

	Arrow:  ">",
	Dot:    "*",
	Play:   "[>]",
	Rocket: "[!]",

	App: ">>",
}

var Icons = ASCIIIcons

func DetectNerdFont() bool {

	if env := os.Getenv("SHOPIFY_TUI_ICONS"); env != "" {
		return strings.ToLower(env) == "nerd" || strings.ToLower(env) == "true"
	}

	if env := os.Getenv("SHOPIFY_TUI_ASCII"); env != "" {
		return false
	}

	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")

	nerdTerminals := []string{
		"kitty", "alacritty", "wezterm", "iterm", "iterm2",
		"hyper", "warp", "ghostty",
	}

	for _, t := range nerdTerminals {
		if strings.Contains(strings.ToLower(term), t) ||
			strings.Contains(strings.ToLower(termProgram), t) {
			return true
		}
	}

	return true
}

func InitIcons() {
	if DetectNerdFont() {
		Icons = NerdIcons
	} else {
		Icons = ASCIIIcons
	}
}

func IconSuccess(msg string) string {
	return Icons.Success + " " + msg
}

func IconError(msg string) string {
	return Icons.Error + " " + msg
}

func IconWarning(msg string) string {
	return Icons.Warning + " " + msg
}

func IconInfo(msg string) string {
	return Icons.Info + " " + msg
}
