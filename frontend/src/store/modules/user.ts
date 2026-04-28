import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { login, getCurrentUser } from '@/api/auth'
import type { LoginParams } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => userInfo.value?.role || '')
  const isAdmin = computed(() => userInfo.value?.role === 'admin')
  const isChef = computed(() => userInfo.value?.role === 'chef')
  const isEmployee = computed(() => userInfo.value?.role === 'employee')

  async function handleLogin(params: LoginParams) {
    const res = await login(params)
    if (res.code === 200 && res.data) {
      token.value = res.data.token
      userInfo.value = res.data.user
      localStorage.setItem('token', res.data.token)
      return true
    }
    return false
  }

  async function fetchUserInfo() {
    if (!token.value) return false
    
    try {
      const res = await getCurrentUser()
      if (res.code === 200 && res.data) {
        userInfo.value = res.data
        return true
      }
    } catch (error) {
      logout()
    }
    return false
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    userRole,
    isAdmin,
    isChef,
    isEmployee,
    handleLogin,
    fetchUserInfo,
    logout
  }
})
