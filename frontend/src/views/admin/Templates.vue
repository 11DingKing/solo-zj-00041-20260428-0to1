<template>
  <div class="templates-page">
    <div class="page-header">
      <h2>周菜单模板</h2>
      <el-button type="primary" :icon="Plus" @click="handleCreateTemplate">新建模板</el-button>
    </div>

    <div class="templates-container" v-loading="loading">
      <el-empty v-if="templates.length === 0" description="暂无模板，请点击「新建模板」创建" />
      
      <div v-else class="template-cards">
        <div class="template-card" v-for="template in templates" :key="template.id">
          <div class="template-header">
            <div class="template-info">
              <h3>{{ template.name }}</h3>
              <p v-if="template.description">{{ template.description }}</p>
            </div>
            <el-tag :type="template.is_active ? 'success' : 'info'">
              {{ template.is_active ? '启用' : '禁用' }}
            </el-tag>
          </div>
          
          <div class="template-preview" v-if="template.items && template.items.length > 0">
            <el-table :data="getTemplatePreview(template.items)" size="small" stripe>
              <el-table-column prop="day" label="星期" width="100" />
              <el-table-column prop="breakfast" label="早餐">
                <template #default="{ row }">
                  <template v-if="row.breakfast.length > 0">
                    <el-tag v-for="(dish, idx) in row.breakfast" :key="idx" size="mini" style="margin: 2px;">
                      {{ dish }}
                    </el-tag>
                  </template>
                  <span v-else class="empty-text">无</span>
                </template>
              </el-table-column>
              <el-table-column prop="lunch" label="午餐">
                <template #default="{ row }">
                  <template v-if="row.lunch.length > 0">
                    <el-tag v-for="(dish, idx) in row.lunch" :key="idx" size="mini" style="margin: 2px;">
                      {{ dish }}
                    </el-tag>
                  </template>
                  <span v-else class="empty-text">无</span>
                </template>
              </el-table-column>
              <el-table-column prop="dinner" label="晚餐">
                <template #default="{ row }">
                  <template v-if="row.dinner.length > 0">
                    <el-tag v-for="(dish, idx) in row.dinner" :key="idx" size="mini" style="margin: 2px;">
                      {{ dish }}
                    </el-tag>
                  </template>
                  <span v-else class="empty-text">无</span>
                </template>
              </el-table-column>
            </el-table>
          </div>
          <el-empty v-else description="暂无菜品配置" :image-size="60" />
          
          <div class="template-actions">
            <el-button size="small" @click="handleEditTemplate(template)">编辑</el-button>
            <el-button size="small" type="success" @click="handleApplyTemplate(template)">应用</el-button>
            <el-button size="small" type="danger" @click="handleDeleteTemplate(template)">删除</el-button>
          </div>
        </div>
      </div>
    </div>

    <el-dialog 
      v-model="showTemplateDialog" 
      :title="isEditTemplate ? '编辑模板' : '新建模板'" 
      width="800px"
      :close-on-click-modal="false"
    >
      <el-form 
        ref="templateFormRef" 
        :model="templateForm" 
        :rules="templateRules" 
        label-width="100px"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="templateForm.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="templateForm.description" 
            type="textarea" 
            :rows="2" 
            placeholder="请输入描述（可选）"
          />
        </el-form-item>
      </el-form>

      <div class="template-week-edit">
        <el-tabs v-model="activeDay" type="card">
          <el-tab-pane 
            v-for="day in weekDays" 
            :key="day.value" 
            :label="day.label" 
            :name="String(day.value)"
          >
            <div class="day-meal-section">
              <h4>早餐</h4>
              <el-select
                v-model="weekConfig[activeDay].breakfast"
                multiple
                filterable
                placeholder="选择早餐菜品"
                style="width: 100%"
              >
                <el-option
                  v-for="dish in availableDishes"
                  :key="dish.id"
                  :label="`${dish.name} (¥${dish.price})`"
                  :value="dish.id"
                />
              </el-select>
            </div>
            <div class="day-meal-section">
              <h4>午餐</h4>
              <el-select
                v-model="weekConfig[activeDay].lunch"
                multiple
                filterable
                placeholder="选择午餐菜品"
                style="width: 100%"
              >
                <el-option
                  v-for="dish in availableDishes"
                  :key="dish.id"
                  :label="`${dish.name} (¥${dish.price})`"
                  :value="dish.id"
                />
              </el-select>
            </div>
            <div class="day-meal-section">
              <h4>晚餐</h4>
              <el-select
                v-model="weekConfig[activeDay].dinner"
                multiple
                filterable
                placeholder="选择晚餐菜品"
                style="width: 100%"
              >
                <el-option
                  v-for="dish in availableDishes"
                  :key="dish.id"
                  :label="`${dish.name} (¥${dish.price})`"
                  :value="dish.id"
                />
              </el-select>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
      <template #footer>
        <el-button @click="showTemplateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSaveTemplate" :loading="savingTemplate">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog 
      v-model="showApplyDialog" 
      title="应用模板" 
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form label-width="100px">
        <el-form-item label="起始日期">
          <el-date-picker
            v-model="applyStartDate"
            type="date"
            placeholder="选择一周起始日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <el-alert 
        title="提示" 
        type="info" 
        :closable="false"
        show-icon
      >
        应用模板将根据所选起始日期，为接下来7天创建每日菜单。
      </el-alert>
      <template #footer>
        <el-button @click="showApplyDialog = false">取消</el-button>
        <el-button type="primary" @click="handleConfirmApply" :loading="applyingTemplate">确认应用</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { WeeklyMenuTemplate, WeeklyMenuTemplateItem, Dish } from '@/types'
