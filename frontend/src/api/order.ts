import request from '@/utils/request'
import type { Order, ApiResponse } from '@/types'

export interface CreateOrderParams {
  menu_date: string
  meal_period: string
  items: {
    daily_menu_dish_id: number
    quantity: number
  }[]
  pickup_time_start: string
  pickup_time_end: string
}

export function createOrder(params: CreateOrderParams) {
  return request.post<ApiResponse<Order>>('/orders', params)
}

export function getMyOrders(status?: string) {
  return request.get<ApiResponse<Order[]>>('/orders', { params: { status } })
}

export function getOrderById(id: number) {
  return request.get<ApiResponse<Order>>(`/orders/${id}`)
}

export function getChefOrders(params?: { status?: string; menu_date?: string; meal_period?: string }) {
  return request.get<ApiResponse<Order[]>>('/chef/orders', { params })
}

export function getOrderByNo(orderNo: string) {
  return request.get<ApiResponse<Order>>(`/chef/orders/no/${orderNo}`)
}

export function updateOrderStatus(orderId: number, status: string) {
  return request.put<ApiResponse<any>>(`/chef/orders/${orderId}/status`, { status })
}

export function confirmPickup(orderNo: string) {
  return request.post<ApiResponse<any>>('/chef/pickup', null, { params: { order_no: orderNo } })
}
