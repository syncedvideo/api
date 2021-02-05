import { UserDto } from '@/api'
import { AxiosPromise } from 'axios'
import client from './client'

export function auth(): AxiosPromise<UserDto> {
  return client.post('/user/auth', undefined, { withCredentials: true })
}