import { getTemplates, createTemplate, updateTemplateItems, applyTemplate } from '@/api/template'
import { getDishes } from '@/api/dish'

const loading = ref(false)
const templates = ref<WeeklyMenuTemplate[]>([])
const availableDishes = ref<Dish[]>([])

const weekDays = [
  { label: '周一', value: 1 },
  { label: '周二', value: 2 },
  { label: '周三', value: 3 },
  { label: '周四', value: 4 },
  { label: '周五', value: 5 },
  { label: '周六', value: 6 },
  { label: '周日', value: 7 }
]

const showTemplateDialog = ref(false)
const isEditTemplate = ref(false)
const savingTemplate = ref(false)
const editTemplateId = ref<number | null>(null)
const activeDay = ref('1')
const templateFormRef = ref<FormInstance>()

const showApplyDialog = ref(false)
const applyingTemplate = ref(false)
const applyTemplateId = ref<number | null>(null)
const applyStartDate = ref('')

const templateForm = reactive({
  name: '',
  description: ''
})

const templateRules: FormRules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }]
}

interface WeekDayConfig {
  breakfast: number[]
  lunch: number[]
  dinner: number[]
}

const weekConfig = reactive<{ [key: string]: WeekDayConfig }>({
  '1': { breakfast: [], lunch: [], dinner: [] },
  '2': { breakfast: [], lunch: [], dinner: [] },
  '3': { breakfast: [], lunch: [], dinner: [] },
  '4': { breakfast: [], lunch: [], dinner: [] },
  '5': { breakfast: [], lunch: [], dinner: [] },
  '6': { breakfast: [], lunch: [], dinner: [] },
  '7': { breakfast: [], lunch: [], dinner: [] }
})

const resetWeekConfig = () => {
  for (let i = 1; i <= 7; i++) {
    weekConfig[String(i)] = { breakfast: [], lunch: [], dinner: [] }
  }
}

const getTemplatePreview = (items: WeeklyMenuTemplateItem[]) => {
  const preview: { [key: number]: any } = {}
  
  for (let i = 1; i <= 7; i++) {
    preview[i] = {
      day: weekDays[i - 1].label,
      breakfast: [] as string[],
      lunch: [] as string[],
      dinner: [] as string[]
    }
  }

  items.forEach(item => {
    const day = item.day_of_week
    const period = item.meal_period as 'breakfast' | 'lunch' | 'dinner'
    const dishName = item.dish?.name || `菜品${item.dish_id}`
    
    if (preview[day]) {
      preview[day][period].push(dishName)
    }
  })

  return Object.values(preview)
}

const loadTemplates = async () => {
  loading.value = true
  try {
    const res = await getTemplates()
    if (res.code === 200 && res.data) {
      templates.value = res.data
    }
  } catch (error) {
    console.error('获取模板列表失败:', error)
  } finally {
    loading.value = false
  }
}

const loadDishes = async () => {
  try {
    const res = await getDishes({ is_available: true })
    if (res.code === 200 && res.data) {
      availableDishes.value = res.data
    }
  } catch (error) {
    console.error('获取菜品列表失败:', error)
  }
}

const handleCreateTemplate = () => {
  isEditTemplate.value = false
  editTemplateId.value = null
  templateForm.name = ''
  templateForm.description = ''
  resetWeekConfig()
  activeDay.value = '1'
  showTemplateDialog.value = true
}

