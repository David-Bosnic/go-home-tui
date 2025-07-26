package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func (config *apiConfig) refreshAccessTokenPost(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("GET /admin/refresh/ Error loading .env", err)
	}
	tokenURL := "https://oauth2.googleapis.com/token"
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	refreshToken := os.Getenv("REFRESH_TOKEN")

	if clientID == "" || clientSecret == "" || refreshToken == "" {
		log.Println("POST /auth/refreshToken .env did not match requirements")
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Println("POST /auth/refreshToken failed to create request", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("POST /auth/refreshToken failed to make request: %w", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("POST /auth/refreshToken failed to read response: %w", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("POST /auth/refresh token refresh failed with status %d: %s\n", resp.StatusCode, string(body))
		return
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Println("POST /auth/refresh failed to parse JSON %w", err)
		return
	}
	config.accessToken = "Bearer " + tokenResp.AccessToken
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error reading .env file:", err)
	}
	envMap["ACCESS_TOKEN"] = tokenResp.AccessToken
	err = godotenv.Write(envMap, ".env")
	if err != nil {
		log.Fatal("Error writing to .env file:", err)
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Successfully refreshed access token"))
}
