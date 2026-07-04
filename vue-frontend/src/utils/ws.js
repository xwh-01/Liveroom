export class LiveWSClient {
  constructor() {
    this.ws = null
    this.listeners = {}
  }

  connect(roomId, userId) {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const url = `${protocol}//${host}/ws?room_id=${roomId}&user_id=${userId}`

    this.ws = new WebSocket(url)

    this.ws.onopen = () => {
      this._emit('open')
    }

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        this._emit(msg.type, msg)
        this._emit('message', msg)
      } catch (e) {
        console.error('parse ws message error', e)
      }
    }

    this.ws.onclose = () => {
      this._emit('close')
    }

    this.ws.onerror = (err) => {
      this._emit('error', err)
    }
  }

  send(msg) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(msg))
    }
  }

  sendChat(roomId, userId, content) {
    this.send({
      type: 'chat',
      room_id: roomId,
      user_id: userId,
      data: { content }
    })
  }

  sendGift(roomId, userId, giftType) {
    this.send({
      type: 'gift',
      room_id: roomId,
      user_id: userId,
      data: { gift_type: giftType }
    })
  }

  on(event, callback) {
    if (!this.listeners[event]) {
      this.listeners[event] = []
    }
    this.listeners[event].push(callback)
  }

  _emit(event, data) {
    if (this.listeners[event]) {
      this.listeners[event].forEach(cb => cb(data))
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }
}
