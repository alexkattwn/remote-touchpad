let params = new URLSearchParams(location.search)
let token = params.get('token')

let ws = new WebSocket(`ws://${location.host}/ws?token=${token}`)

let isReady = false

let startX = 0
let startY = 0
let moved = false

let lastX = null
let lastY = null

let lastTapTime = 0
let lastTapX = 0
let lastTapY = 0

let touchCount = 0
let longPressTimer = null
let longPressTriggered = false

const errorScreen = document.getElementById('error-screen')
const errorText = document.getElementById('error-text')
const pad = document.getElementById('pad')

function showError(text) {
    errorText.innerText = text
    errorScreen.style.display = 'flex'
}

function hideError() {
    errorScreen.style.display = 'none'
}

ws.onmessage = (e) => {
    const msg = JSON.parse(e.data)

    if (msg.type === 'error') {
        showError(msg.message)
        ws.close()
        return
    }

    if (msg.type === 'ready') {
        isReady = true
        hideError()
    }
}

ws.onclose = () => {
    if (isReady) {
        showError('Соединение потеряно')
    }
}

ws.onerror = () => {
    showError('Ошибка подключения')
}

pad.addEventListener('touchstart', (e) => {
    if (!isReady) return

    touchCount = e.touches.length

    const t = e.touches[0]

    startX = t.clientX
    startY = t.clientY

    lastX = t.clientX
    lastY = t.clientY

    moved = false
    longPressTriggered = false

    // долгое нажатие (600мс)
    longPressTimer = setTimeout(() => {
        if (!moved && touchCount === 1 && !longPressTriggered) {
            ws.send(JSON.stringify({ type: 'right_click' }))
            longPressTriggered = true
        }
    }, 600)
})

pad.addEventListener('touchmove', (e) => {
    if (!isReady) return

    e.preventDefault()

    const t = e.touches[0]

    const dx = t.clientX - lastX
    const dy = t.clientY - lastY

    // считаем движение
    if (
        Math.abs(t.clientX - startX) > 10 ||
        Math.abs(t.clientY - startY) > 10
    ) {
        moved = true
        // отменяем long press если начали двигаться
        if (longPressTimer) {
            clearTimeout(longPressTimer)
        }
    }

    if (e.touches.length === 1) {
        lastX = t.clientX
        lastY = t.clientY

        ws.send(
            JSON.stringify({
                type: 'move',
                dx: dx,
                dy: dy,
            }),
        )
    }

    if (e.touches.length === 2) {
        lastY = t.clientY

        ws.send(
            JSON.stringify({
                type: 'scroll',
                dy: dy,
            }),
        )
    }
})

pad.addEventListener('touchend', (e) => {
    if (!isReady) return

    // стоп таймера
    if (longPressTimer) {
        clearTimeout(longPressTimer)
    }

    // если был long press — больше ничего не делаем
    if (longPressTriggered) return

    // если было движение — это не тап
    if (moved) return

    const now = Date.now()

    // двойной тап
    const isDoubleTap =
        now - lastTapTime < 300 &&
        Math.abs(startX - lastTapX) < 20 &&
        Math.abs(startY - lastTapY) < 20

    if (isDoubleTap) {
        ws.send(JSON.stringify({ type: 'double_click' }))
        lastTapTime = 0
        return
    }

    // обычный тап
    ws.send(JSON.stringify({ type: 'left_click' }))

    lastTapTime = now
    lastTapX = startX
    lastTapY = startY
})

window.addEventListener('beforeunload', () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close()
    }
})

// для мобильных
window.addEventListener('pagehide', () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close()
    }
})
