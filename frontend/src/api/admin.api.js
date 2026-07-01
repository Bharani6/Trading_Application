import api from './config'

export const adminApi = {
  getUsers: () => {
    return api.get('/admin/users')
  },
  
  approveUser: (userId) => {
    return api.put(`/admin/users/${userId}/approve`)
  },

  rejectUser: (userId) => {
    return api.put(`/admin/users/${userId}/reject`)
  },

  updateUserAction: (userId, action) => {
    return api.put(`/admin/users/${userId}/${action}`)
  },

  deleteShares: () => {
    return api.delete('/admin/shares')
  },

  uploadShares: (formData) => {
    return api.post('/admin/shares/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }
}
