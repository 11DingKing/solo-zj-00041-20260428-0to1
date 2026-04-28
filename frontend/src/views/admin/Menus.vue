<template>
  <div class="menus-page">
    <div class="page-header">
      <h2>菜单管理</h2>
      <div class="header-actions">
        <el-date-picker
          v-model="selectedDate"
          type="date"
          placeholder="选择日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          :clearable="false"
          @change="loadMenus"
        />
        <el-button type="primary" :icon="Plus" @click="handleCreate">新建菜单</el-button>
      </div>
    </div>

    <div class="menus-container" v-loading="loading">
      <el-empty v-if="menus.length === 0" description="当日暂无菜单，请点击「新建菜单」添加" />
      
      <div v-else class="period-cards">
        <div class="period-card" v-for="menu in menus" :key="menu.id">
          <div class="period-header" :class="menu.meal_period">
            <h3>{{ getMealPeriodText(menu.meal_period) }}</h3>
            <span class="time">{{ menu.start_time }} - {{ menu.end_time }}</span>
            <el-tag :type="menu.is_active ? 'success' : 'info'" size="small">
              {{ menu.is_active ? '启用' : '禁用' }}
            </el-tag>
          </div>
          
          <div class="period-dishes" v-if="menu.dishes && menu.dishes.length > 0">
            <div class="dish-item" v-for="item in menu.dishes" :key="item.id">
              <div class="dish-info">
                <img :src="item.dish?.image || defaultImage" :alt="item.dish?.name" />
                <div class="dish-detail">
                  <span class="name">{{ item.dish?.name }}</span>
                  <span class="price">¥{{ item.dish?.price?.toFixed(2) }}</span>
                </div>
              </div>
              <div class="dish-quantity">
                <span>剩余：{{ item.remaining_quantity }} 份</span>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无菜品" :image-size="60" />
          
          <div class="period-actions">
            <el-button size="small" @click="handleEdit(menu)">编辑菜品</el-button>
            <el-button 
              size="small" 
              :type="menu.is_active ? 'warning' : 'success'"
              @click="toggleActive(menu)"
            >
              {{ menu.is_active ? '禁用' : '启用' }}
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <el-dialog 
      v-model="showDialog" 
      :title="isEdit ? '编辑菜单菜品' : '新建菜单'" 
      width="700px"
      :close-on-click-modal="false"
    >
      <el-form 
        ref="formRef" 
        :model="formData" 
        :rules="formRules" 
        label-width="100px"
        v-if="!isEdit"
      >
        <el-form-item label="菜单日期" prop="menu_date">
          <el-date-picker
            v-model="formData.menu_date"
            type="date"
            placeholder="选择日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="时段" prop="meal_period">
          <el-radio-group v-model="formData.meal_period">
            <el-radio-button label="breakfast">早餐</el-radio-button>
            <el-radio-button label="lunch">午餐</el-radio-button>
            <el-radio-button label="dinner">晚餐</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="供应时间">
          <el-time-select
            v-model="formData.start_time"
            placeholder="开始时间"
            :picker-options="{ start: '00:00', step: '00:15', end: '23:45' }"
            style="width: 180px"
          />
          <span class="time-separator">至</span>
          <el-time-select
            v-model="formData.end_time"
            placeholder="结束时间"
            :picker-options="{ start: '00:00', step: '00:15', end: '23:45' }"
            style="width: 180px"
          />
        </el-form-item>
      </el-form>

      <div class="dish-selection">
        <h4>选择菜品</h4>
        <el-transfer
          v-model="selectedDishes"
          :data="dishesData"
          :titles="['可选菜品', '已选菜品']"
          :button-texts="['添加', '移除']"
          filterable
          filter-placeholder="搜索菜品"
          :props="{
            key: 'id',
            label: 'label'
          }"
        />
      </div>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules, type TransferDataItem } from 'element-plus'
import type { DailyMenu, DailyMenuDish, Dish } from '@/types'
import { getDailyMenus, createDailyMenu, updateMenuDishes, getDishes } from '@/api/dish'

const loading = ref(false)
const submitting = ref(false)
const menus = ref<DailyMenu[]>([])
const dishes = ref<Dish[]>([])
const selectedDate = ref(new Date().toISOString().split('T')[0])

const showDialog = ref(false)
const isEdit = ref(false)
const editMenuId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const defaultImage = 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=food%20dish%20placeholder%20image&image_size=square'

const formData = reactive({
  menu_date: '',
  meal_period: 'lunch' as string,
  start_time: '11:00',
  end_time: '13:00'
})

