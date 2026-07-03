import { reactive } from 'vue'
import api from '../services/api'

const state = reactive({
  watchlist: [],
  loading: false
})

export const useWatchlistStore = () => {
  const fetchWatchlist = async () => {
    state.loading = true
    try {
      const res = await api.get('/watchlist')
      if (res.data.success) {
        state.watchlist = res.data.data || []
      }
    } catch (error) {
      console.error('Failed to fetch watchlist:', error)
    } finally {
      state.loading = false
    }
  }

  const addToWatchlist = async (stockId) => {
    try {
      const res = await api.post('/watchlist', { stockId })
      if (res.data.success) {
        await fetchWatchlist()
        return { success: true, message: res.data.message }
      }
      return { success: false, message: 'Add failed' }
    } catch (error) {
      console.error('Failed to add to watchlist:', error.response?.data || error)
      return { success: false, message: error.response?.data?.message || error.message || 'Add error' }
    }
  }

  const removeFromWatchlist = async (id) => {
    try {
      const res = await api.delete(`/watchlist/${id}`)
      if (res.data.success) {
        await fetchWatchlist()
        return { success: true, message: res.data.message }
      }
      return { success: false, message: 'Remove failed' }
    } catch (error) {
      return { success: false, message: error.response?.data?.message || 'Remove error' }
    }
  }

  const updateFavorite = async (id, isFavorite) => {
    try {
      const res = await api.put(`/watchlist/${id}/favorite`, { isFavorite })
      if (res.data.success) {
        await fetchWatchlist()
        return { success: true, message: res.data.message }
      }
      return { success: false, message: 'Update failed' }
    } catch (error) {
      return { success: false, message: error.response?.data?.message || 'Update error' }
    }
  }

  return {
    state,
    fetchWatchlist,
    addToWatchlist,
    removeFromWatchlist,
    updateFavorite
  }
}
