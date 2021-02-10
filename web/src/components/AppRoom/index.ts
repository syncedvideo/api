import * as api from '@/api'
import { reactive } from 'vue'
import YouTubePlayer from 'yt-player'
import AppRoom from './AppRoom.vue'

export { AppRoom }
export default AppRoom

export interface State {
  room?: api.RoomDto
  connected: boolean
  ytPlayer?: YouTubePlayer
  chatMessages: api.ChatMessageDto[]
}

export function createState(): State {
  return {
    room: undefined,
    connected: false,
    ytPlayer: undefined,
    chatMessages: []
  }
}

export const state: State = reactive(createState())

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
  state: State

  constructor(url: string, state: State) {
    super(url)
    this.state = state
  }

  onopen = (ev: Event) => {
    console.log('open', ev)
    this.state.connected = true
  }

  onclose = (ev: CloseEvent) => {
    console.log('closed', ev)
    this.state.connected = false
  }

  onmessage = (ev: MessageEvent) => {
    console.log('message', ev)
    if (ev.data === 'ping') {
      return
    }
    const msg: api.WebSocketMessage = JSON.parse(ev.data)
    if (msg.t === undefined) {
      return
    }
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
        this.handleSyncUsers(msg.d)
        break
      case api.WebSocketMessageType.Chat:
        this.handleChatMessage(msg.d)
        break
    }
  }

  onerror = (ev: Event) => {
    console.error('error', ev)
  }

  handleSyncUsers(users: api.UserDto[]) {
    if (this.state.room) {
      this.state.room.users = users
    }
  }

  handleChatMessage(msg: api.ChatMessageDto) {
    this.state.chatMessages.push(msg)
  }
}