const formRules: FormRules = {
  menu_date: [{ required: true, message: '请选择日期', trigger: 'change' }],
  meal_period: [{ required: true, message: '请选择时段', trigger: 'change' }]
}

const selectedDishes = ref<number[]>([])

const dishesData = computed<TransferDataItem[]>(() => {
  return dishes.value.map(dish => ({
    id: dish.id,
    label: `${dish.name} (¥${dish.price.toFixed(2)})`,
    disabled: false
  }))
})

const getMealPeriodText = (period: string) => {
  const map: { [key: string]: string } = {
    breakfast: '早餐',
    lunch: '午餐',
    dinner: '晚餐'
  }
  return map[period] || period
}

const loadMenus = async () => {
  loading.value = true
  try {
    const res = await getDailyMenus({ date: selectedDate.value })
    if (res.code === 200 && res.data) {
      menus.value = res.data
    }
  } catch (error) {
    console.error('获取菜单失败:', error)
  } finally {
    loading.value = false
  }
}

const loadDishes = async () => {
  try {
    const res = await getDishes({ is_available: true })
    if (res.code === 200 && res.data) {
      dishes.value = res.data
    }
  } catch (error) {
    console.error('获取菜品列表失败:', error)
  }
}

const resetForm = () => {
  formData.menu_date = selectedDate.value
  formData.meal_period = 'lunch'
  formData.start_time = '11:00'
  formData.end_time = '13:00'
  selectedDishes.value = []
  editMenuId.value = null
  isEdit.value = false
}

const handleCreate = () => {
  resetForm()
  formData.menu_date = selectedDate.value
  showDialog.value = true
}

const handleEdit = (menu: DailyMenu) => {
  isEdit.value = true
  editMenuId.value = menu.id
  selectedDishes.value = menu.dishes?.map(d => d.dish_id) || []
  showDialog.value = true
}

const toggleActive = async (menu: DailyMenu) => {
  // TODO: 实现启用/禁用菜单的API调用
  ElMessage.info('该功能开发中')
}

const handleSubmit = async () => {
  if (selectedDishes.value.length === 0) {
    ElMessage.warning('请至少选择一个菜品')
    return
  }

  submitting.value = true
  try {
    if (isEdit.value && editMenuId.value) {
      const res = await updateMenuDishes(editMenuId.value, selectedDishes.value)
      if (res.code === 200) {
        ElMessage.success('更新成功')
        showDialog.value = false
        loadMenus()
      } else {
        ElMessage.error(res.message || '更新失败')
      }
    } else {
      if (!formRef.value) return
      
      await formRef.value.validate(async (valid) => {
        if (!valid) {
          submitting.value = false
          return
        }

        const res = await createDailyMenu({
          menu_date: formData.menu_date!,
          meal_period: formData.meal_period,
          start_time: formData.start_time,
          end_time: formData.end_time,
          dish_ids: selectedDishes.value
        })

        if (res.code === 200) {
          ElMessage.success('创建成功')
          showDialog.value = false
          loadMenus()
        } else {
          ElMessage.error(res.message || '创建失败')
        }
      })
    }
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadDishes()
  loadMenus()
})
</script>

<style scoped>
.menus-page {
  min-height: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 15px;
  align-items: center;
}

.menus-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  min-height: 400px;
}

.period-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.period-card {
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  overflow: hidden;
}

.period-header {
  padding: 15px 20px;
  display: flex;
  align-items: center;
  gap: 15px;
}

.period-header.breakfast {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.period-header.lunch {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.period-header.dinner {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.period-header h3 {
  margin: 0;
  color: white;
  font-size: 18px;
}

.period-header .time {
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
}

.period-dishes {
  padding: 15px;
  max-height: 400px;
  overflow-y: auto;
}

.dish-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.dish-item:last-child {
  border-bottom: none;
}

.dish-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.dish-info img {
  width: 50px;
  height: 50px;
  border-radius: 8px;
  object-fit: cover;
}

.dish-detail {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.dish-detail .name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.dish-detail .price {
  font-size: 13px;
  color: #f56c6c;
}

.dish-quantity {
  font-size: 13px;
  color: #909399;
}

.period-actions {
  display: flex;
  gap: 10px;
  padding: 15px;
  border-top: 1px solid #ebeef5;
}

.period-actions .el-button {
  flex: 1;
}

.dish-selection {
  margin-top: 20px;
}

.dish-selection h4 {
  margin: 0 0 15px 0;
  font-size: 14px;
  color: #606266;
}

.time-separator {
  margin: 0 10px;
  color: #909399;
}
</style>
