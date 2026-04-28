<template>
  <div class="dishes-page">
    <div class="page-header">
      <h2>菜品管理</h2>
      <el-button type="primary" :icon="Plus" @click="handleCreate">添加菜品</el-button>
    </div>

    <div class="filter-bar">
      <el-select v-model="filterCategory" placeholder="选择分类" clearable @change="loadDishes">
        <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
      </el-select>
      <el-select v-model="filterAvailable" placeholder="上下架状态" clearable @change="loadDishes">
        <el-option label="已上架" :value="true" />
        <el-option label="已下架" :value="false" />
      </el-select>
      <el-button :icon="Refresh" @click="loadDishes">刷新</el-button>
    </div>

    <div class="dishes-grid" v-loading="loading">
      <el-empty v-if="dishes.length === 0" description="暂无菜品" />
      
      <div v-else class="grid-container">
        <div class="dish-card" v-for="dish in dishes" :key="dish.id">
          <div class="dish-image">
            <img :src="dish.image || defaultImage" :alt="dish.name" />
            <div class="dish-status">
              <el-tag :type="dish.is_available ? 'success' : 'info'" size="small">
                {{ dish.is_available ? '上架' : '下架' }}
              </el-tag>
            </div>
          </div>
          <div class="dish-info">
            <h3 class="dish-name">{{ dish.name }}</h3>
            <p class="dish-category">分类：{{ getCategoryName(dish.category_id) }}</p>
            <p class="dish-price">
              <span class="price">¥{{ dish.price.toFixed(2) }}</span>
              <span class="limit">每日限量：{{ dish.daily_limit }} 份</span>
            </p>
            <div class="dish-allergens" v-if="dish.allergens && dish.allergens.length > 0">
              <el-tag v-for="allergen in dish.allergens" :key="allergen" size="mini" type="warning">
                {{ getLabelName('allergens', allergen) }}
              </el-tag>
            </div>
            <p class="dish-desc" v-if="dish.description">{{ dish.description }}</p>
          </div>
          <div class="dish-actions">
            <el-button size="small" @click="handleEdit(dish)">编辑</el-button>
            <el-button 
              size="small" 
              :type="dish.is_available ? 'warning' : 'success'"
              @click="toggleAvailable(dish)"
            >
              {{ dish.is_available ? '下架' : '上架' }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(dish)">删除</el-button>
          </div>
        </div>
      </div>
    </div>

    <el-dialog 
      v-model="showDialog" 
      :title="isEdit ? '编辑菜品' : '添加菜品'" 
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form 
        ref="formRef" 
        :model="formData" 
        :rules="formRules" 
        label-width="100px"
      >
        <el-form-item label="菜品名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入菜品名称" />
        </el-form-item>
        <el-form-item label="分类" prop="category_id">
          <el-select v-model="formData.category_id" placeholder="请选择分类" style="width: 100%">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number 
            v-model="formData.price" 
            :min="0" 
            :precision="2" 
            placeholder="请输入价格"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="每日限量" prop="daily_limit">
          <el-input-number 
            v-model="formData.daily_limit" 
            :min="0" 
            placeholder="请输入每日限量"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="图片URL" prop="image">
          <el-input v-model="formData.image" placeholder="请输入图片URL（可选）" />
        </el-form-item>
        <el-form-item label="过敏原标签" prop="allergens">
          <el-checkbox-group v-model="formData.allergens">
            <el-checkbox 
              v-for="allergen in allergenOptions" 
              :key="allergen.value" 
              :label="allergen.value"
            >
              {{ allergen.label }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="formData.description" 
            type="textarea" 
            :rows="3" 
            placeholder="请输入菜品描述（可选）"
          />
        </el-form-item>
        <el-form-item label="是否上架" prop="is_available">
          <el-switch v-model="formData.is_available" active-text="上架" inactive-text="下架" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { Dish, Category } from '@/types'
import { getCategories, getDishes, createDish, updateDish, deleteDish } from '@/api/dish'

const loading = ref(false)
const submitting = ref(false)
const dishes = ref<Dish[]>([])
const categories = ref<Category[]>([])
const filterCategory = ref<number | null>(null)
const filterAvailable = ref<boolean | null>(null)

const showDialog = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const editId = ref<number | null>(null)

const defaultImage = 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=food%20dish%20placeholder%20image&image_size=square'

const allergenOptions = [
  { label: '花生', value: 'peanut' },
  { label: '海鲜', value: 'seafood' },
  { label: '乳制品', value: 'dairy' },
  { label: '麸质', value: 'gluten' },
  { label: '坚果', value: 'nut' },
  { label: '蛋类', value: 'egg' },
  { label: '大豆', value: 'soy' }
]

const formData = reactive({
  name: '',
  category_id: null as number | null,
  price: 0,
  daily_limit: 100,
  image: '',
  description: '',
  allergens: [] as string[],
  is_available: true
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入菜品名称', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }]
}

const getCategoryName = (categoryId: number) => {
  const cat = categories.value.find(c => c.id === categoryId)
  return cat?.name || '未知分类'
}

const getLabelName = (type: string, value: string) => {
  if (type === 'allergens') {
    const option = allergenOptions.find(o => o.value === value)
    return option?.label || value
  }
  return value
}

const loadCategories = async () => {
  try {
    const res = await getCategories()
    if (res.code === 200 && res.data) {
      categories.value = res.data
    }
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

const loadDishes = async () => {
  loading.value = true
  try {
    const params: { category_id?: number; is_available?: boolean } = {}
    if (filterCategory.value) {
      params.category_id = filterCategory.value
    }
    if (filterAvailable.value !== null) {
      params.is_available = filterAvailable.value
    }
    const res = await getDishes(params)
    if (res.code === 200 && res.data) {
      dishes.value = res.data
    }
  } catch (error) {
    console.error('获取菜品列表失败:', error)
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  formData.name = ''
  formData.category_id = null
  formData.price = 0
  formData.daily_limit = 100
  formData.image = ''
  formData.description = ''
  formData.allergens = []
  formData.is_available = true
  editId.value = null
  isEdit.value = false
}

const handleCreate = () => {
  resetForm()
  showDialog.value = true
}

const handleEdit = (dish: Dish) => {
  isEdit.value = true
  editId.value = dish.id
  formData.name = dish.name
  formData.category_id = dish.category_id
  formData.price = dish.price
  formData.daily_limit = dish.daily_limit
  formData.image = dish.image || ''
  formData.description = dish.description || ''
  formData.allergens = [...dish.allergens]
  formData.is_available = dish.is_available
  showDialog.value = true
}

const handleDelete = async (dish: Dish) => {
  try {
    await ElMessageBox.confirm(`确定要删除菜品「${dish.name}」吗？`, '确认删除', {
      type: 'warning'
    })
    const res = await deleteDish(dish.id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      loadDishes()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

const toggleAvailable = async (dish: Dish) => {
  try {
    const res = await updateDish(dish.id, { is_available: !dish.is_available })
    if (res.code === 200) {
      ElMessage.success(dish.is_available ? '已下架' : '已上架')
      loadDishes()
    }
  } catch (error) {
    console.error('更新状态失败:', error)
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      const params = {
        name: formData.name,
        category_id: formData.category_id!,
        price: formData.price,
        daily_limit: formData.daily_limit,
        image: formData.image || undefined,
        description: formData.description || undefined,
        allergens: formData.allergens,
        is_available: formData.is_available
      }

      let res
      if (isEdit.value && editId.value) {
        res = await updateDish(editId.value, params)
      } else {
        res = await createDish(params)
      }

      if (res.code === 200) {
        ElMessage.success(isEdit.value ? '编辑成功' : '添加成功')
        showDialog.value = false
        loadDishes()
      } else {
        ElMessage.error(res.message || '操作失败')
      }
    } catch (error) {
      console.error('提交失败:', error)
      ElMessage.error('操作失败')
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  loadCategories()
  loadDishes()
})
</script>

<style scoped>
.dishes-page {
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

.filter-bar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.dishes-grid {
  background: white;
  border-radius: 12px;
  padding: 20px;
  min-height: 400px;
}

.grid-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.dish-card {
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s;
}

.dish-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.dish-image {
  position: relative;
  height: 160px;
  background: #f5f5f5;
}

.dish-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.dish-status {
  position: absolute;
  top: 10px;
  right: 10px;
}

.dish-info {
  padding: 15px;
}

.dish-name {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #333;
  font-weight: 500;
}

.dish-category {
  margin: 0 0 8px 0;
  font-size: 13px;
  color: #909399;
}

.dish-price {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.dish-price .price {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}

.dish-price .limit {
  font-size: 12px;
  color: #909399;
}

.dish-allergens {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
  margin-bottom: 8px;
}

.dish-desc {
  margin: 0;
  font-size: 13px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.dish-actions {
  display: flex;
  gap: 8px;
  padding: 10px 15px 15px;
  border-top: 1px solid #ebeef5;
}

.dish-actions .el-button {
  flex: 1;
}
</style>
