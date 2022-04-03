package main

import (
	"example.com/interactive_cli/cmd"
	"example.com/interactive_cli/note"
)

func main() {
	note.OpenDatabase()
	cmd.Execute()
}
