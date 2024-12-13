package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Auth0Config struct {
	ClientID     string
	ClientSecret string
	Audience     string
	Domain       string
}

// reaed:invoiceスコープを持つClientAのトークンを生成
func NewClientAToken() (string, error) {
	config := Auth0Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID_A"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET_A"),
		Audience:     os.Getenv("AUTH0_AUDIENCE"),
		Domain:       os.Getenv("AUTH0_DOMAIN"),
	}
	return generateAuth0JWT("read:invoice", config)
}

// write:invoiceスコープを持つClientBのトークンを生成
func NewClientBToken() (string, error) {
	config := Auth0Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID_B"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET_B"),
		Audience:     os.Getenv("AUTH0_AUDIENCE"),
		Domain:       os.Getenv("AUTH0_DOMAIN"),
	}
	return generateAuth0JWT("write:invoice", config)
}

// スコープを持たないClientCのトークンを生成
func NewClientCToken() (string, error) {
	config := Auth0Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID_C"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET_C"),
		Audience:     os.Getenv("AUTH0_AUDIENCE"),
		Domain:       os.Getenv("AUTH0_DOMAIN"),
	}
	return generateAuth0JWT("", config)
}

func generateAuth0JWT(scope string, config Auth0Config) (string, error) {
	url := "https://" + config.Domain + "/oauth/token"

	// Payload
	payload := map[string]string{
		"client_id":     config.ClientID,
		"client_secret": config.ClientSecret,
		"audience":      config.Audience,
		"grant_type":    "client_credentials",
	}
	if scope != "" {
		payload["scope"] = scope
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response JSON
	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return response.AccessToken, nil
}

// GenerateMockJWT creates a mock JWT with the given scope
func GenerateMockJWT() (string, error) {
	secret := "test-secret"
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(), // Token expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
