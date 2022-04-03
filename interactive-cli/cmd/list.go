package cmd

import (
	"example.com/interactive_cli/note"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "See a list of all notes you've added",
	Long:  `See a list of all notes you've added`,
	Run: func(cmd *cobra.Command, args []string) {
		note.DisplayAllNotes()
	},
}

func init() {
	noteCmd.AddCommand(listCmd)
}
