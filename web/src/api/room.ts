import { RoomDto } from '@/api'
import { AxiosPromise } from 'axios'
import { client, webSocketBaseURL } from './client'

export function getRoom(roomId: string): AxiosPromise<RoomDto> {
  return client.get('/room/' + roomId)
}

export function createWebSocket(roomId: string): WebSocket {
  return new WebSocket(webSocketBaseURL + '/room/' + roomId + '/websocket')
}
