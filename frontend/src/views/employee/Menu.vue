<template>
  <div class="menu-page">
    <div class="date-selector">
      <el-radio-group v-model="selectedDate" size="large" @change="loadMenus">
        <el-radio-button
          v-for="date in availableDates"
          :key="date.date"
          :label="date.date"
        >
          {{ date.label }}
        </el-radio-button>
      </el-radio-group>
    </div>

    <div class="meal-tabs">
      <el-tabs v-model="activeMealPeriod" type="card" @tab-change="onMealChange">
        <el-tab-pane label="早餐" name="breakfast">
          <template #label>
            <el-icon><Sunrise /></el-icon>
            早餐
          </template>
        </el-tab-pane>
        <el-tab-pane label="午餐" name="lunch">
          <template #label>
            <el-icon><Sunny /></el-icon>
            午餐
          </template>
        </el-tab-pane>
        <el-tab-pane label="晚餐" name="dinner">
          <template #label>
            <el-icon><Moon /></el-icon>
            晚餐
          </template>
        </el-tab-pane>
      </el-tabs>
    </div>

    <div class="menu-content" v-loading="loading">
      <el-empty v-if="currentMenuDishes.length === 0" description="该时段暂无菜单" />
      
      <div class="dish-grid" v-else>
        <div class="dish-card" v-for="dish in currentMenuDishes" :key="dish.id">
          <div class="dish-image">
            <img :src="dish.dish?.image || getDefaultImage()" :alt="dish.dish?.name" />
            <el-tag v-if="dish.remaining_quantity <= 5" type="danger" class="stock-tag">
              仅剩 {{ dish.remaining_quantity }} 份
            </el-tag>
          </div>
          <div class="dish-info">
            <h3 class="dish-name">{{ dish.dish?.name }}</h3>
            <div class="dish-tags">
              <el-tag size="small" type="info" v-if="dish.dish?.category?.name">
                {{ dish.dish.category.name }}
              </el-tag>
              <el-tag
                v-for="allergen in dish.dish?.allergens"
                :key="allergen"
                size="small"
                type="warning"
              >
                {{ allergen }}
              </el-tag>
            </div>
            <p class="dish-desc" v-if="dish.dish?.description">
              {{ dish.dish.description }}
            </p>
            <div class="dish-action">
              <span class="dish-price">¥{{ dish.dish?.price?.toFixed(2) }}</span>
              <div class="quantity-control">
                <el-button
                  :disabled="dish.remaining_quantity <= 0 || getCartQuantity(dish.id) <= 0"
                  size="small"
                  circle
                  @click="removeFromCart(dish)"
                >
                  <el-icon><Minus /></el-icon>
                </el-button>
                <span class="quantity-text">{{ getCartQuantity(dish.id) }}</span>
                <el-button
                  :disabled="dish.remaining_quantity <= 0 || getCartQuantity(dish.id) >= dish.remaining_quantity"
                  size="small"
                  circle
                  type="primary"
                  @click="addToCart(dish)"
                >
                  <el-icon><Plus /></el-icon>
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="cart-bar" v-if="cartTotalCount > 0">
      <div class="cart-summary">
        <span class="cart-count">{{ cartTotalCount }} 件</span>
        <span class="cart-price">¥{{ cartTotalPrice.toFixed(2) }}</span>
      </div>
      <el-button type="primary" @click="submitOrder">去结算</el-button>
    </div>

    <el-dialog
      v-model="showAllergenDialog"
      title="过敏原警告"
      width="400px"
    >
      <p>您选择的以下菜品含有您的过敏原：</p>
      <div class="allergen-list">
        <div v-for="item in allergenWarningItems" :key="item.dish_id" class="allergen-item">
          <strong>{{ item.dish_name }}</strong>
          <el-tag v-for="a in item.allergens" :key="a" size="small" type="warning" style="margin-left: 5px;">
            {{ a }}
          </el-tag>
        </div>
      </div>
      <template #footer>
        <el-button @click="showAllergenDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmSubmit">仍要继续</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { DailyMenu, DailyMenuDish, AvailableDate, AllergenDishInfo } from '@/types'
import { getDailyMenus, getAvailableDates, checkAllergens } from '@/api/dish'
import { createOrder } from '@/api/order'

const router = useRouter()

const loading = ref(false)
const availableDates = ref<AvailableDate[]>([])
const selectedDate = ref('')
const activeMealPeriod = ref('lunch')
const menus = ref<DailyMenu[]>([])
const cart = ref<{ [key: number]: DailyMenuDish & { quantity: number } }>({})
const showAllergenDialog = ref(false)
const allergenWarningItems = ref<AllergenDishInfo[]>([])
const pendingSubmit = ref(false)

const currentMenuDishes = computed(() => {
  const menu = menus.value.find(m => m.meal_period === activeMealPeriod.value)
  return menu?.dishes || []
})

const cartTotalCount = computed(() => {
  return Object.values(cart.value).reduce((sum, item) => sum + item.quantity, 0)
})

const cartTotalPrice = computed(() => {
  return Object.values(cart.value).reduce((sum, item) => sum + (item.dish?.price || 0) * item.quantity, 0)
})

