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
          <h1 class="auth-title">Reset Password</h1>
          <p style="color: rgba(255, 255, 255, 0.7); font-size: 0.9rem; text-align: center; margin-top: 8px;">
            Create a new strong password for your account.
          </p>
        </div>

        <!-- Token Validation State -->
        <div v-if="isValidating" class="validation-loading" style="text-align:center; padding: 20px;">
          <i class="fa fa-circle-notch fa-spin fa-2x" style="color:#00e676;"></i>
          <p style="color:#fff; margin-top: 10px;">Verifying token...</p>
        </div>

        <div v-else-if="!tokenValid" class="validation-error" style="text-align:center; padding: 20px;">
          <i class="fa fa-times-circle fa-2x" style="color:#ff1744;"></i>
          <p style="color:#fff; margin-top: 10px;">{{ tokenError }}</p>
          <router-link to="/forgot-password" class="submit-btn" style="display:inline-block; margin-top:15px; text-decoration:none;">Request New Link</router-link>
        </div>

        <!-- Reset Form -->
        <form v-else @submit.prevent="handleResetPassword" class="auth-form" novalidate>
          <div v-if="errorMessage" class="error-banner">
            <i class="fa fa-exclamation-circle"></i>
            {{ errorMessage }}
          </div>

          <div class="form-group">
            <label class="form-label">New Password</label>
            <div class="input-wrapper">
              <i class="fa fa-lock input-icon"></i>
              <input
                :type="showPassword ? 'text' : 'password'"
                class="form-input"
                v-model="password"
                placeholder="••••••••"
                required
              />
              <i
                :class="showPassword ? 'fa fa-eye-slash' : 'fa fa-eye'"
                class="password-toggle-icon"
                @click="showPassword = !showPassword"
              ></i>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Confirm Password</label>
            <div class="input-wrapper">
              <i class="fa fa-lock input-icon"></i>
              <input
                :type="showConfirmPassword ? 'text' : 'password'"
                class="form-input"
                v-model="confirmPassword"
                placeholder="••••••••"
                required
              />
              <i
                :class="showConfirmPassword ? 'fa fa-eye-slash' : 'fa fa-eye'"
                class="password-toggle-icon"
                @click="showConfirmPassword = !showConfirmPassword"
              ></i>
            </div>
          </div>

          <button type="submit" class="submit-btn" :disabled="loading">
            <span v-if="!loading">
              <i class="fa fa-check" style="margin-right:8px;"></i>Reset Password
            </span>
            <i v-else class="fa fa-spinner fa-spin"></i>
          </button>
        </form>

        <div class="register-prompt" v-if="!isValidating && !tokenValid">
          Remembered your password? 
          <router-link to="/login" class="register-link">Back to login</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import '../../assets/css/login.css'
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import { authApi } from '../../api/auth.api'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const token = ref(route.query.token || '')
const isValidating = ref(true)
const tokenValid = ref(false)
const tokenError = ref('')

const loading = ref(false)
const errorMessage = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)

onMounted(async () => {
  if (!token.value) {
    isValidating.value = false
    tokenValid.value = false
    tokenError.value = 'No reset token provided in the URL.'
    return
  }

  try {
    const res = await authApi.verifyResetToken({ token: token.value })
    if (res.data.success && res.data.data.valid) {
      tokenValid.value = true
    }
  } catch (error) {
    tokenValid.value = false
    tokenError.value = error.response?.data?.message || 'Invalid or expired reset token.'
  } finally {
    isValidating.value = false
  }
})

const showError = (msg) => {
  errorMessage.value = msg
  loading.value = false
}

const handleResetPassword = async () => {
  loading.value = true
  errorMessage.value = ''

  if (!password.value) return showError('Password is required.')
  if (password.value !== confirmPassword.value) return showError('Passwords do not match.')

  try {
    const res = await authApi.resetPassword({ 
      token: token.value,
      password: password.value
    })
    
    if (res.data.success) {
      toast.success(res.data.message || 'Password updated successfully. Please login.')
      router.push('/login')
    }
  } catch (error) {
    const msg = error.response?.data?.message || 'Failed to reset password.'
    errorMessage.value = msg
    toast.error(msg)
  } finally {
    loading.value = false
  }
}
</script>
