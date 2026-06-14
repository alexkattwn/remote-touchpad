const qrImage = document.getElementById('qr-image')
const touchpadLink = document.getElementById('touchpad-link')
const closeBtn = document.getElementById('close-btn')

let originalURL = ''

async function loadTouchpadData() {
    try {
        const response = await fetch('http://127.0.0.1:8080/api/data')
        if (!response.ok) throw new Error('Ошибка сервера')

        const data = await response.json()

        originalURL = data.url
        touchpadLink.textContent = data.url
        qrImage.src = `data:image/png;base64,${data.qr}`
    } catch (err) {
        setTimeout(loadTouchpadData, 300)
    }
}

touchpadLink.addEventListener('click', async () => {
    if (!originalURL) return

    try {
        await navigator.clipboard.writeText(originalURL)

        touchpadLink.textContent = ' ССЫЛКА СКОПИРОВАНА! '
        touchpadLink.style.color = '#ff0055'
        touchpadLink.style.borderColor = '#00ffcc'

        setTimeout(() => {
            touchpadLink.textContent = originalURL
            touchpadLink.style.color = '#00ffcc'
            touchpadLink.style.borderColor = '#ff0055'
        }, 1500)
    } catch (err) {
        console.error('Не удалось скопировать:', err)
    }
})

window.addEventListener('DOMContentLoaded', loadTouchpadData)

closeBtn.addEventListener('click', () => {
    if (
        window.go &&
        window.go.main &&
        window.go.main.App &&
        window.go.main.App.CloseApp
    ) {
        window.go.main.App.CloseApp()
    } else {
        if (window.runtime && window.runtime.Quit) {
            window.runtime.Quit()
        }
    }
})
