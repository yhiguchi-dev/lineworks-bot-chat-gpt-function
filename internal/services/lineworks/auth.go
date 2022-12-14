package lineworks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AuthApi interface {
	RequestToken(accessTokenRequest AccessTokenRequest) (AccessTokenResponse, error)
	RequestTokenRefresh(tokenRefreshRequest TokenRefreshRequest) (TokenRefreshResponse, error)
}

type authApi struct {
	client *http.Client
	url    string
}

func (a authApi) RequestToken(accessTokenRequest AccessTokenRequest) (AccessTokenResponse, error) {
	form := url.Values{}
	form.Add("assertion", accessTokenRequest.Assertion)
	form.Add("grant_type", accessTokenRequest.GrantType)
	form.Add("client_id", accessTokenRequest.ClientId)
	form.Add("client_secret", accessTokenRequest.ClientSecret)
	form.Add("scope", accessTokenRequest.Scope)
	fmt.Println(form.Encode())
	body := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth2/v2.0/token", a.url), body)
	if err != nil {
		return AccessTokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := a.client.Do(req)
	if err != nil {
		return AccessTokenResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return AccessTokenResponse{}, fmt.Errorf("不正なAPIリクエストです")
	}
	fmt.Println(res)
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("fail!!")
		return AccessTokenResponse{}, err
	}
	var accessTokenResponse AccessTokenResponse
	err = json.Unmarshal(bytes, &accessTokenResponse)
	if err != nil {
		fmt.Println("fail!!")
		return AccessTokenResponse{}, err
	}
	return accessTokenResponse, nil
}

func (a authApi) RequestTokenRefresh(tokenRefreshRequest TokenRefreshRequest) (TokenRefreshResponse, error) {
	form := url.Values{}
	form.Add("refresh_token", tokenRefreshRequest.RefreshToken)
	form.Add("grant_type", tokenRefreshRequest.GrantType)
	form.Add("client_id", tokenRefreshRequest.ClientId)
	form.Add("client_secret", tokenRefreshRequest.ClientSecret)
	body := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth2/v2.0/token", a.url), body)
	if err != nil {
		return TokenRefreshResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := a.client.Do(req)
	if err != nil {
		return TokenRefreshResponse{}, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return TokenRefreshResponse{}, err
	}
	var tokenRefreshResponse TokenRefreshResponse
	err = json.Unmarshal(bytes, &tokenRefreshResponse)
	if err != nil {
		return TokenRefreshResponse{}, err
	}
	return tokenRefreshResponse, nil
}

func NewAuthApi(client *http.Client, url string) AuthApi {
	return &authApi{client, url}
}

type AccessTokenRequest struct {
	Assertion    string
	GrantType    string
	ClientId     string
	ClientSecret string
	Scope        string
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	Scope        string `json:"scope"`
}

type TokenRefreshRequest struct {
	RefreshToken string
	GrantType    string
	ClientId     string
	ClientSecret string
}

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
	Scope       string `json:"scope"`
}
