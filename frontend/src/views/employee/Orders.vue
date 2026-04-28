<template>
  <div class="orders-page">
    <el-tabs v-model="activeStatus" @tab-change="loadOrders">
      <el-tab-pane label="全部" name="">
        <template #label>
          <el-icon><Document /></el-icon>
          全部
        </template>
      </el-tab-pane>
      <el-tab-pane label="待制作" name="pending_confirm">
        <template #label>
          <el-tag type="info">待制作</el-tag>
        </template>
      </el-tab-pane>
      <el-tab-pane label="制作中" name="in_production">
        <template #label>
          <el-tag type="warning">制作中</el-tag>
        </template>
      </el-tab-pane>
      <el-tab-pane label="待取餐" name="ready_for_pickup">
        <template #label>
          <el-tag type="success">待取餐</el-tag>
        </template>
      </el-tab-pane>
      <el-tab-pane label="已完成" name="picked_up">
        <template #label>
          <el-tag>已完成</el-tag>
        </template>
      </el-tab-pane>
    </el-tabs>

    <div class="orders-list" v-loading="loading">
      <el-empty v-if="orders.length === 0" description="暂无订单" />
      
      <div class="order-card" v-for="order in orders" :key="order.id" @click="goToDetail(order.id)">
        <div class="order-header">
          <span class="order-no">{{ order.order_no }}</span>
          <el-tag :type="getStatusType(order.status)" size="small">
            {{ getStatusText(order.status) }}
          </el-tag>
        </div>
        <div class="order-items">
          <div class="order-item" v-for="item in order.items?.slice(0, 3)" :key="item.id">
            <span>{{ item.dish_name }} x{{ item.quantity }}</span>
            <span>¥{{ item.subtotal.toFixed(2) }}</span>
          </div>
          <div v-if="order.items && order.items.length > 3" class="more-text">
            还有 {{ order.items.length - 3 }} 件菜品
          </div>
        </div>
        <div class="order-footer">
          <span class="order-time">{{ formatDate(order.created_at) }}</span>
          <span class="order-total">共 {{ order.items?.length || 0 }} 件，合计 ¥{{ order.total_amount.toFixed(2) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { Order } from '@/types'
import { getMyOrders } from '@/api/order'

const router = useRouter()

const loading = ref(false)
const activeStatus = ref('')
const orders = ref<Order[]>([])

const getStatusType = (status: string) => {
  const typeMap: { [key: string]: string } = {
    'pending_confirm': 'info',
    'in_production': 'warning',
    'ready_for_pickup': 'success',
    'picked_up': '',
    'reviewed': ''
  }
  return typeMap[status] || ''
}

const getStatusText = (status: string) => {
  const textMap: { [key: string]: string } = {
    'pending_confirm': '待制作',
    'in_production': '制作中',
    'ready_for_pickup': '待取餐',
    'picked_up': '已取餐',
    'reviewed': '已评价'
  }
  return textMap[status] || status
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await getMyOrders(activeStatus.value || undefined)
    if (res.code === 200 && res.data) {
      orders.value = res.data
    }
  } catch (error) {
    console.error('获取订单失败:', error)
  } finally {
    loading.value = false
  }
}

const goToDetail = (orderId: number) => {
  router.push(`/employee/order/${orderId}`)
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.orders-page {
  min-height: 100%;
}

.orders-list {
  margin-top: 20px;
}

.order-card {
  background: white;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  cursor: pointer;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.order-no {
  font-size: 14px;
  color: #909399;
}

.order-items {
  margin-bottom: 12px;
}

.order-item {
  display: flex;
  justify-content: space-between;
  padding: 5px 0;
  font-size: 14px;
}

.more-text {
  font-size: 12px;
  color: #909399;
  text-align: center;
  padding: 5px 0;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.order-time {
  font-size: 13px;
  color: #909399;
}

.order-total {
  font-size: 14px;
  font-weight: bold;
  color: #333;
}
</style>
