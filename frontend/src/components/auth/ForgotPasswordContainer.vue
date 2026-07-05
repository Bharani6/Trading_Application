<template>
  <div class="loginpage">
    <div class="bg-particles">
      <div class="particle p1"></div>
      <div class="particle p2"></div>
      <div class="particle p3"></div>
      <div class="particle p4"></div>
      <div class="particle p5"></div>
      <div class="particle p6"></div>
      <div class="particle p7"></div>
      <div class="particle p8"></div>
    </div>

    <div class="auth-section">
      <div class="auth-card">
        <div class="auth-card-header">
          <div class="brand-header">
            <img src="../../assets/icons/favicon.svg" alt="MoneymakePro" class="auth-brand-logo" />
            <span class="auth-brand-name">MoneymakePro</span>
          </div>
          <h1 class="auth-title">Forgot Password</h1>
          <p style="color: rgba(255, 255, 255, 0.7); font-size: 0.9rem; text-align: center; margin-top: 8px;">
            Enter your email to receive a password reset link.
          </p>
        </div>

        <form @submit.prevent="handleForgotPassword" class="auth-form" novalidate>
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
                v-model="email"
                placeholder="you@example.com"
                required
                pattern="[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$"
                title="Enter a valid email address"
                autocomplete="email"
              />
            </div>
          </div>

          <button type="submit" class="submit-btn" :disabled="loading">
            <span v-if="!loading">
              <i class="fa fa-paper-plane" style="margin-right:8px;"></i>Send Reset Link
            </span>
            <i v-else class="fa fa-spinner fa-spin"></i>
          </button>
        </form>

        <div class="register-prompt">
          Remembered your password? 
          <router-link to="/login" class="register-link">Back to login</router-link>
        </div>

        <div v-if="mockLink" class="mock-link-container" style="margin-top: 20px; padding: 15px; background: rgba(0, 230, 118, 0.1); border-left: 4px solid #00e676; border-radius: 4px;">
          <p style="color: #fff; font-size: 0.9rem; margin-bottom: 8px;">
            <i class="fa fa-info-circle"></i> Mock Email Link (Visible for {{ countdown }}s):
          </p>
          <a :href="mockLink" style="color: #00e676; word-break: break-all; font-size: 0.85rem;">{{ mockLink }}</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import '../../assets/css/login.css'
import { ref } from 'vue'
import { useToast } from 'vue-toastification'
import { authApi } from '../../api/auth.api'

const toast = useToast()
const loading = ref(false)
const errorMessage = ref('')
const email = ref('')

const mockLink = ref('')
const countdown = ref(10)
let timer = null

const showError = (msg) => {
  errorMessage.value = msg
  loading.value = false
}

const handleForgotPassword = async () => {
  loading.value = true
  errorMessage.value = ''

  if (!email.value) return showError('Email address is required.')
  if (!/[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$/i.test(email.value)) return showError('Please enter a valid email address.')

  try {
    const res = await authApi.forgotPassword({ email: email.value })
    if (res.data.success) {
      toast.success(res.data.message)
      
      // Mocking email delivery for testing
      if (res.data.data && res.data.data.mock_token) {
        const resetLink = `${window.location.origin}/reset-password?token=${res.data.data.mock_token}`
        
        mockLink.value = resetLink
        countdown.value = 10
        
        if (timer) clearInterval(timer)
        
        timer = setInterval(() => {
          countdown.value--
          if (countdown.value <= 0) {
            clearInterval(timer)
            mockLink.value = ''
          }
        }, 1000)
      }
      
      email.value = ''
    }
  } catch (error) {
    const msg = error.response?.data?.message || 'Failed to send reset link.'
    errorMessage.value = msg
    toast.error(msg)
  } finally {
    loading.value = false
  }
}
</script>
