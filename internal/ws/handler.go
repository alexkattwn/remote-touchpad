package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"remote-touchpad/internal/mouse"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var activeConn *websocket.Conn
var mu sync.Mutex
var appToken string

const (
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

func SetToken(t string) {
	appToken = t
}

type Move struct {
	Type string  `json:"type"`
	DX   float64 `json:"dx"`
	DY   float64 `json:"dy"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка апгрейда:", err)
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("Клиент закрыл вкладку")
		return nil
	})

	conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(pongWait))

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	token := strings.TrimSpace(r.URL.Query().Get("token"))

	// неверный токен
	if token != appToken {
		fmt.Println("Неверный токен:", token)

		conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": "Неверный токен",
		})

		conn.Close()
		return
	}

	// один активный клиент
	mu.Lock()
	if activeConn != nil {
		mu.Unlock()

		fmt.Println("Уже есть активный пользователь")

		conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": "Уже подключен другой пользователь",
		})

		conn.Close()
		return
	}

	activeConn = conn
	mu.Unlock()

	// ping для поддержания соединения
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		for {
			<-ticker.C

			mu.Lock()
			if activeConn != conn {
				mu.Unlock()
				return
			}
			mu.Unlock()

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("Ping ошибка:", err)
				conn.Close()
				return
			}
		}
	}()

	fmt.Println("Пользователь подключился")

	conn.WriteJSON(map[string]string{
		"type": "ready",
	})

	defer func() {
		mu.Lock()
		if activeConn == conn {
			activeConn = nil
		}
		mu.Unlock()

		fmt.Println("Пользователь отключился")
		conn.Close()
	}()

	// основной цикл
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Клиент отключился:", err)
			break
		}

		var msg Move
		if err := json.Unmarshal(data, &msg); err != nil {
			fmt.Println("Ошибка JSON:", err)
			continue
		}

		switch msg.Type {

		case "move":
			x, y := mouse.GetPos()

			speed := 1.5

			mouse.Move(
				x+int(msg.DX*speed),
				y+int(msg.DY*speed),
			)

		case "scroll":
			mouse.Scroll(int(-msg.DY * 10))

		case "left_click":
			mouse.LeftClick()

		case "right_click":
			mouse.RightClick()

		case "double_click":
			mouse.LeftClick()
			time.Sleep(50 * time.Millisecond)
			mouse.LeftClick()
		}
	}
}