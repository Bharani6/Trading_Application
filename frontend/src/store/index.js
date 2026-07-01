import { reactive } from 'vue'

const savedUser = (() => {
  try { return JSON.parse(localStorage.getItem('user') || 'null') } catch { return null }
})()

const state = reactive({
  user: savedUser,
  isAuthenticated: !!localStorage.getItem('token'),
  token: localStorage.getItem('token') || null
})

export const useAuthStore = () => {
  const login = (userData, token) => {
    state.user = userData
    state.isAuthenticated = true
    state.token = token
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  const logout = () => {
    state.user = null
    state.isAuthenticated = false
    state.token = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    state,
    login,
    logout
  }
}
