<template>
  <div class="chef-orders-page">
    <div class="page-header">
      <h2>订单看板</h2>
      <div class="filter-bar">
        <el-date-picker
          v-model="selectedDate"
          type="date"
          placeholder="选择日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          @change="loadOrders"
        />
        <el-radio-group v-model="selectedStatus" @change="loadOrders">
          <el-radio-button label="">全部</el-radio-button>
          <el-radio-button label="pending_confirm">待制作</el-radio-button>
          <el-radio-button label="in_production">制作中</el-radio-button>
          <el-radio-button label="ready_for_pickup">待取餐</el-radio-button>
        </el-radio-group>
        <el-button type="primary" @click="loadOrders" :icon="Refresh">刷新</el-button>
      </div>
    </div>

    <div class="orders-container" v-loading="loading">
      <el-empty v-if="orders.length === 0" description="暂无订单" />
      
      <div class="orders-grid" v-else>
        <div class="order-card" v-for="order in orders" :key="order.id">
          <div class="order-header">
            <span class="order-no">{{ order.order_no }}</span>
            <el-tag :type="getStatusType(order.status)" size="large">
              {{ getStatusText(order.status) }}
            </el-tag>
          </div>
          
          <div class="order-user">
            <el-avatar :size="32">
              <el-icon><User /></el-icon>
            </el-avatar>
            <span class="user-name">{{ order.user?.name }}</span>
            <span class="pickup-time">
              取餐: {{ formatTime(order.pickup_time_start) }}
            </span>
          </div>

          <div class="order-items">
            <div v-for="item in order.items" :key="item.id" class="order-item">
              <span class="item-name">{{ item.dish_name }}</span>
              <span class="item-quantity">x{{ item.quantity }}</span>
            </div>
          </div>

          <div class="order-actions" v-if="order.status === 'pending_confirm'">
            <el-button type="primary" @click="updateStatus(order.id, 'in_production')">
              开始制作
            </el-button>
          </div>
          <div class="order-actions" v-else-if="order.status === 'in_production'">
            <el-button type="success" @click="updateStatus(order.id, 'ready_for_pickup')">
              制作完成
            </el-button>
          </div>
          <div class="order-actions" v-else-if="order.status === 'ready_for_pickup'">
            <el-button type="warning" @click="viewOrder(order)">
              查看详情
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="showOrderDetail" title="订单详情" width="400px">
      <div v-if="currentOrder" class="order-detail">
        <div class="detail-row">
          <span class="label">订单号:</span>
          <span>{{ currentOrder.order_no }}</span>
        </div>
        <div class="detail-row">
          <span class="label">用户:</span>
          <span>{{ currentOrder.user?.name }}</span>
        </div>
        <div class="detail-row">
          <span class="label">状态:</span>
          <el-tag :type="getStatusType(currentOrder.status)">
            {{ getStatusText(currentOrder.status) }}
          </el-tag>
        </div>
        <div class="detail-row">
          <span class="label">取餐时间:</span>
          <span>{{ formatDateTime(currentOrder.pickup_time_start) }}</span>
        </div>
        <div class="detail-section">
          <h4>菜品列表</h4>
          <div v-for="item in currentOrder.items" :key="item.id" class="detail-item">
            <span>{{ item.dish_name }}</span>
            <span>x{{ item.quantity }} - ¥{{ item.subtotal.toFixed(2) }}</span>
          </div>
        </div>
        <div class="detail-total">
          <span>合计:</span>
          <span class="total-price">¥{{ currentOrder.total_amount.toFixed(2) }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="showOrderDetail = false">关闭</el-button>
        <el-button 
          v-if="currentOrder?.status === 'ready_for_pickup'" 
          type="primary" 
          @click="confirmPickup"
        >
          确认取餐
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { Order } from '@/types'
import { getChefOrders, updateOrderStatus, confirmPickup as apiConfirmPickup } from '@/api/order'

const loading = ref(false)
const selectedDate = ref<string>('2026-04-28')
const selectedStatus = ref('')
const orders = ref<Order[]>([])
const showOrderDetail = ref(false)
const currentOrder = ref<Order | null>(null)

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

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

const formatDateTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getFullYear()}/${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

const loadOrders = async () => {
  loading.value = true
  try {
    const params: { status?: string; menu_date?: string } = {}
    if (selectedStatus.value) {
      params.status = selectedStatus.value
    }
    if (selectedDate.value) {
      params.menu_date = selectedDate.value
    }
    const res = await getChefOrders(params)
    if (res.code === 200 && res.data) {
      orders.value = res.data
    }
  } catch (error) {
    console.error('获取订单失败:', error)
  } finally {
    loading.value = false
  }
}

const updateStatus = async (orderId: number, status: string) => {
  try {
    const res = await updateOrderStatus(orderId, status)
    if (res.code === 200) {
      ElMessage.success('状态更新成功')
      loadOrders()
    }
  } catch (error) {
    console.error('更新状态失败:', error)
  }
}

const viewOrder = (order: Order) => {
  currentOrder.value = order
  showOrderDetail.value = true
}

const confirmPickup = async () => {
  if (!currentOrder.value) return
  try {
    const res = await apiConfirmPickup(currentOrder.value.order_no)
    if (res.code === 200) {
      ElMessage.success('取餐确认成功')
      showOrderDetail.value = false
      loadOrders()
    }
  } catch (error) {
    console.error('确认取餐失败:', error)
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.chef-orders-page {
  min-height: 100%;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0 0 15px 0;
  font-size: 20px;
  color: #333;
}

.filter-bar {
  display: flex;
  gap: 15px;
  align-items: center;
  flex-wrap: wrap;
}

.orders-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  min-height: 400px;
}

.orders-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.order-card {
  background: #f9f9f9;
  border-radius: 10px;
  padding: 15px;
  border: 1px solid #eee;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.order-no {
  font-size: 13px;
  color: #909399;
}

.order-user {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  padding-bottom: 10px;
  border-bottom: 1px dashed #eee;
}

.user-name {
  font-weight: bold;
  color: #333;
}

.pickup-time {
  margin-left: auto;
  font-size: 13px;
  color: #909399;
}

.order-items {
  margin-bottom: 15px;
}

.order-item {
  display: flex;
  justify-content: space-between;
  padding: 5px 0;
  font-size: 14px;
}

.item-name {
  color: #333;
}

.item-quantity {
  color: #606266;
}

.order-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.order-detail .detail-row {
  display: flex;
  padding: 8px 0;
  font-size: 14px;
}

.order-detail .detail-row .label {
  width: 80px;
  color: #909399;
}

.order-detail .detail-section {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
}

.order-detail .detail-section h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #333;
}

.order-detail .detail-item {
  display: flex;
  justify-content: space-between;
  padding: 5px 0;
  font-size: 14px;
}

.order-detail .detail-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
  font-size: 16px;
  font-weight: bold;
}

.order-detail .total-price {
  color: #f56c6c;
  font-size: 20px;
}
</style>
