<template>
  <div class="dashboard-page">
    <div class="page-header">
      <h2>数据看板</h2>
      <el-button type="primary" :icon="Refresh" @click="loadStats" :loading="loading">刷新数据</el-button>
    </div>

    <div class="stats-cards" v-loading="loading">
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="stat-card today-orders">
            <div class="card-icon">
              <el-icon><Document /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">{{ stats?.today_orders || 0 }}</div>
              <div class="card-label">今日订单数</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card today-revenue">
            <div class="card-icon">
              <el-icon><Money /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">¥{{ stats?.today_revenue?.toFixed(2) || '0.00' }}</div>
              <div class="card-label">今日营收</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card top-dishes">
            <div class="card-icon">
              <el-icon><Food /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">{{ stats?.top_dishes?.length || 0 }}</div>
              <div class="card-label">热门菜品</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card ingredients">
            <div class="card-icon">
              <el-icon><Box /></el-icon>
            </div>
            <div class="card-info">
              <div class="card-value">{{ stats?.tomorrow_ingredients?.length || 0 }}</div>
              <div class="card-label">明日食材需求</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <el-row :gutter="20" class="charts-row" v-loading="loading">
      <el-col :span="12">
        <div class="chart-card">
          <h3 class="chart-title">时段订单分布</h3>
          <div ref="periodChartRef" class="chart-container"></div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="chart-card">
          <h3 class="chart-title">近30天营收趋势</h3>
          <div ref="revenueChartRef" class="chart-container"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="tables-row" v-loading="loading">
      <el-col :span="12">
        <div class="table-card">
          <h3 class="table-title">菜品销量排行</h3>
          <el-table :data="stats?.top_dishes || []" stripe>
            <el-table-column prop="dish_name" label="菜品名称" />
            <el-table-column prop="total_orders" label="销量" width="100" />
            <el-table-column label="销售额" width="120">
              <template #default="{ row }">
                ¥{{ row.total_revenue?.toFixed(2) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="table-card">
          <h3 class="table-title">菜品评分排行</h3>
          <el-table :data="stats?.top_rated_dishes || []" stripe>
            <el-table-column prop="dish_name" label="菜品名称" />
            <el-table-column label="评分" width="120">
              <template #default="{ row }">
                <el-rate v-model="row.avg_rating" disabled show-text :texts="['1', '2', '3', '4', '5']" />
              </template>
            </el-table-column>
            <el-table-column prop="review_count" label="评价数" width="100" />
          </el-table>
        </div>
      </el-col>
    </el-row>

    <div class="table-card" v-loading="loading">
      <h3 class="table-title">明日食材消耗预估</h3>
      <el-table :data="stats?.tomorrow_ingredients || []" stripe>
        <el-table-column prop="dish_id" label="菜品ID" width="100" />
        <el-table-column prop="dish_name" label="菜品名称" />
        <el-table-column prop="quantity" label="预估需求量" width="150">
          <template #default="{ row }">
            <el-tag type="warning">{{ row.quantity }} 份</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import * as echarts from 'echarts'
import type { DashboardStats } from '@/types'
import { getDashboardStats } from '@/api/stats'

const loading = ref(false)
const stats = ref<DashboardStats | null>(null)
const periodChartRef = ref<HTMLElement>()
const revenueChartRef = ref<HTMLElement>()

let periodChart: echarts.ECharts | null = null
let revenueChart: echarts.ECharts | null = null

const loadStats = async () => {
  loading.value = true
  try {
    const res = await getDashboardStats()
    if (res.code === 200 && res.data) {
      stats.value = res.data
      await nextTick()
      renderCharts()
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

const renderCharts = () => {
  renderPeriodChart()
  renderRevenueChart()
}

const renderPeriodChart = () => {
  if (!periodChartRef.value || !stats.value) return

  if (!periodChart) {
    periodChart = echarts.init(periodChartRef.value)
  }

  const data = stats.value.meal_period_distribution || []
  const names: { [key: string]: string } = {
    breakfast: '早餐',
    lunch: '午餐',
    dinner: '晚餐'
  }

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} 笔 ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: '10%',
      top: 'center'
    },
    series: [
      {
        name: '时段分布',
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['35%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: data.map(item => ({
          value: item.count,
          name: names[item.meal_period] || item.meal_period
        }))
      }
    ],
    color: ['#67C23A', '#409EFF', '#F56C6C']
  }

  periodChart.setOption(option)
}

const renderRevenueChart = () => {
  if (!revenueChartRef.value || !stats.value) return

  if (!revenueChart) {
    revenueChart = echarts.init(revenueChartRef.value)
  }

  const data = stats.value.last_30_days_revenue || []

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const date = params[0].axisValue
        const value = params[0].data
        return `${date}<br/>营收: ¥${value?.toFixed(2) || '0.00'}`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.map(item => {
        const parts = item.date.split('-')
        return `${parseInt(parts[1])}/${parseInt(parts[2])}`
      }),
      axisLabel: {
        rotate: 45,
        interval: Math.floor(data.length / 6)
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '¥{value}'
      }
    },
    series: [
      {
        name: '营收',
        type: 'line',
        smooth: true,
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
            { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
          ])
        },
        lineStyle: {
          color: '#409EFF',
          width: 2
        },
        itemStyle: {
          color: '#409EFF'
        },
        data: data.map(item => item.revenue)
      }
    ]
  }

  revenueChart.setOption(option)
}

const handleResize = () => {
  periodChart?.resize()
  revenueChart?.resize()
}

onMounted(() => {
  loadStats()
  window.addEventListener('resize', handleResize)
})

watch(() => stats.value, () => {
  nextTick(() => renderCharts())
})
</script>

<style scoped>
.dashboard-page {
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

.stats-cards {
  margin-bottom: 20px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 15px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.card-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-icon .el-icon {
  font-size: 28px;
  color: white;
}

.today-orders .card-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.today-revenue .card-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.top-dishes .card-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.ingredients .card-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.card-info {
  flex: 1;
}

.card-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.card-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.charts-row {
  margin-bottom: 20px;
}

.chart-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.chart-title {
  margin: 0 0 15px 0;
  font-size: 16px;
  color: #333;
  font-weight: 500;
}

.chart-container {
  height: 300px;
  width: 100%;
}

.tables-row {
  margin-bottom: 20px;
}

.table-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.table-title {
  margin: 0 0 15px 0;
  font-size: 16px;
  color: #333;
  font-weight: 500;
}
</style>
