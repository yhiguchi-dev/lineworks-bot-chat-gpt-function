package main

import (
	"fmt"
	"net/http"
	"time"
	"yhiguchi.dev/lineworksbotchatgpt/internal/services"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

func main() {

	service := services.NewBotService(client)
	sendMessageRequest := services.SendMessageRequest{
		BotId:     "botId",
		ChannelId: "request.Source.ChannelId",
		UserId:    "request.Source.UserId",
	}
	err := service.SendMessage(sendMessageRequest, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	//err := service.VerifySignature("=", "")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
}
