import request from '@/utils/request'
import type { Dish, Category, DailyMenu, ApiResponse, AllergenCheckResponse } from '@/types'

export interface CreateDishParams {
  name: string
  category_id: number
  price: number
  image?: string
  description?: string
  daily_limit?: number
  allergens?: string[]
  is_available?: boolean
}

export function getCategories() {
  return request.get<ApiResponse<Category[]>>('/categories')
}

export function createCategory(name: string) {
  return request.post<ApiResponse<Category>>('/categories', { name })
}

export function getDishes(params?: { category_id?: number; is_available?: boolean }) {
  return request.get<ApiResponse<Dish[]>>('/dishes', { params })
}

export function getDishById(id: number) {
  return request.get<ApiResponse<Dish>>(`/dishes/${id}`)
}

export function createDish(params: CreateDishParams) {
  return request.post<ApiResponse<Dish>>('/dishes', params)
}

export function updateDish(id: number, params: Partial<CreateDishParams>) {
  return request.put<ApiResponse<Dish>>(`/dishes/${id}`, params)
}

export function deleteDish(id: number) {
  return request.delete<ApiResponse<any>>(`/dishes/${id}`)
}

export function getDailyMenus(params?: { date?: string; meal_period?: string }) {
  return request.get<ApiResponse<DailyMenu[]>>('/menus', { params })
}

export function getDailyMenuDetail(params: { date: string; meal_period: string }) {
  return request.get<ApiResponse<DailyMenu>>('/menus/detail', { params })
}

export function getAvailableDates() {
  return request.get<ApiResponse<any[]>>('/menus/dates')
}

export interface CreateMenuParams {
  menu_date: string
  meal_period: string
  start_time: string
  end_time: string
  dish_ids: number[]
}

export function createDailyMenu(params: CreateMenuParams) {
  return request.post<ApiResponse<DailyMenu>>('/menus', params)
}

export function updateMenuDishes(menuId: number, dishIds: number[]) {
  return request.put<ApiResponse<any>>(`/menus/${menuId}/dishes`, { dish_ids: dishIds })
}

export function checkAllergens(dailyMenuDishIds: number[]) {
  return request.post<ApiResponse<AllergenCheckResponse>>('/menus/check-allergens', {
    daily_menu_dish_ids: dailyMenuDishIds
  })
}
