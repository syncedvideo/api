import * as api from '@/api'
import { reactive } from 'vue'
import YouTubePlayer from 'yt-player'
import AppRoom from './AppRoom.vue'

export { AppRoom }
export default AppRoom

export interface State {
  room?: api.RoomDto
  connected: boolean
  ytPlaer?: YouTubePlayer
  chatMessages: api.ChatMessageDto[]
}

export const state: State = reactive({
  room: undefined,
  connected: false,
  ytPlaer: undefined,
  chatMessages: []
})

export function toDurationString(n: number): string {
  const minutes = Math.floor(n / 60)
  const seconds = n - minutes * 60
  let str = ''
  if (minutes < 10) {
    str += '0' + minutes + ':'
  } else {
    str += minutes + ':'
  }
  if (seconds < 10) {
    str += '0' + seconds
  } else {
    str += seconds
  }
  return str
}

export class RoomWebSocket extends WebSocket {
  roomState: State

  constructor(url: string, roomState: State) {
    super(url)
    this.roomState = roomState
  }

  onopen = (ev: Event) => {
    console.log('open', ev)
    this.roomState.connected = true
  }

  onclose = (ev: CloseEvent) => {
    console.log('closed', ev)
    this.roomState.connected = false
  }

  onmessage = (ev: MessageEvent) => {
    console.log('message', ev)
    if (ev.data === 'ping') {
      return
    }
    const msg: api.WebSocketMessage = JSON.parse(ev.data)
    switch (msg.t) {
      case api.WebSocketMessageType.Ping:
        console.log('received a PING!')
        break
      case api.WebSocketMessageType.Join:
        console.log('handle Join')
        break
      case api.WebSocketMessageType.Leave:
        console.log('handle Leave')
        break
      case api.WebSocketMessageType.SyncUsers:
        console.log('handle SyncUsers')
        break
      case api.WebSocketMessageType.Chat:
        this.handleChatMessage(msg.d)
        break
    }
  }

  onerror = (ev: Event) => {
    console.error('error', ev)
  }

  handleChatMessage(msg: api.ChatMessageDto) {
    this.roomState.chatMessages.push(msg)
  }
}