const getCartQuantity = (dailyMenuDishId: number) => {
  return cart.value[dailyMenuDishId]?.quantity || 0
}

const getDefaultImage = () => {
  return 'https://picsum.photos/200/150'
}

const loadAvailableDates = async () => {
  try {
    const res = await getAvailableDates()
    if (res.code === 200 && res.data) {
      availableDates.value = res.data
      if (availableDates.value.length > 0) {
        selectedDate.value = availableDates.value[0].date
      }
    }
  } catch (error) {
    console.error('获取可订日期失败:', error)
  }
}

const loadMenus = async () => {
  if (!selectedDate.value) return
  
  loading.value = true
  try {
    const res = await getDailyMenus({ date: selectedDate.value })
    if (res.code === 200 && res.data) {
      menus.value = res.data
      if (menus.value.length > 0 && !menus.value.find(m => m.meal_period === activeMealPeriod.value)) {
        activeMealPeriod.value = menus.value[0].meal_period
      }
    }
  } catch (error) {
    console.error('获取菜单失败:', error)
  } finally {
    loading.value = false
  }
}

const onMealChange = () => {
  cart.value = {}
}

const addToCart = (dish: DailyMenuDish) => {
  if (cart.value[dish.id]) {
    cart.value[dish.id].quantity++
  } else {
    cart.value[dish.id] = { ...dish, quantity: 1 }
  }
}

const removeFromCart = (dish: DailyMenuDish) => {
  if (cart.value[dish.id]) {
    cart.value[dish.id].quantity--
    if (cart.value[dish.id].quantity <= 0) {
      delete cart.value[dish.id]
    }
  }
}

const submitOrder = async () => {
  if (Object.keys(cart.value).length === 0) {
    ElMessage.warning('请先选择菜品')
    return
  }

  const dailyMenuDishIds = Object.keys(cart.value).map(Number)
  
  try {
    const res = await checkAllergens(dailyMenuDishIds)
    if (res.code === 200 && res.data && res.data.has_allergens) {
      allergenWarningItems.value = res.data.allergen_dishes
      showAllergenDialog.value = true
      return
    }
    
    await confirmSubmit()
  } catch (error) {
    console.error('检查过敏原失败:', error)
  }
}

const confirmSubmit = async () => {
  showAllergenDialog.value = false
  
  try {
    await ElMessageBox.confirm('确认提交订单？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
  } catch {
    return
  }

  const currentMenu = menus.value.find(m => m.meal_period === activeMealPeriod.value)
  if (!currentMenu) {
    ElMessage.error('菜单信息错误')
    return
  }

  pendingSubmit.value = true
  try {
    const now = new Date()
    const pickupStart = new Date(selectedDate.value + ' ' + currentMenu.start_time)
    const pickupEnd = new Date(selectedDate.value + ' ' + currentMenu.end_time)
    
    const res = await createOrder({
      menu_date: selectedDate.value,
      meal_period: activeMealPeriod.value,
      items: Object.values(cart.value).map(item => ({
        daily_menu_dish_id: item.id,
        quantity: item.quantity
      })),
      pickup_time_start: pickupStart.toISOString().replace('T', ' ').substring(0, 19),
      pickup_time_end: pickupEnd.toISOString().replace('T', ' ').substring(0, 19)
    })

    if (res.code === 200) {
      ElMessage.success('下单成功')
      cart.value = {}
      router.push('/employee/orders')
    }
  } catch (error) {
    console.error('下单失败:', error)
  } finally {
    pendingSubmit.value = false
  }
}

onMounted(() => {
  loadAvailableDates()
})
</script>

<style scoped>
.menu-page {
  padding-bottom: 80px;
}

.date-selector {
  margin-bottom: 20px;
  display: flex;
  justify-content: center;
}

.meal-tabs {
  margin-bottom: 20px;
}

.menu-content {
  min-height: 300px;
}

.dish-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.dish-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.dish-image {
  position: relative;
  height: 150px;
  overflow: hidden;
}

.dish-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.stock-tag {
  position: absolute;
  top: 10px;
  right: 10px;
}

.dish-info {
  padding: 15px;
}

.dish-name {
  font-size: 16px;
  font-weight: bold;
  margin: 0 0 10px 0;
  color: #333;
}

.dish-tags {
  margin-bottom: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.dish-desc {
  font-size: 13px;
  color: #909399;
  margin: 0 0 10px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dish-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dish-price {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}

.quantity-control {
  display: flex;
  align-items: center;
  gap: 10px;
}

.quantity-text {
  min-width: 30px;
  text-align: center;
  font-weight: bold;
}

.cart-bar {
  position: fixed;
  bottom: 60px;
  left: 0;
  right: 0;
  background: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
  z-index: 100;
}

.cart-summary {
  display: flex;
  align-items: center;
  gap: 15px;
}

.cart-count {
  background: #409eff;
  color: white;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.cart-price {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}

.allergen-list {
  margin-top: 10px;
}

.allergen-item {
  padding: 10px;
  background: #fdf6ec;
  border-radius: 6px;
  margin-bottom: 8px;
}
</style>
