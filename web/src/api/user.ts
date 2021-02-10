import { UserDto } from '@/api'
import { AxiosPromise } from 'axios'
import client from './client'

export function auth(): AxiosPromise<UserDto> {
  return client.post('/auth')
}

export function updateUser(user: UserDto): AxiosPromise<UserDto> {
  return client.patch('/user', user)
}
