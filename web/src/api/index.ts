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
  currentUser: UserDto
  users: UserDto[]
}

export enum WebSocketMessageType {
  Ping = 0,
  Join = 1000,
  Leave = 1001,
  Chat = 2000,
  SyncUsers = 3000
}

export interface WebSocketMessage<T = WebSocketMessageType, D = any> {
  t: T
  d: D
}

export interface ChatMessageDto {
  id: string
  text: string
  user: UserDto
}

export interface VideoDto {
  id: string
  provider: number
  providerId: string
  title: string
  author: string
  thumbnail: string
  statistics: {
    views: number
    likes: number
    dislikes: number
    lengthInSeconds: number
  }
}
