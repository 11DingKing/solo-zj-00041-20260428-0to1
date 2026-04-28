<template>
  <div class="scan-page">
    <div class="page-header">
      <h2>扫码取餐</h2>
      <p class="desc">扫描订单二维码或手动输入订单号确认取餐</p>
    </div>

    <div class="scan-container">
      <div class="scan-box">
        <div class="scan-area">
          <el-icon class="scan-icon"><Camera /></el-icon>
          <p class="scan-hint">请对准订单二维码</p>
          <p class="scan-note">（注：浏览器环境暂不支持摄像头扫码，<br>请使用下方手动输入方式）</p>
        </div>
      </div>

      <div class="manual-input-box">
        <h3>手动输入订单号</h3>
        <el-form label-position="top">
          <el-form-item label="订单号">
            <el-input 
              v-model="orderNo" 
              placeholder="请输入订单号，例如：ORD20260428000001"
              clearable
              @keyup.enter="searchOrder"
            >
              <template #append>
                <el-button @click="searchOrder" :loading="searching">查询</el-button>
              </template>
            </el-input>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <el-dialog v-model="showOrderDetail" title="订单确认" width="450px">
      <div v-if="searchResult" class="order-confirm">
        <el-alert 
          v-if="searchResult.status === 'ready_for_pickup'" 
          type="success" 
          :closable="false"
          show-icon
        >
          该订单已准备完成，可以取餐
        </el-alert>
        <el-alert 
          v-else-if="searchResult.status === 'picked_up' || searchResult.status === 'reviewed'"
          type="warning" 
          :closable="false"
          show-icon
        >
          该订单已取餐
        </el-alert>
        <el-alert 
          v-else
          type="info" 
          :closable="false"
          show-icon
        >
          订单状态：{{ getStatusText(searchResult.status) }}
        </el-alert>

        <div class="order-info">
          <div class="info-row">
            <span class="label">订单号：</span>
            <span class="value">{{ searchResult.order_no }}</span>
          </div>
          <div class="info-row">
            <span class="label">用户：</span>
            <span class="value">{{ searchResult.user?.name }}</span>
          </div>
          <div class="info-row">
            <span class="label">取餐时间：</span>
            <span class="value">{{ formatTime(searchResult.pickup_time_start) }}</span>
          </div>
        </div>

        <div class="order-dishes">
          <h4>菜品列表</h4>
          <div v-for="item in searchResult.items" :key="item.id" class="dish-item">
            <span>{{ item.dish_name }}</span>
            <span>x{{ item.quantity }}</span>
          </div>
        </div>

        <div class="order-total">
          <span>订单金额：</span>
          <span class="price">¥{{ searchResult.total_amount.toFixed(2) }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="showOrderDetail = false">取消</el-button>
        <el-button 
          v-if="searchResult?.status === 'ready_for_pickup'"
          type="primary" 
          @click="handleConfirmPickup"
          :loading="confirming"
        >
          确认取餐
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { Order } from '@/types'
import { getOrderByNo, confirmPickup } from '@/api/order'

const orderNo = ref('')
const searching = ref(false)
const confirming = ref(false)
const showOrderDetail = ref(false)
const searchResult = ref<Order | null>(null)

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
  return `${date.getFullYear()}/${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

const searchOrder = async () => {
  if (!orderNo.value.trim()) {
    ElMessage.warning('请输入订单号')
    return
  }

  searching.value = true
  try {
    const res = await getOrderByNo(orderNo.value.trim())
    if (res.code === 200 && res.data) {
      searchResult.value = res.data
      showOrderDetail.value = true
    } else {
      ElMessage.error(res.message || '未找到该订单')
    }
  } catch (error) {
    ElMessage.error('查询订单失败')
  } finally {
    searching.value = false
  }
}

const handleConfirmPickup = async () => {
  if (!searchResult.value) return

  confirming.value = true
  try {
    const res = await confirmPickup(searchResult.value.order_no)
    if (res.code === 200) {
      ElMessage.success('取餐确认成功')
      showOrderDetail.value = false
      orderNo.value = ''
    } else {
      ElMessage.error(res.message || '确认失败')
    }
  } catch (error) {
    ElMessage.error('确认取餐失败')
  } finally {
    confirming.value = false
  }
}
</script>

<style scoped>
.scan-page {
  min-height: 100%;
}

.page-header {
  margin-bottom: 30px;
}

.page-header h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
  color: #333;
}

.page-header .desc {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.scan-container {
  display: flex;
  gap: 30px;
  flex-wrap: wrap;
}

.scan-box {
  flex: 1;
  min-width: 300px;
}

.scan-area {
  background: #f9f9f9;
  border: 2px dashed #dcdfe6;
  border-radius: 12px;
  padding: 60px 20px;
  text-align: center;
}

.scan-icon {
  font-size: 64px;
  color: #c0c4cc;
}

.scan-hint {
  margin: 15px 0 5px 0;
  font-size: 16px;
  color: #606266;
}

.scan-note {
  margin: 0;
  font-size: 12px;
  color: #909399;
  line-height: 1.8;
}

.manual-input-box {
  flex: 1;
  min-width: 300px;
  background: white;
  border-radius: 12px;
  padding: 25px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.manual-input-box h3 {
  margin: 0 0 20px 0;
  font-size: 16px;
  color: #333;
}

.order-confirm {
  padding: 10px 0;
}

.order-info {
  margin-top: 20px;
  padding: 15px;
  background: #f9f9f9;
  border-radius: 8px;
}

.info-row {
  display: flex;
  padding: 6px 0;
  font-size: 14px;
}

.info-row .label {
  color: #909399;
  width: 80px;
}

.info-row .value {
  color: #333;
  font-weight: 500;
}

.order-dishes {
  margin-top: 15px;
}

.order-dishes h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #606266;
}

.dish-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px dashed #eee;
  font-size: 14px;
}

.order-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #eee;
  font-size: 16px;
  font-weight: bold;
}

.order-total .price {
  color: #f56c6c;
  font-size: 20px;
}
</style>
