<template>
  <div class="order-detail-page">
    <div class="order-card" v-loading="loading">
      <div class="order-header">
        <div class="order-status">
          <el-tag :type="getStatusType(order.status)" size="large">
            {{ getStatusText(order.status) }}
          </el-tag>
        </div>
        <div class="order-no">
          订单号: {{ order.order_no }}
        </div>
      </div>

      <div class="order-section" v-if="order.status === 'ready_for_pickup'">
        <h3>取餐二维码</h3>
        <div class="qr-code">
          <img :src="qrCodeUrl" alt="取餐二维码" />
        </div>
        <p class="qr-hint">请出示此二维码取餐</p>
      </div>

      <div class="order-section">
        <h3>订单详情</h3>
        <div class="order-items">
          <div class="order-item" v-for="item in order.items" :key="item.id">
            <div class="item-info">
              <span class="item-name">{{ item.dish_name }}</span>
              <span class="item-price">¥{{ item.dish_price.toFixed(2) }}</span>
            </div>
            <div class="item-quantity">
              x{{ item.quantity }}
            </div>
            <div class="item-subtotal">
              ¥{{ item.subtotal.toFixed(2) }}
            </div>
          </div>
        </div>
      </div>

      <div class="order-section">
        <h3>取餐信息</h3>
        <div class="info-row">
          <span class="label">取餐时间:</span>
          <span>{{ formatDateTime(order.pickup_time_start) }} - {{ formatTime(order.pickup_time_end) }}</span>
        </div>
        <div class="info-row">
          <span class="label">下单时间:</span>
          <span>{{ formatDateTime(order.created_at) }}</span>
        </div>
      </div>

      <div class="order-total">
        <span>合计</span>
        <span class="total-price">¥{{ order.total_amount.toFixed(2) }}</span>
      </div>

      <div class="order-actions" v-if="order.status === 'picked_up' || order.status === 'reviewed'">
        <el-button type="primary" @click="showReviewDialog = true" :disabled="order.status === 'reviewed'">
          {{ order.status === 'reviewed' ? '已评价' : '去评价' }}
        </el-button>
      </div>
    </div>

    <el-dialog v-model="showReviewDialog" title="评价订单" width="500px">
      <div v-for="item in order.items" :key="item.id" class="review-item">
        <div class="review-dish-name">{{ item.dish_name }}</div>
        <el-rate
          :model-value="getRating(item.dish_id)"
          :max="5"
          show-score
          text-color="#ff9900"
          score-template="{value} 分"
          @change="onRatingChange(item.dish_id, $event)"
        />
        <el-input
          :model-value="getComment(item.dish_id)"
          type="textarea"
          placeholder="写下你的评价..."
          :rows="2"
          @input="onCommentChange(item.dish_id, $event)"
        />
      </div>
      <template #footer>
        <el-button @click="showReviewDialog = false">取消</el-button>
        <el-button type="primary" @click="submitReview" :loading="submitting">提交评价</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import QRCode from 'qrcode'
import type { Order, OrderItem } from '@/types'
import { getOrderById } from '@/api/order'
import { createReview, type CreateReviewParams } from '@/api/stats'

const route = useRoute()

const loading = ref(false)
const showReviewDialog = ref(false)
const submitting = ref(false)
const order = ref<Order>({
  id: 0,
  order_no: '',
  user_id: 0,
  total_amount: 0,
  status: 'pending_confirm',
  pickup_time_start: '',
  pickup_time_end: '',
  items: [],
  created_at: '',
  updated_at: ''
})

const reviewData = reactive<{ [key: number]: { rating: number; comment: string } }>({})

const qrCodeUrl = computed(() => {
  if (!order.value.qr_code_content) return ''
  try {
    return QRCode.toDataURL(order.value.qr_code_content, { width: 200 })
  } catch {
    return ''
  }
})

const getStatusType = (status: string) => {
  const typeMap: { [key: string]: string } = {
    'pending_confirm': 'info',
    'in_production': 'warning',
    'ready_for_pickup': 'success',
    'picked_up': '',
    'reviewed': 'success'
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

const formatDateTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getFullYear()}/${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

const getRating = (dishId: number) => {
  return reviewData[dishId]?.rating ?? 5
}

const getComment = (dishId: number) => {
  return reviewData[dishId]?.comment ?? ''
}

const onRatingChange = (dishId: number, rating: number) => {
  if (!reviewData[dishId]) {
    reviewData[dishId] = { rating, comment: '' }
  } else {
    reviewData[dishId].rating = rating
  }
}

const onCommentChange = (dishId: number, comment: string) => {
  if (!reviewData[dishId]) {
    reviewData[dishId] = { rating: 5, comment }
  } else {
    reviewData[dishId].comment = comment
  }
}

const submitReview = async () => {
  const reviews: CreateReviewParams[] = []
  
  for (const [dishId, data] of Object.entries(reviewData)) {
    if (data.rating > 0) {
      reviews.push({
        dish_id: Number(dishId),
        rating: data.rating,
        comment: data.comment
      })
    }
  }

  if (reviews.length === 0) {
    ElMessage.warning('请至少对一个菜品进行评分')
    return
  }

  submitting.value = true
  try {
    const res = await createReview(order.value.id, reviews)
    if (res.code === 200) {
      ElMessage.success('评价成功')
      showReviewDialog.value = false
      loadOrder()
    }
  } catch (error) {
    console.error('提交评价失败:', error)
  } finally {
    submitting.value = false
  }
}

const loadOrder = async () => {
  const id = Number(route.params.id)
  if (!id) return

  loading.value = true
  try {
    const res = await getOrderById(id)
    if (res.code === 200 && res.data) {
      order.value = res.data
      
      order.value.items?.forEach(item => {
        reviewData[item.dish_id] = { rating: 5, comment: '' }
      })
    }
  } catch (error) {
    console.error('获取订单详情失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadOrder()
})
</script>

<style scoped>
.order-detail-page {
  max-width: 600px;
  margin: 0 auto;
}

.order-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 15px;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 20px;
}

.order-no {
  font-size: 13px;
  color: #909399;
}

.order-section {
  margin-bottom: 20px;
}

.order-section h3 {
  font-size: 15px;
  color: #333;
  margin-bottom: 12px;
}

.qr-code {
  text-align: center;
  padding: 20px;
}

.qr-code img {
  width: 200px;
  height: 200px;
}

.qr-hint {
  text-align: center;
  color: #909399;
  font-size: 13px;
}

.order-items {
  background: #f9f9f9;
  border-radius: 8px;
  padding: 10px;
}

.order-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 5px;
  border-bottom: 1px solid #eee;
}

.order-item:last-child {
  border-bottom: none;
}

.item-info {
  flex: 1;
}

.item-name {
  display: block;
  font-size: 14px;
  color: #333;
}

.item-price {
  font-size: 12px;
  color: #909399;
}

.item-quantity {
  width: 50px;
  text-align: center;
  color: #606266;
}

.item-subtotal {
  width: 80px;
  text-align: right;
  font-weight: bold;
  color: #333;
}

.info-row {
  display: flex;
  padding: 8px 0;
  font-size: 14px;
}

.info-row .label {
  width: 80px;
  color: #909399;
}

.order-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 15px;
  border-top: 1px solid #f0f0f0;
  font-size: 16px;
  font-weight: bold;
}

.total-price {
  color: #f56c6c;
  font-size: 20px;
}

.order-actions {
  margin-top: 20px;
  text-align: center;
}

.review-item {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.review-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.review-dish-name {
  font-weight: bold;
  margin-bottom: 10px;
}
</style>
