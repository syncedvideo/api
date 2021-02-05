export interface UserDto {
  id: string
  name: string
  color: string
  isRoomOwner: boolean
  isAdmin: boolean
}

export interface RoomDto {
  id: string
  name: string
  description: string
  ownerUserId: string
  users: UserDto[]
}

export enum WebSocketMessageType {
  Join = 1000,
  Leave = 1001,
  Chat = 2000,
  SyncUsers = 3000
}

export interface WebSocketMessage<T = WebSocketMessageType, D = any> {
  t: T
  d: D
}

export interface ChatMessage {
  id: string
  text: string
  user: UserDto
}
