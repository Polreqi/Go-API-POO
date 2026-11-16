package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	fbauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

const cloudflareAPI = "https://api.cloudflare.com/client/v4"

// api guarda los clientes base. Las rutas se agregan en handler.
type api struct {
	auth      *fbauth.Client
	http      *http.Client
	baseURL   string
	accountID string
	apiToken  string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	cred := os.Getenv("FIREBASE_CREDENTIALS_FILE")
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if cred == "" || accountID == "" || apiToken == "" {
		return fmt.Errorf("define FIREBASE_CREDENTIALS_FILE, CLOUDFLARE_ACCOUNT_ID y CLOUDFLARE_API_TOKEN")
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(cred))
	if err != nil {
		return fmt.Errorf("firebase: %w", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("firebase auth: %w", err)
	}

	api := &api{
		auth:      authClient,
		http:      &http.Client{Timeout: 30 * time.Second},
		baseURL:   cloudflareAPI,
		accountID: accountID,
		apiToken:  apiToken,
	}

	log.Printf("escuchando en :%s", port)
	return http.ListenAndServe(":"+port, api.handler())
}

func (a *api) handler() http.Handler {
	return http.NewServeMux()
}
