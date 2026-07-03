<template>
  <div class="userlayout">
    <!-- Top App Bar -->
    <header class="top-header">
      <div class="header-left">
        <div class="logo-area">
          <i class="fas fa-chart-line logo-icon"></i>
          <h2 class="logo-text">MoneymakePro</h2>
        </div>
        <nav class="top-nav">
          <router-link to="/dashboard" class="nav-item" active-class="active">
            <div class="nav-icon"><i class="fas fa-home"></i></div>
            <span>Home</span>
          </router-link>
          <router-link to="/wallet" class="nav-item" active-class="active">
            <div class="nav-icon"><i class="fas fa-wallet"></i></div>
            <span>Wallet</span>
          </router-link>
          <router-link to="/trade" class="nav-item" active-class="active">
            <div class="nav-icon"><i class="fas fa-chart-line"></i></div>
            <span>Trade</span>
          </router-link>

          <router-link to="/billing" class="nav-item" active-class="active">
            <i class="fas fa-file-invoice-dollar nav-icon"></i>
            <span>Billing</span>
          </router-link>
        </nav>
      </div>

      <div class="header-right">
        <div class="notification-wrapper">
          <button class="icon-btn notification-btn" @click="toggleNotifications">
            <i class="fas fa-bell"></i>
            <span class="badge-pill" v-if="notifications.length">{{ notifications.length }}</span>
          </button>
          <div class="notification-dropdown" v-if="showNotifications">
            <div class="dropdown-header">
              <h4>Notifications</h4>
            </div>
            <div class="dropdown-body">
              <div v-for="n in notifications" :key="n.id" class="notif-item unread">
                <div class="notif-icon"><i :class="n.icon"></i></div>
                <div class="notif-content">
                  <p>{{ n.text }}</p>
                  <span>{{ n.time }}</span>
                </div>
              </div>
              <div v-if="!notifications.length" class="p-20 text-center text-muted">
                No new notifications
              </div>
            </div>
            <div class="dropdown-footer">
              <button @click="markAllAsRead">Mark all as read</button>
            </div>
          </div>
        </div>
        
        <div class="user-profile" @click="router.push('/profile')">
          <div class="avatar-icon">
            <i class="fas fa-user-circle"></i>
          </div>
          <div class="user-info">
            <span class="user-name">{{ authStore.state.user?.name || 'Trader' }}</span>
          </div>
        </div>

        <button v-if="authStore.state.user?.role === 'admin'" @click="$router.push('/admin/dashboard')" class="logout-btn-header" style="margin-right: 8px; color: #f59e0b;" title="Admin Panel">
          <i class="fas fa-shield-alt"></i>
        </button>
        <button @click="handleLogout" class="logout-btn-header" title="Logout">
          <i class="fas fa-sign-out-alt"></i>
        </button>
      </div>
    </header>

    <div class="main-content-area">
      <!-- Sidebar ONLY for Profile Route -->
      <aside class="sidebar-section profile-sidebar" v-if="route.path === '/profile'">
        <nav class="sidebar-nav">
          <router-link to="/dashboard" class="nav-item">
            <div class="nav-icon"><i class="fas fa-arrow-left"></i></div>
            <span>Back to Dashboard</span>
          </router-link>
          <a href="#" class="nav-item" :class="{ active: activeSection === 'personal-details' }" @click.prevent="changeTab('personal-details')">
            <i class="fas fa-user nav-icon"></i>
            <span>Personal Details</span>
          </a>
          <a href="#" class="nav-item" :class="{ active: activeSection === 'change-password' }" @click.prevent="changeTab('change-password')">
            <i class="fas fa-lock nav-icon"></i>
            <span>Change Password</span>
          </a>
          <a href="#" class="nav-item" :class="{ active: activeSection === 'bank-accounts' }" @click.prevent="changeTab('bank-accounts')">
            <i class="fas fa-university nav-icon"></i>
            <span>Bank Accounts</span>
          </a>
          <a href="#" class="nav-item" :class="{ active: activeSection === 'nominee-details' }" @click.prevent="changeTab('nominee-details')">
            <i class="fas fa-users nav-icon"></i>
            <span>Nominee Details</span>
          </a>
        </nav>
      </aside>

      <main class="page-container" :class="{ 'with-sidebar': route.path === '/profile' }">
        <router-view></router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import '../assets/css/user-layout.css'
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../store'
import { useToast } from 'vue-toastification'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const toast = useToast()

const activeSection = ref(route.query.tab || 'personal-details')

const changeTab = (id) => {
  activeSection.value = id
  router.push({ path: '/profile', query: { tab: id } })
}

const showNotifications = ref(false)
const notifications = ref([])

watch(() => route.path, () => {
  showNotifications.value = false
})

const markAllAsRead = () => {
  notifications.value = []
  showNotifications.value = false
  toast.success('All notifications marked as read')
}

const toggleNotifications = () => {
  showNotifications.value = !showNotifications.value
}

const handleLogout = () => {
  authStore.logout()
  toast.info("Logged out successfully")
  router.push('/login')
}
</script>
