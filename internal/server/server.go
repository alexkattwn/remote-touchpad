package server

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"remote-touchpad/internal/auth"
	"remote-touchpad/internal/utils"
	"remote-touchpad/internal/ws"

	"github.com/skip2/go-qrcode"
)

//go:embed static/*
var embeddedFiles embed.FS

var globalURL string

func Start() string {
	if globalURL != "" {
		return globalURL
	}

	content, _ := fs.Sub(embeddedFiles, "static")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(content)))
	mux.HandleFunc("/ws", ws.Handler)

	token := auth.GenerateToken(6)
	ws.SetToken(token)

	ip := utils.GetLocalIP()
	globalURL = fmt.Sprintf("http://%s:8080/?token=%s", ip, token)

	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// разрешение wails брать данные (CORS)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		// генерация QR-кода
		qrBytes, err := qrcode.Encode(globalURL, qrcode.Medium, 256)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		qrBase64 := base64.StdEncoding.EncodeToString(qrBytes)

		json.NewEncoder(w).Encode(map[string]string{
			"url": globalURL,
			"qr":  qrBase64,
		})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка сервера тачпада: %v", err)
		}
	}()

	return globalURL
}