package tui

import (
	"log"

	"github.com/charmbracelet/huh"
)

func NewConfirm(prompt string) bool {
	confirm := false

	form := huh.NewConfirm().
		Title(prompt).
		Affirmative("Yes").
		Negative("Cancel").
		Value(&confirm)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	return confirm
}
