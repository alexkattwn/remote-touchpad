//go:build windows

package mouse

import (
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSetCursorPos = user32.NewProc("SetCursorPos")
	procMouseEvent   = user32.NewProc("mouse_event")
	procGetCursorPos = user32.NewProc("GetCursorPos")
)

// ===== структура точки =====
type point struct {
	X, Y int32
}

// ===== движение =====
func Move(x, y int) {
	procSetCursorPos.Call(uintptr(x), uintptr(y))
}

// ===== клики =====
func LeftClick() {
	const (
		MOUSEEVENTF_LEFTDOWN = 0x0002
		MOUSEEVENTF_LEFTUP   = 0x0004
	)

	procMouseEvent.Call(MOUSEEVENTF_LEFTDOWN, 0, 0, 0, 0)
	procMouseEvent.Call(MOUSEEVENTF_LEFTUP, 0, 0, 0, 0)
}

func RightClick() {
	const (
		MOUSEEVENTF_RIGHTDOWN = 0x0008
		MOUSEEVENTF_RIGHTUP   = 0x0010
	)

	procMouseEvent.Call(MOUSEEVENTF_RIGHTDOWN, 0, 0, 0, 0)
	procMouseEvent.Call(MOUSEEVENTF_RIGHTUP, 0, 0, 0, 0)
}

// ===== скролл =====
func Scroll(dy int) {
	const MOUSEEVENTF_WHEEL = 0x0800

	procMouseEvent.Call(
		MOUSEEVENTF_WHEEL,
		0,
		0,
		uintptr(int32(dy)),
		0,
	)
}

// ===== позиция =====
func GetPos() (int, int) {
	var p point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	return int(p.X), int(p.Y)
}