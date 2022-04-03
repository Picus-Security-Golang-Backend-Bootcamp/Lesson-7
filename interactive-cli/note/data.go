package note

import (
	"log"
)

type Note struct {
	Word       string
	Definition string
	Category   string
}

var notes []Note

func OpenDatabase() error {
	notes = []Note{}
	return nil
}

func InsertNote(word string, definition string, category string) {
	note := Note{
		Word:       word,
		Definition: definition,
		Category:   category,
	}
	notes = append(notes, note)
	log.Println("Inserted study note successfully")
}

func DisplayAllNotes() {
	for _, v := range notes {
		log.Println("[", v.Category, "] ", v.Word, "—", v.Definition)
	}
}
