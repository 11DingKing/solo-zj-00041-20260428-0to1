import request from '@/utils/request'
import type { DashboardStats, UserStats, ApiResponse, Review } from '@/types'

export function getDashboardStats() {
  return request.get<ApiResponse<DashboardStats>>('/admin/dashboard')
}

export function getUserStats() {
  return request.get<ApiResponse<UserStats>>('/user/stats')
}

export interface CreateReviewParams {
  dish_id: number
  rating: number
  comment?: string
}

export function createReview(orderId: number, reviews: CreateReviewParams[]) {
  return request.post<ApiResponse<any>>(`/reviews/order/${orderId}`, reviews)
}

export function getDishReviews(dishId: number, params?: { page?: number; page_size?: number }) {
  return request.get<ApiResponse<{
    total: number
    page: number
    page_size: number
    reviews: Review[]
  }>>(`/reviews/dish/${dishId}`, { params })
}
