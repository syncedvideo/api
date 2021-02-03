import {
  createWebSocket,
  RoomDto,
  sendAction,
  VideoDto,
  RoomActionName,
  RoomEvent,
  RoomEventName
} from '@/api'
import * as roomApi from '@/api/room'
import { reactive, Ref, ref } from 'vue'
import { useRoute } from 'vue-router'
import YouTubePlayer from 'yt-player'
import AppRoom from './AppRoom.vue'

export { AppRoom }
export default AppRoom

export const state = ref({}) as Ref<RoomEvent<RoomDto>>
export const ytPlayer: Ref<YouTubePlayer | undefined> = ref(undefined)

interface SocketState {
  connection?: WebSocket
  connected: boolean
}
export const socketState: SocketState = reactive({
  connection: undefined,
  connected: false
})

export function onOpenHandler() {
  socketState.connected = true
}

export function onErrorHandler() {
  socketState.connected = false
}

export function onMessageHandler(this: WebSocket, e: MessageEvent) {
  const event: RoomEvent = JSON.parse(e.data)
  switch (event.event) {
    case RoomEventName.Sync:
      state.value = event
      break
    case RoomEventName.PlayerSeeked:
      if (ytPlayer.value) {
        ytPlayer.value.seek(parseInt(event.data))
      }
      break
  }
}

export function connect() {
  const roomId = useRoute().params.id as string
  if (!roomId) {
    return
  }
  socketState.connection = roomApi.createWebSocket(roomId)
  socketState.connection.onopen = onOpenHandler
  socketState.connection.onerror = onErrorHandler
  socketState.connection.onclose = onErrorHandler
  socketState.connection.onmessage = onMessageHandler
}

export function disconnect() {
  socketState.connection?.close()
}

export function sendChatMessage(msg: string) {
  sendAction(socketState.connection, {
    action: RoomActionName.ChatMessage,
    data: msg
  })
}

export function sendClientColor(color: string) {
  sendAction(socketState.connection, {
    action: RoomActionName.UserSetColor,
    data: color
  })
}

export function sendClientUsername(username: string) {
  sendAction(socketState.connection, {
    action: RoomActionName.UserSetUsername,
    data: username
  })
}

export function addQueueVideo(video: VideoDto) {
  sendAction(socketState.connection, {
    action: RoomActionName.QueueAdd,
    data: video
  })
}

export function voteQueueItem(queueItemId: string) {
  sendAction(socketState.connection, {
    action: RoomActionName.QueueVote,
    data: queueItemId
  })
}

export function removeQueueItem(queueItemId: string) {
  sendAction(socketState.connection, {
    action: RoomActionName.QueueRemove,
    data: queueItemId
  })
}

// interface RoomActionMap {
//   [key: string]: {
//     [key: string]: Function
//   }
// }
export const actions = {
  player: {
    play() {
      sendAction(socketState.connection, {
        action: RoomActionName.PlayerPlay
      })
    },
    pause() {
      sendAction(socketState.connection, {
        action: RoomActionName.PlayerPause
      })
    },
    skip() {
      sendAction(socketState.connection, {
        action: RoomActionName.PlayerSkip
      })
    },
    seek(time: number) {
      sendAction(socketState.connection, {
        action: RoomActionName.PlayerSeek,
        data: time
      })
    }
  },

  queue: {
    add(video: VideoDto) {
      sendAction(socketState.connection, {
        action: RoomActionName.QueueAdd,
        data: video
      })
    },
    vote(video: VideoDto) {
      sendAction(socketState.connection, {
        action: RoomActionName.QueueVote,
        data: video.id
      })
    },
    remove(video: VideoDto) {
      sendAction(socketState.connection, {
        action: RoomActionName.QueueRemove,
        data: video.id
      })
    }
  },

  user: {
    setBuffering(buffering: boolean) {
      sendAction(socketState.connection, {
        action: RoomActionName.UserSetBuffering,
        data: buffering
      })
    }
  }
}

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
