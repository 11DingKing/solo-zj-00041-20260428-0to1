import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/modules/user'

const service: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

service.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    if (res.code !== 200) {
      let errorMsg = res.message || '请求失败'
      
      if (res.code === 401) {
        errorMsg = '登录已过期'
        const userStore = useUserStore()
        userStore.logout()
        window.location.reload()
      } else if (res.code === 403) {
        errorMsg = '没有操作权限'
      } else if (res.code === 500) {
        errorMsg = '服务器内部错误'
      }
      
      ElMessage.error(errorMsg)
      
      return Promise.reject(new Error(errorMsg))
    }
    return res
  },
  (error) => {
    console.error('请求错误:', error)
    let errorMsg = '网络请求失败，请稍后重试'
    
    if (error.response) {
      const status = error.response.status
      if (status === 401) {
        errorMsg = '登录已过期'
        const userStore = useUserStore()
        userStore.logout()
        window.location.reload()
      } else if (status === 403) {
        errorMsg = '没有操作权限'
      } else if (status === 500) {
        errorMsg = '服务器内部错误'
      } else if (status === 404) {
        errorMsg = '请求的资源不存在'
      }
    }
    
    ElMessage.error(errorMsg)
    return Promise.reject(new Error(errorMsg))
  }
)

export default service
