import api from './config'

export const userApi = {
  getProfile: () => {
    return api.get('/users/me')
  },
  
  updateKyc: (payload) => {
    return api.post('/users/kyc', payload)
  }
}
