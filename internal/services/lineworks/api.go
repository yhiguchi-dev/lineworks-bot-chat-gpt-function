package lineworks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Api interface {
	Messages(botId string, channelId string, messageRequest MessageRequest) error
}

type api struct {
	client      *http.Client
	accessToken string
	url         string
}

func (a api) Messages(botId string, channelId string, messageRequest MessageRequest) error {
	requestBytes, err := json.Marshal(messageRequest)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1.0/bots/%s/channels/%s/messages", a.url, botId, channelId),
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(res.StatusCode)
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("不正なAPIリクエストです")
	}
	return nil
}

func NewApi(client *http.Client, accessToken string, url string) Api {
	return &api{client, accessToken, url}
}

type MessageRequest struct {
	Content struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}
