<template>
  <div class="auth-card">
    <div class="auth-card-header" style="margin-bottom: 2rem;">
      <div style="display: flex; align-items: center; gap: 1rem; margin-bottom: 8px;">
        <div style="width: 56px; height: 56px; background: rgba(15, 23, 42, 0.6); border: 1px solid rgba(139, 92, 246, 0.2); border-radius: 18px; display: flex; align-items: center; justify-content: center; box-shadow: inset 0 0 20px rgba(139, 92, 246, 0.05);">
          <i class="far fa-user" style="color: #a5b4fc; font-size: 24px;"></i>
        </div>
        <h1 style="margin: 0; text-align: left; font-size: 32px; color: #ffffff; font-weight: 700; letter-spacing: -1px;">Welcome back</h1>
      </div>
      <p style="color: #94a3b8; font-size: 15px; margin: 0 0 0 72px; text-align: left;">Login to continue to your account</p>
    </div>

    <form @submit.prevent="handleLogin" class="auth-form" novalidate>
      <div v-if="errorMessage" class="error-banner">
        <i class="fa fa-exclamation-circle"></i>
        {{ errorMessage }}
      </div>
      <div class="form-group">
        <label class="form-label">Email Address</label>
        <div class="input-wrapper">
          <i class="fa fa-envelope input-icon"></i>
          <input
            type="email"
            class="form-input"
            v-model="form.email"
            placeholder="you@example.com"
            required
            pattern="[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$"
            title="Enter a valid email address"
            autocomplete="email"
          />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">Password</label>
        <div class="input-wrapper">
          <i class="fa fa-lock input-icon"></i>
          <input
            :type="showPassword ? 'text' : 'password'"
            class="form-input"
            v-model="form.password"
            placeholder="********"
            required
            autocomplete="current-password"
          />
          <i
            :class="showPassword ? 'fa fa-eye-slash' : 'fa fa-eye'"
            class="password-toggle-icon"
            @click="showPassword = !showPassword"
          ></i>
        </div>
      </div>

      <div class="form-actions">
        <label class="remember-me">
          <input type="checkbox" v-model="form.remember" />
          <span>Remember me</span>
        </label>
        <router-link to="/forgot-password" class="forgot-link">Forgot Password?</router-link>
      </div>

      <button type="submit" class="submit-btn" :disabled="loading" style="background: linear-gradient(90deg, #8B5CF6, #3B82F6); color: white; border: none; padding: 14px; border-radius: 8px; width: 100%; font-weight: 600; font-size: 15px; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 8px;">
        <span v-if="!loading">
          <i class="fas fa-sign-in-alt"></i> Sign In
        </span>
        <i v-else class="fa fa-spinner fa-spin"></i>
      </button>
    </form>

    <div class="register-prompt">
      Don't have an account? 
      <router-link to="/register" class="register-link">Create one free</router-link>
      
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import { useAuthStore } from '../../store'
import { authApi } from '../../api/auth.api'

const router = useRouter()
const toast = useToast()
const authStore = useAuthStore()

const loading = ref(false)
const errorMessage = ref('')
const showPassword = ref(false)
const form = reactive({
  email: '',
  password: '',
  remember: false
})

const showError = (msg) => {
  errorMessage.value = msg
  loading.value = false
}



const handleLogin = async () => {
  loading.value = true
  errorMessage.value = ''

  if (!form.email) return showError('Email address is required.')
  if (!/[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$/i.test(form.email)) return showError('Please enter a valid email address.')
  if (!form.password) return showError('Password is required.')

  try {
    const res = await authApi.login({
      email: form.email,
      password: form.password
    })

    if (res.data.success) {
      const loginData = res.data.data
      authStore.login(loginData.data, loginData.access_token)
      toast.success(`Welcome back, ${loginData.data?.name || 'Trader'}!`)
      
      if (loginData.data?.role === 'admin') {
        router.push('/admin/dashboard')
      } else {
        router.push('/dashboard')
      }
    }
  } catch (error) {
    const msg = error.response?.data?.message || 'Login failed. Please check credentials.'
    errorMessage.value = msg
    toast.error(msg)
  } finally {
    loading.value = false
  }
}
</script>
