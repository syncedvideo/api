import { AxiosPromise } from 'axios'
import { client, webSocketBaseURL } from './client'
import { UserDto } from './user'

export interface RoomDto {
  id: string
  name: string
  description: string
  ownerUserId: string
  users: UserDto
}

export function getRoom(roomId: string): AxiosPromise<RoomDto> {
  return client.get('/room/' + roomId)
}

export function createWebSocket(roomId: string): WebSocket {
  return new WebSocket(webSocketBaseURL + '/room/' + roomId + '/websocket')
}
