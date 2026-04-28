import request from '@/utils/request'
import type { WeeklyMenuTemplate, ApiResponse } from '@/types'

export interface CreateTemplateParams {
  name: string
  description?: string
}

export interface TemplateItemParams {
  day_of_week: number
  meal_period: string
  dish_ids: number[]
}

export function getTemplates() {
  return request.get<ApiResponse<WeeklyMenuTemplate[]>>('/templates')
}

export function createTemplate(params: CreateTemplateParams) {
  return request.post<ApiResponse<WeeklyMenuTemplate>>('/templates', params)
}

export function updateTemplateItems(templateId: number, items: TemplateItemParams[]) {
  return request.put<ApiResponse<any>>(`/templates/${templateId}/items`, items)
}

export function applyTemplate(templateId: number, startDate: string) {
  return request.post<ApiResponse<any>>(`/templates/${templateId}/apply`, { start_date: startDate })
}
