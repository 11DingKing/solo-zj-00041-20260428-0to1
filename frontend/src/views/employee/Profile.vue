<template>
  <div class="profile-page">
    <div class="user-card">
      <el-avatar :size="80" :src="userStore.userInfo?.avatar">
        <el-icon><User /></el-icon>
      </el-avatar>
      <div class="user-info">
        <h3>{{ userStore.userInfo?.name }}</h3>
        <p class="user-role">{{ getRoleText(userStore.userInfo?.role) }}</p>
      </div>
    </div>

    <div class="stats-card" v-if="stats">
      <h3 class="section-title">本月统计</h3>
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-value">¥{{ stats.month_amount.toFixed(2) }}</div>
          <div class="stat-label">消费金额</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ stats.order_count }}</div>
          <div class="stat-label">点餐次数</div>
        </div>
      </div>
      <div class="top-dishes" v-if="stats.top_dishes.length > 0">
        <h4 class="sub-title">最爱菜品 Top5</h4>
        <div class="dish-list">
          <div v-for="(dish, index) in stats.top_dishes" :key="dish.dish_id" class="dish-item">
            <span class="rank">{{ index + 1 }}</span>
            <span class="name">{{ dish.dish_name }}</span>
            <span class="count">{{ dish.count }} 次</span>
          </div>
        </div>
      </div>
    </div>

    <div class="allergen-card">
      <h3 class="section-title">过敏原设置</h3>
      <p class="allergen-hint">设置您的过敏原信息，点餐时会自动提醒您</p>
      <el-checkbox-group v-model="selectedAllergens">
        <el-checkbox label="花生" border>花生</el-checkbox>
        <el-checkbox label="海鲜" border>海鲜</el-checkbox>
        <el-checkbox label="乳制品" border>乳制品</el-checkbox>
        <el-checkbox label="麸质" border>麸质</el-checkbox>
        <el-checkbox label="鸡蛋" border>鸡蛋</el-checkbox>
        <el-checkbox label="肉类" border>肉类</el-checkbox>
        <el-checkbox label="大蒜" border>大蒜</el-checkbox>
      </el-checkbox-group>
      <el-button type="primary" @click="saveAllergens" :loading="saving" style="margin-top: 15px;">
        保存设置
      </el-button>
    </div>

    <div class="action-card">
      <el-button type="danger" plain @click="handleLogout">退出登录</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { getUserStats } from '@/api/stats'
import { updateAllergens } from '@/api/auth'
import type { UserStats } from '@/types'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const saving = ref(false)
const stats = ref<UserStats | null>(null)
const selectedAllergens = ref<string[]>([])

const getRoleText = (role?: string) => {
  const textMap: { [key: string]: string } = {
    'admin': '管理员',
    'chef': '厨师',
    'employee': '员工'
  }
  return textMap[role || ''] || role
}

const loadStats = async () => {
  loading.value = true
  try {
    const res = await getUserStats()
    if (res.code === 200 && res.data) {
      stats.value = res.data
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

const loadAllergens = () => {
  if (userStore.userInfo?.allergens) {
    selectedAllergens.value = userStore.userInfo.allergens.map(a => a.allergen)
  }
}

const saveAllergens = async () => {
  saving.value = true
  try {
    const res = await updateAllergens(selectedAllergens.value)
    if (res.code === 200) {
      ElMessage.success('保存成功')
      if (userStore.userInfo) {
        userStore.userInfo.allergens = selectedAllergens.value.map((a, i) => ({
          id: i + 1,
          user_id: userStore.userInfo!.id,
          allergen: a,
          created_at: new Date().toISOString()
        }))
      }
    }
  } catch (error) {
    console.error('保存过敏原失败:', error)
  } finally {
    saving.value = false
  }
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

onMounted(() => {
  loadStats()
  loadAllergens()
})
</script>

<style scoped>
.profile-page {
  max-width: 600px;
  margin: 0 auto;
}

.user-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 15px;
}

.user-info h3 {
  margin: 0 0 5px 0;
  font-size: 18px;
  color: #333;
}

.user-role {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

.stats-card, .allergen-card, .action-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 15px;
}

.section-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin: 0 0 15px 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.stat-item {
  text-align: center;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.sub-title {
  font-size: 14px;
  color: #606266;
  margin: 0 0 10px 0;
}

.dish-list {
  background: #f9f9f9;
  border-radius: 8px;
  padding: 10px;
}

.dish-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #eee;
}

.dish-item:last-child {
  border-bottom: none;
}

.rank {
  width: 24px;
  height: 24px;
  background: #409eff;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  margin-right: 10px;
}

.name {
  flex: 1;
  font-size: 14px;
  color: #333;
}

.count {
  font-size: 13px;
  color: #909399;
}

.allergen-hint {
  font-size: 13px;
  color: #909399;
  margin-bottom: 15px;
}

.allergen-card .el-checkbox {
  margin-bottom: 10px;
}

.action-card {
  text-align: center;
}
</style>