const handleEditTemplate = (template: WeeklyMenuTemplate) => {
  isEditTemplate.value = true
  editTemplateId.value = template.id
  templateForm.name = template.name
  templateForm.description = template.description || ''
  
  resetWeekConfig()
  if (template.items) {
    template.items.forEach(item => {
      const dayKey = String(item.day_of_week)
      const period = item.meal_period as 'breakfast' | 'lunch' | 'dinner'
      if (weekConfig[dayKey]) {
        weekConfig[dayKey][period].push(item.dish_id)
      }
    })
  }
  
  activeDay.value = '1'
  showTemplateDialog.value = true
}

const handleDeleteTemplate = async (template: WeeklyMenuTemplate) => {
  try {
    await ElMessageBox.confirm(`确定要删除模板「${template.name}」吗？`, '确认删除', {
      type: 'warning'
    })
    // TODO: 实现删除模板API
    ElMessage.info('删除功能开发中')
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

const handleSaveTemplate = async () => {
  if (!templateFormRef.value) return

  await templateFormRef.value.validate(async (valid) => {
    if (!valid) return

    savingTemplate.value = true
    try {
      const items: { day_of_week: number; meal_period: string; dish_ids: number[] }[] = []
      
      for (let i = 1; i <= 7; i++) {
        const dayKey = String(i)
        const config = weekConfig[dayKey]
        
        if (config.breakfast.length > 0) {
          items.push({ day_of_week: i, meal_period: 'breakfast', dish_ids: config.breakfast })
        }
        if (config.lunch.length > 0) {
          items.push({ day_of_week: i, meal_period: 'lunch', dish_ids: config.lunch })
        }
        if (config.dinner.length > 0) {
          items.push({ day_of_week: i, meal_period: 'dinner', dish_ids: config.dinner })
        }
      }

      if (isEditTemplate.value && editTemplateId.value) {
        const res = await updateTemplateItems(editTemplateId.value, items)
        if (res.code === 200) {
          ElMessage.success('更新成功')
          showTemplateDialog.value = false
          loadTemplates()
        } else {
          ElMessage.error(res.message || '更新失败')
        }
      } else {
        const createRes = await createTemplate({
          name: templateForm.name,
          description: templateForm.description || undefined
        })
        
        if (createRes.code === 200 && createRes.data) {
          if (items.length > 0) {
            await updateTemplateItems(createRes.data.id, items)
          }
          ElMessage.success('创建成功')
          showTemplateDialog.value = false
          loadTemplates()
        } else {
          ElMessage.error(createRes.message || '创建失败')
        }
      }
    } catch (error) {
      console.error('保存失败:', error)
      ElMessage.error('操作失败')
    } finally {
      savingTemplate.value = false
    }
  })
}

const handleApplyTemplate = (template: WeeklyMenuTemplate) => {
  applyTemplateId.value = template.id
  applyStartDate.value = new Date().toISOString().split('T')[0]
  showApplyDialog.value = true
}

const handleConfirmApply = async () => {
  if (!applyTemplateId.value || !applyStartDate.value) {
    ElMessage.warning('请选择起始日期')
    return
  }

  applyingTemplate.value = true
  try {
    const res = await applyTemplate(applyTemplateId.value, applyStartDate.value)
    if (res.code === 200) {
      ElMessage.success('应用成功')
      showApplyDialog.value = false
    } else {
      ElMessage.error(res.message || '应用失败')
    }
  } catch (error) {
    console.error('应用模板失败:', error)
    ElMessage.error('应用失败')
  } finally {
    applyingTemplate.value = false
  }
}

onMounted(() => {
  loadTemplates()
  loadDishes()
})
</script>

<style scoped>
.templates-page {
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

.templates-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  min-height: 400px;
}

.template-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(500px, 1fr));
  gap: 20px;
}

.template-card {
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  overflow: hidden;
}

.template-header {
  padding: 15px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.template-info h3 {
  margin: 0;
  color: white;
  font-size: 16px;
}

.template-info p {
  margin: 5px 0 0 0;
  color: rgba(255, 255, 255, 0.8);
  font-size: 13px;
}

.template-preview {
  padding: 15px;
}

.empty-text {
  color: #c0c4cc;
  font-size: 12px;
}

.template-actions {
  display: flex;
  gap: 10px;
  padding: 15px;
  border-top: 1px solid #ebeef5;
}

.template-actions .el-button {
  flex: 1;
}

.template-week-edit {
  margin-top: 20px;
}

.day-meal-section {
  margin-bottom: 20px;
}

.day-meal-section h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #606266;
}
</style>
