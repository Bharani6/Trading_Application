import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: () => import('../views/Landing.vue')
  },
  {
    path: '/',
    component: () => import('../layouts/AuthLayout.vue'),
    children: [
      { path: 'login', name: 'Login', component: () => import('../views/Login.vue') },
      { path: 'register', name: 'Register', component: () => import('../views/Register.vue') }
    ]
  },
  {
    path: '/',
    component: () => import('../layouts/UserLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('../views/Dashboard.vue') },
      { path: 'wallet', name: 'Wallet', component: () => import('../views/Wallet.vue') },
      { path: 'trade', name: 'Trade', component: () => import('../views/Trade.vue') },
      { path: 'billing', name: 'Billing', component: () => import('../views/Billing.vue') },
      { path: 'profile', name: 'Profile', component: () => import('../views/Profile.vue') },
    ]
  },
  {
    path: '/admin',
    component: () => import('../layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      { path: '', redirect: '/admin/dashboard' },
      { path: 'dashboard', name: 'AdminDashboard', component: () => import('../views/AdminDashboard.vue') },
      { path: 'users', name: 'UserManagement', component: () => import('../views/UserManagement.vue') },
      { path: 'stocks', name: 'AdminStockUpload', component: () => import('../views/AdminStockUpload.vue') }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation Guards
router.beforeEach((to, from, next) => {
  // We will implement auth check here using our reactive store
  next()
})

export default router
