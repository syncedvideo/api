import { AxiosPromise } from 'axios'
import client from './client'

export interface UserDto {
  id: string
  name: string
  color: string
  isRoomOwner: boolean
  isAdmin: boolean
}

export function auth(): AxiosPromise<UserDto> {
  return client.post('/user/auth', undefined, { withCredentials: true })
}
