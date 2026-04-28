import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/modules/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    redirect: '/employee/menu'
  },
  {
    path: '/employee',
    component: () => import('@/layouts/EmployeeLayout.vue'),
    meta: { requiresAuth: true, role: 'employee' },
    children: [
      {
        path: 'menu',
        name: 'EmployeeMenu',
        component: () => import('@/views/employee/Menu.vue'),
        meta: { title: '今日菜单' }
      },
      {
        path: 'orders',
        name: 'EmployeeOrders',
        component: () => import('@/views/employee/Orders.vue'),
        meta: { title: '我的订单' }
      },
      {
        path: 'order/:id',
        name: 'EmployeeOrderDetail',
        component: () => import('@/views/employee/OrderDetail.vue'),
        meta: { title: '订单详情' }
      },
      {
        path: 'profile',
        name: 'EmployeeProfile',
        component: () => import('@/views/employee/Profile.vue'),
        meta: { title: '个人中心' }
      }
    ]
  },
  {
    path: '/chef',
    component: () => import('@/layouts/ChefLayout.vue'),
    meta: { requiresAuth: true, role: 'chef' },
    children: [
      {
        path: 'orders',
        name: 'ChefOrders',
        component: () => import('@/views/chef/Orders.vue'),
        meta: { title: '订单看板' }
      },
      {
        path: 'scan',
        name: 'ChefScan',
        component: () => import('@/views/chef/Scan.vue'),
        meta: { title: '扫码取餐' }
      }
    ]
  },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, role: 'admin' },
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { title: '数据看板' }
      },
      {
        path: 'dishes',
        name: 'AdminDishes',
        component: () => import('@/views/admin/Dishes.vue'),
        meta: { title: '菜品管理' }
      },
      {
        path: 'menus',
        name: 'AdminMenus',
        component: () => import('@/views/admin/Menus.vue'),
        meta: { title: '菜单管理' }
      },
      {
        path: 'templates',
        name: 'AdminTemplates',
        component: () => import('@/views/admin/Templates.vue'),
        meta: { title: '周菜单模板' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, _from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - 食堂点餐系统` : '食堂点餐系统'

  const userStore = useUserStore()

  if (to.meta.requiresAuth) {
    if (!userStore.isLoggedIn) {
      next('/login')
      return
    }

    if (!userStore.userInfo) {
      await userStore.fetchUserInfo()
    }

    const requiredRole = to.meta.role as string
    if (requiredRole && userStore.userRole !== requiredRole && userStore.userRole !== 'admin') {
      if (userStore.isAdmin) {
        next('/admin/dashboard')
      } else if (userStore.isChef) {
        next('/chef/orders')
      } else {
        next('/employee/menu')
      }
      return
    }
  }

  if (to.path === '/login' && userStore.isLoggedIn) {
    if (userStore.isAdmin) {
      next('/admin/dashboard')
    } else if (userStore.isChef) {
      next('/chef/orders')
    } else {
      next('/employee/menu')
    }
    return
  }

  next()
})

export default router
