package main

import (
	"context"
	"embed"
	"os"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"golang.org/x/sys/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed icon.ico
var iconBytes []byte

const mutexName = "Global\\RemoteTouchpadMutex"

func main() {
var lpName *uint16
	var handle windows.Handle
	var mutexErr error

	lpName, mutexErr = windows.UTF16PtrFromString(mutexName)
	if mutexErr != nil {
		os.Exit(1)
	}

	handle, mutexErr = windows.CreateMutex(nil, false, lpName)
	if mutexErr != nil {
		if mutexErr == windows.ERROR_ALREADY_EXISTS {
			os.Exit(0)
		}
		os.Exit(1)
	}
	
	defer windows.CloseHandle(handle)

	app := NewApp()

	windowWidth := 380
	windowHeight := 420

	err := wails.Run(&options.App{
		Title:         "Remote Touchpad",
		Width:         windowWidth,
		Height:        windowHeight,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 245, G: 245, B: 245, A: 255},

		OnStartup: func(ctx context.Context) {
			app.startup(ctx)

			// запуск в системном трее
			go func() {
				systray.Run(func() {
					systray.SetIcon(iconBytes)
					systray.SetTooltip("Remote Touchpad")

					mOpen := systray.AddMenuItem("Показать", "Открыть окно")
					mQuit := systray.AddMenuItem("Выход", "Закрыть полностью")

					for {
						select {
						case <-mOpen.ClickedCh:
							runtime.WindowShow(ctx)
						
						case <-mQuit.ClickedCh:
							systray.Quit()
							os.Exit(0)
							return
						}
					}
				}, func() {})
			}()
		},

		OnBeforeClose: func(ctx context.Context) bool {
			runtime.WindowHide(ctx)
			return true
		},

		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}
}