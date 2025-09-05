package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func (config *apiConfig) startOauthFlow(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?"+
			"client_id=%s&redirect_uri=%s&response_type=code&"+
			"scope=%s&access_type=offline&prompt=consent",
		config.clientID,
		url.QueryEscape("http://localhost:8080/auth/callback"),
		url.QueryEscape("https://www.googleapis.com/auth/calendar.events"))

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}
func (cfg *apiConfig) exchangeCodeForTokens(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", cfg.clientID)
	data.Set("client_secret", cfg.clientSecret)
	data.Set("redirect_uri", "http://localhost:8080/auth/callback")
	data.Set("grant_type", "authorization_code")

	resp, err := http.Post("https://oauth2.googleapis.com/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", body)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (config *apiConfig) handleOauthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No authorization code received", http.StatusBadRequest)
		return
	}

	tokens, err := config.exchangeCodeForTokens(code)
	if err != nil {
		http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	config.accessToken = "Bearer " + tokens.AccessToken
	config.refreshToken = tokens.RefreshToken

	configDir, err := os.UserConfigDir()
	if err != nil {
		http.Error(w, "Faild to get config Dir", http.StatusInternalServerError)
		return
	}
	url := configDir + "/go-home/.env"
	envMap, err := godotenv.Read(url)
	if err != nil {
		log.Fatal("Error reading .env file:", err)
	}
	envMap["ACCESS_TOKEN"] = tokens.AccessToken
	envMap["REFRESH_TOKEN"] = config.refreshToken
	err = godotenv.Write(envMap, url)
	if err != nil {
		log.Fatal("Error writing to .env file:", err)
	}

	w.Write([]byte("Authorization successful. You can close this window."))
}
