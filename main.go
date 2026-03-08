package main

import (
	"dndgoldtracker/ui"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fileName := "logFile.log"

	// open log file
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer func(){
		if err := logFile.Close(); err != nil{
			fmt.Println("Error:", fmt.Errorf("closing log file: %w", err))
		}
	}()

	// set log out put
	log.SetOutput(logFile)

	// optional: log date-time, filename, and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// Initialize and run the program
	p := tea.NewProgram(ui.NewModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", fmt.Errorf("application runtime: %w", err))
		os.Exit(1)
	}
}
