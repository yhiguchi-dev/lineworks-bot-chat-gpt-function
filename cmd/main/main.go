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
	err := service.SendMessage("4643482", "71ac64b3-850c-37b4-419c-ee3646dfb974", "バンダについて教えてください")
	if err != nil {
		fmt.Println(err)
		return
	}

	//err := service.VerifySignature("H3qoCRcjAtlKLbYiGoT3bfFyLqcLhS4U6vZ8JSqzYHs=", "原神について教えてください")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
}
