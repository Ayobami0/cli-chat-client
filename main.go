package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Ayobami0/cli-chat/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	var username string

	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)

	if len(os.Args) < 2 {
		os.Exit(-1)
	}

	switch os.Args[1] {
	case "login":
		if len(os.Args) == 2 {
			// TODO: Work on retriving username from a local storage if exist
		}
		loginCmd.StringVar(&username, "u", "", "Username")
		loginCmd.Parse(os.Args[2:])
		// TODO: Add passwordless login to config
		if _, err := tea.NewProgram(ui.NewLoginModel(username, "immanuel")).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	case "create":
		if len(os.Args) > 2 {
			// TODO: Implement a usage print statement
			os.Exit(-1)
		}
		if _, err := tea.NewProgram(ui.NewCreateModel()).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}
