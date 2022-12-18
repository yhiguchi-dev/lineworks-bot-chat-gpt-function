package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"yhiguchi.dev/lineworksbotchatgpt/internal/services/lineworks"
	"yhiguchi.dev/lineworksbotchatgpt/internal/services/openai"
)

type BotService interface {
	VerifySignature(signature string, body any) error
	SendMessage(botId, channelId, message string) error
}

type botService struct {
	client *http.Client
}

func (b botService) VerifySignature(signature string, body any) error {
	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	botSecret := os.Getenv("BOT_SECRET")
	verifier := lineworks.NewVerifier(botSecret, signature)
	if verifier.Verify(bytes) {
		return fmt.Errorf("不正なリクエストです")
	}
	return nil
}

func (b botService) SendMessage(botId, channelId, message string) error {
	if strings.HasPrefix(message, "/off") {
		return nil
	}
	privateKey := os.Getenv("PRIVATE_KEY")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	serviceAccountId := os.Getenv("SERVICE_ACCOUNT_ID")
	openaiApiKey := os.Getenv("OPENAI_API_KEY")

	openaiApi := openai.NewApi(b.client, "https://api.openai.com", openaiApiKey)
	completionRequest := openai.CompletionRequest{
		Model:     "text-davinci-003",
		Prompt:    message,
		MaxTokens: 2000,
	}
	completionResponse, err := openaiApi.Completions(completionRequest)
	if err != nil {
		return err
	}

	key, err := lineworks.CreatePrivateKey(privateKey)
	if err != nil {
		return err
	}
	signer := lineworks.NewSigner(key)

	jwt := lineworks.NewJwt(signer, clientId, serviceAccountId)
	value, err := jwt.Create()
	if err != nil {
		return err
	}

	authApi := lineworks.NewAuthApi(b.client, "https://auth.worksmobile.com")

	accessTokenRequest := lineworks.AccessTokenRequest{
		Assertion:    value,
		GrantType:    "urn:ietf:params:oauth:grant-type:jwt-bearer",
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        "bot",
	}

	token, err := authApi.RequestToken(accessTokenRequest)
	lineworksApi := lineworks.NewApi(b.client, token.AccessToken, "https://www.worksapis.com")
	choices := completionResponse.Choices
	fmt.Println(choices[0].Text)
	messageRequest := lineworks.MessageRequest{
		Content: struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}(struct {
			Type string
			Text string
		}{Type: "text", Text: choices[0].Text}),
	}
	err = lineworksApi.Messages(botId, channelId, messageRequest)
	if err != nil {
		return err
	}
	return nil
}

func NewBotService(client *http.Client) BotService {
	return &botService{client}
}
