<template>
  <div class="userlayout adminlayout">
    <main class="main-section">
      <header class="top-header">
        
        <div class="header-left">
          <i class="fas fa-shield-alt logo-icon text-primary"></i>
          <h2 class="logo-text">Admin Panel</h2>
        </div>

        <div class="header-actions">
          <button @click="$router.push('/dashboard')" class="header-icon-btn" title="Back to User View">
            <i class="fas fa-user-circle"></i>
          </button>
          <button @click="handleLogout" class="header-icon-btn text-danger" title="Logout">
            <i class="fas fa-sign-out-alt"></i>
          </button>
          <div class="user-profile">
            <div class="avatar-icon bg-gradient-admin">
              <i class="fas fa-user-shield"></i>
            </div>
            <div class="user-info">
              <span class="user-name">{{ authStore.state.user?.name || 'Admin' }}</span>
              <span class="user-status text-danger">Administrator</span>
            </div>
          </div>
        </div>

      </header>

      <nav class="sub-nav-bar">
        <router-link to="/admin/dashboard" class="nav-item" active-class="active">
          <i class="fas fa-tachometer-alt nav-icon"></i>
          <span>Overview</span>
        </router-link>
        <router-link to="/admin/users" class="nav-item" active-class="active">
          <i class="fas fa-users-cog nav-icon"></i>
          <span>User Approvals</span>
        </router-link>
        <router-link to="/admin/stocks" class="nav-item" active-class="active">
          <i class="fas fa-file-csv nav-icon"></i>
          <span>Stock Management</span>
        </router-link>
        <router-link to="/admin/support" class="nav-item" active-class="active">
          <i class="fas fa-envelope nav-icon"></i>
          <span>Support</span>
        </router-link>
      </nav>

      <div class="page-container">
        <router-view></router-view>
      </div>
    </main>
  </div>
</template>

<script setup>
import '../assets/css/admin-layout.css'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store'
import { useToast } from 'vue-toastification'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const handleLogout = () => {
  authStore.logout()
  toast.info("Logged out successfully")
  router.push('/login')
}
</script>
