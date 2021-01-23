import Axios, { AxiosPromise } from 'axios'

const httpBaseURL = process.env.VUE_APP_HTTP_BASE_URL
const webSocketBaseURL = process.env.VUE_APP_WEBSOCKET_BASE_URL

const client = Axios.create({ baseURL: httpBaseURL })

export interface RoomDto {
  id: string
  connectionHub: ConnectionHubDto
  player: PlayerDto
  chat: ChatDto
}

export function createRoom(): AxiosPromise<RoomDto> {
  return client.post('/room')
}

export interface ConnectionHubDto {
  connections: ConnectionMap
}

interface ConnectionDto {
  user: UserDto
}

export interface ConnectionMap {
  [userId: string]: ConnectionDto
}

export interface VideoDto {
  id: string
  providerId: string
  provider: string
  title: string
  description?: string
  duration: number
  thumbnail: string
  addedBy: UserDto
  votes: VoteMap
  statistics: {
    viewCount: number
    likeCount: number
    dislikeCount: number
  }
}

export interface VoteMap {
  [userId: string]: UserDto
}

export interface PlayerDto {
  video?: VideoDto
  time: number
  playing: boolean
  queue: VideoQueueDto
}

export interface VideoQueueDto {
  videos: VideoDto[]
}

interface ChatDto {
  messages: ChatMessageDto[]
}

export interface ChatMessageDto {
  id: string
  user: UserDto
  timestamp: string
  text: string
}

export enum VideoService {
  YouTube = 'youtube'
}

export interface UserDto {
  id: string
  username: string
  chatColor: string
  buffering: boolean
  time: number
}

export function getRoom(id: string): AxiosPromise<RoomDto> {
  return client.get('/room/' + id)
}

export function createWebSocket(roomId: string): WebSocket {
  return new WebSocket(webSocketBaseURL + '/room/' + roomId)
}

// RoomEvent is broadcasted by WebSocket
export interface RoomEvent<T = any> {
  event: RoomEventName
  user: UserDto
  data: T
}

// RoomAction is sent to WebSocket
export interface RoomAction<T = any> {
  action: RoomActionName
  data?: T
}

export enum RoomEventName {
  Sync = 'sync',
  PlayerSeeked = 'player:seeked'
}

export enum RoomActionName {
  // User actions
  UserSetBuffering = 'user:set:buffering',
  UserSetColor = 'user:set:color',
  UserSetUsername = 'user:set:username',

  // Player actions
  PlayerPlay = 'player:play',
  PlayerPause = 'player:pause',
  PlayerSkip = 'player:skip',
  PlayerSeek = 'player:seek',

  // Chat actions
  ChatMessage = 'chat:message',

  // Video queue actions
  QueueAdd = 'queue:add',
  QueueRemove = 'queue:remove',
  QueueVote = 'queue:vote'
}

export function sendAction(
  connection: WebSocket | undefined,
  action: RoomAction
) {
  if (connection) {
    connection.send(JSON.stringify(action))
  }
}

export interface VideoSearchResponseDto {
  query: string
  videos: VideoDto[]
}

export function search(query: string): AxiosPromise<VideoSearchResponseDto> {
  return client.get('/search/youtube?query=' + query)
}
