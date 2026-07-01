<template>
  <div class="auth-card">
    <div class="auth-card-header">
      <div class="brand-header">
        <img src="../../assets/icons/favicon.svg" alt="MoneymakePro" class="auth-brand-logo" />
        <span class="auth-brand-name">MoneymakePro</span>
      </div>
      <h1 class="auth-title">Welcome back</h1>
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
            placeholder="••••••••"
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
        <a href="#" class="forgot-link" @click.prevent="handleForgotPassword">Forgot Password?</a>
      </div>

      <button type="submit" class="submit-btn" :disabled="loading">
        <span v-if="!loading">
          <i class="fa fa-sign-in-alt" style="margin-right:8px;"></i>Sign In
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

const handleForgotPassword = () => {
  toast.info('Please Contact Administrator')
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
