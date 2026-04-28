import request from '@/utils/request'
import type { User, ApiResponse } from '@/types'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  user: User
}

export function login(params: LoginParams) {
  return request.post<ApiResponse<LoginResult>>('/auth/login', params)
}

export function getCurrentUser() {
  return request.get<ApiResponse<User>>('/user/me')
}

export function updateAllergens(allergens: string[]) {
  return request.put<ApiResponse<any>>('/user/allergens', { allergens })
}
