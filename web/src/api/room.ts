import { ChatMessageDto, RoomDto, VideoDto } from '@/api'
import { AxiosPromise } from 'axios'
import { client, webSocketBaseURL } from './client'

export function createRoom(): AxiosPromise<RoomDto> {
  return client.post('/room')
}

export function getRoom(roomId: string): AxiosPromise<RoomDto> {
  return client.get('/room/' + roomId)
}

export function getRoomWebSocketUrl(roomId: string): string {
  return webSocketBaseURL + '/room/' + roomId + '/websocket'
}

export function createWebSocket(roomId: string): WebSocket {
  return new WebSocket(webSocketBaseURL + '/room/' + roomId + '/websocket')
}

export function sendChatMessage(
  roomId: string,
  message: string
): AxiosPromise<ChatMessageDto> {
  return client.post('/room/' + roomId + '/chat', { message })
}

export function getVideoInfo(
  roomId: string,
  videoUrl: string
): AxiosPromise<VideoDto> {
  return client.get('/room/' + roomId + '/video-info', {
    params: { url: videoUrl }
  })
}
