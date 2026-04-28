<template>
  <el-container class="layout-container">
    <el-aside width="200px" class="layout-aside">
      <div class="logo">
        <el-icon class="logo-icon"><Food /></el-icon>
        <span class="logo-text">食堂点餐系统</span>
      </div>
      <el-menu :default-active="activeMenu" router background-color="#304156" text-color="#bfcbd9" active-text-color="#409EFF">
        <el-menu-item index="/chef/orders">
          <el-icon><Document /></el-icon>
          <span>订单看板</span>
        </el-menu-item>
        <el-menu-item index="/chef/scan">
          <el-icon><Scan /></el-icon>
          <span>扫码取餐</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="layout-header">
        <div class="header-right">
          <span class="user-name">{{ userStore.userInfo?.name }}</span>
          <el-button type="text" @click="handleLogout">退出</el-button>
        </div>
      </el-header>
      <el-main class="layout-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/modules/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  min-height: 100vh;
}

.layout-aside {
  background: #304156;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: #242f42;
}

.logo-icon {
  font-size: 28px;
  color: #409eff;
}

.logo-text {
  color: white;
  font-size: 16px;
  font-weight: bold;
}

.layout-header {
  background: white;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 20px;
  border-bottom: 1px solid #e6e6e6;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-name {
  color: #606266;
}

.layout-main {
  background: #f0f2f5;
  padding: 20px;
}
</style>
