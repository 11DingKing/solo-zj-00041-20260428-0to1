export interface User {
  id: number
  username: string
  name: string
  role: 'admin' | 'chef' | 'employee'
  avatar?: string
  allergens?: UserAllergen[]
  created_at: string
  updated_at: string
}

export interface UserAllergen {
  id: number
  user_id: number
  allergen: string
  created_at: string
}

export interface Category {
  id: number
  name: string
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Dish {
  id: number
  name: string
  category_id: number
  category?: Category
  price: number
  image?: string
  description?: string
  daily_limit: number
  allergens: string[]
  is_available: boolean
  created_at: string
  updated_at: string
}

export interface DailyMenu {
  id: number
  menu_date: string
  meal_period: 'breakfast' | 'lunch' | 'dinner'
  start_time: string
  end_time: string
  is_active: boolean
  dishes?: DailyMenuDish[]
  created_at: string
  updated_at: string
}

export interface DailyMenuDish {
  id: number
  daily_menu_id: number
  dish_id: number
  dish?: Dish
  remaining_quantity: number
  created_at: string
}

export interface Order {
  id: number
  order_no: string
  user_id: number
  user?: User
  total_amount: number
  status: OrderStatus
  pickup_time_start: string
  pickup_time_end: string
  pickup_code?: string
  qr_code_content?: string
  items?: OrderItem[]
  created_at: string
  updated_at: string
}

export type OrderStatus = 
  | 'pending_confirm'
  | 'in_production' 
  | 'ready_for_pickup'
  | 'picked_up'
  | 'reviewed'

export interface OrderItem {
  id: number
  order_id: number
  dish_id: number
  dish_name: string
  dish_price: number
  quantity: number
  subtotal: number
  created_at: string
}

export interface Review {
  id: number
  order_id: number
  user_id: number
  dish_id: number
  dish?: Dish
  rating: number
  comment?: string
  created_at: string
}

export interface WeeklyMenuTemplate {
  id: number
  name: string
  description?: string
  is_active: boolean
  items?: WeeklyMenuTemplateItem[]
  created_at: string
  updated_at: string
}

export interface WeeklyMenuTemplateItem {
  id: number
  template_id: number
  day_of_week: number
  meal_period: string
  dish_id: number
  dish?: Dish
  created_at: string
}

export interface DashboardStats {
  today_orders: number
  today_revenue: number
  meal_period_distribution: MealPeriodCount[]
  top_dishes: TopDish[]
  last_30_days_revenue: DailyRevenue[]
  top_rated_dishes: RatedDish[]
  tomorrow_ingredients: IngredientDemand[]
}

export interface MealPeriodCount {
  meal_period: string
  count: number
}

export interface TopDish {
  dish_id: number
  dish_name: string
  total_orders: number
  total_revenue: number
}

export interface DailyRevenue {
  date: string
  revenue: number
}

export interface RatedDish {
  dish_id: number
  dish_name: string
  avg_rating: number
  review_count: number
}

export interface IngredientDemand {
  dish_id: number
  dish_name: string
  quantity: number
}

export interface UserStats {
  month_amount: number
  order_count: number
  top_dishes: UserTopDish[]
}

export interface UserTopDish {
  dish_id: number
  dish_name: string
  count: number
}

export interface AllergenCheckResponse {
  has_allergens: boolean
  allergen_dishes: AllergenDishInfo[]
}

export interface AllergenDishInfo {
  dish_id: number
  dish_name: string
  allergens: string[]
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface CartItem {
  daily_menu_dish_id: number
  dish: Dish
  daily_menu_dish: DailyMenuDish
  quantity: number
}

export interface AvailableDate {
  date: string
  label: string
  weekday: string
}
