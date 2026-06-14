package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"remote-touchpad/internal/server"

	"github.com/skip2/go-qrcode"
)

type App struct {
	ctx         context.Context
	touchpadURL string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	url := server.Start()
	a.touchpadURL = url
	
	fmt.Printf("[WAILS] Сервер тачпада успешно запущен на адресе: %s\n", url)
}

// GetTouchpadData отдает ссылку и QR во фронтенд
func (a *App) GetTouchpadData() (string, string, error) {
	if a.touchpadURL == "" {
		return "", "", fmt.Errorf("URL пустой. Сервер еще не сгенерировал адрес")
	}

	qrBytes, err := qrcode.Encode(a.touchpadURL, qrcode.Medium, 256)
	if err != nil {
		return "", "", err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(qrBytes)
	return a.touchpadURL, qrBase64, nil
}

func (a *App) CloseApp() {
	os.Exit(0)
}