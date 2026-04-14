package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"remote-touchpad/cmd/internal/auth"
	"remote-touchpad/cmd/internal/utils"
	"remote-touchpad/cmd/internal/ws"

	"github.com/skip2/go-qrcode"
)

//go:embed static/*
var embeddedFiles embed.FS

func Start() {
	content, _ := fs.Sub(embeddedFiles, "static")

	http.Handle("/", http.FileServer(http.FS(content)))
	http.HandleFunc("/ws", ws.Handler)

	token := auth.GenerateToken(6)
	ws.SetToken(token)

	ip := utils.GetLocalIP()
	url := fmt.Sprintf("http://%s:8080/?token=%s", ip, token)

	fmt.Println("Открыть в браузере:")
	fmt.Println(url + "\n")

	qr, _ := qrcode.New(url, qrcode.Medium)
	fmt.Println(qr.ToSmallString(false))

	log.Fatal(http.ListenAndServe(":8080", nil))
}