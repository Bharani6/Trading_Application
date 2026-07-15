import api from './config'

export const authApi = {
  login: (credentials) => {
    return api.post('/auth/login', credentials)
  },
  
  register: (userData) => {
    return api.post('/auth/register', userData)
  },

  forgotPassword: (data) => {
    return api.post('/auth/forgot-password', data)
  },

  verifyResetToken: (data) => {
    return api.post('/auth/verify-reset-token', data)
  },

  resetPassword: (data) => {
    return api.post('/auth/reset-password', data)
  },

  sendOTP: (data) => {
    return api.post('/auth/send-otp', data)
  },

  verifyOTP: (data) => {
    return api.post('/auth/verify-otp', data)
  },

  sendEmailOTP: (data) => {
    return api.post('/auth/send-email-otp', data)
  },

  verifyEmailOTP: (data) => {
    return api.post('/auth/verify-email-otp', data)
  }
}
