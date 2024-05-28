package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Ayobami0/cli-chat/pb"
	"github.com/Ayobami0/cli-chat/ui"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	USAGE        = "Usage: cli-chat <command> [arguments]"
	LOGIN_USAGE  = "Usage: cli-chat login [[-h | --help] | [-u | --username] <username>]"
	CREATE_USAGE = "Usage: cli-chat create [-h | --help]"
	HELP         = "Chat with friends from you terminal.\n\n%s\n\nAvaliable Commands:\n\tcreate: create a new account to chat with\n\tlogin: log into an existing account\n"
)

func main() {
	var username string
	var help bool

	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	loginCmd.Usage = func() {
		fmt.Printf("cli-chat: invalid argument to login\n%s\n", LOGIN_USAGE)
		return
	}

	loginCmd.StringVar(&username, "u", "", "username")
	loginCmd.StringVar(&username, "username", "", "username")

	loginCmd.BoolVar(&help, "h", false, "help")
	loginCmd.BoolVar(&help, "help", false, "help")

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createCmd.Usage = func() {
		fmt.Printf("cli-chat: invalid argument to create\n%s\n", CREATE_USAGE)
		return
	}
	createCmd.BoolVar(&help, "h", false, "help")
	createCmd.BoolVar(&help, "help", false, "help")

	if len(os.Args) < 2 {
		fmt.Printf("cli-chat: invalid argument\n%s\n", USAGE)
		return
	}

	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		fmt.Println("No server address specified")
		return
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	switch os.Args[1] {
	case "login":
		loginCmd.Parse(os.Args[2:])
		if help {
			fmt.Printf("Login command.\n\n%s\n\nArguments:\n\t-h, --help: show help\n\t-u, --username: account username\n", LOGIN_USAGE)
			return
		}
		if username == "" {
			fmt.Printf("cli-chat: invalid argument passed to <username>.\n\n%s\n", LOGIN_USAGE)
			return
		}
		if _, err := tea.NewProgram(ui.NewLoginModel(username, client)).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
		}
	case "create":
		createCmd.Parse(os.Args[2:])
		if help {
			fmt.Printf("Create command.\n\n%s\n\nArguments:\n\t-h, --help: show help\n", CREATE_USAGE)
			return
		}
		if _, err := tea.NewProgram(ui.NewCreateModel(client)).Run(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
		}
	case "-h", "--help":
		fmt.Printf(HELP, USAGE)
	default:
		fmt.Printf("cli-chat: invalid argument %s\n%s\n", os.Args[1], USAGE)
	}
}
