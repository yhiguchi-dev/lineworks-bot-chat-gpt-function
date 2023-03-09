package lineworksbotchatgpt

import (
	"encoding/json"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"log"
	"net/http"
	"time"
	"yhiguchi.dev/lineworksbotchatgpt/internal/services"
)

func init() {
	functions.HTTP("MessageEventCallback", messageEventCallback)
}

var client = &http.Client{
	Timeout: 30 * time.Second,
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
func messageEventCallback(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Type   string `json:"type"`
		Source struct {
			UserId    string `json:"userId"`
			ChannelId string `json:"channelId"`
			DomainId  int    `json:"domainId"`
		} `json:"source"`
		IssuedTime time.Time `json:"issuedTime"`
		Content    struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	log.Printf("start")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	service := services.NewBotService(client)
	signature := r.Header.Get("X-WORKS-Signature")
	log.Printf("signature:" + signature)
	err = service.VerifySignature(signature, request)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	botId := r.Header.Get("X-WORKS-BotId")
	log.Printf("channelId:" + request.Source.ChannelId)
	sendMessageRequest := services.SendMessageRequest{
		BotId:     botId,
		ChannelId: request.Source.ChannelId,
		UserId:    request.Source.UserId,
	}
	err = service.SendChatMessage(sendMessageRequest, request.Content.Text)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("ok")
}
