import api from './config'

export const userApi = {
  getProfile: () => {
    return api.get('/users/me')
  },
  
  updateKyc: (payload) => {
    return api.post('/users/kyc', payload)
  },

  requestAccountClosure: () => {
    return api.post('/users/closure')
  },

  changePassword: (payload) => {
    return api.post('/users/change-password', payload)
  }
}
