import { reactive } from 'vue'
import api from '../services/api'

const state = reactive({
  shares: [],
  history: [],
  loading: false
})

export const useTradeStore = () => {
  const fetchShares = async (search = '') => {
    try {
      const res = await api.get('/shares', { params: { search } })
      if (res.data.success) {
        state.shares = res.data.data || []
      }
    } catch (error) {
      console.error('Failed to fetch shares:', error)
    }
  }

  const fetchHistory = async () => {
    try {
      const res = await api.get('/trades/history')
      if (res.data.success) {
        state.history = res.data.data || []
      }
    } catch (error) {
      console.error('Failed to fetch trade history:', error)
    }
  }

  const buyShare = async (shareId, quantity) => {
    state.loading = true
    try {
      const res = await api.post('/trades/buy', { share_id: shareId, quantity: parseInt(quantity) })
      if (res.data.success) {
        await fetchShares()
        await fetchHistory()
        return { success: true, message: res.data.message }
      }
      return { success: false, message: 'Buy failed' }
    } catch (error) {
      return { success: false, message: error.response?.data?.message || 'Buy error' }
    } finally {
      state.loading = false
    }
  }

  const sellShare = async (shareId, quantity) => {
    state.loading = true
    try {
      const res = await api.post('/trades/sell', { share_id: shareId, quantity: parseInt(quantity) })
      if (res.data.success) {
        await fetchShares()
        await fetchHistory()
        return { success: true, message: res.data.message }
      }
      return { success: false, message: 'Sell failed' }
    } catch (error) {
      return { success: false, message: error.response?.data?.message || 'Sell error' }
    } finally {
      state.loading = false
    }
  }

  const cancelTrade = async (tradeId) => {
    state.loading = true
    try {
      const res = await api.post(`/trades/${tradeId}/cancel`)
      if (res.data.success) {
        await fetchHistory()
        return { success: true, message: res.data.message }
      }
      return { success: false, message: 'Cancel failed' }
    } catch (error) {
      return { success: false, message: error.response?.data?.message || 'Cancel error' }
    } finally {
      state.loading = false
    }
  }

  return {
    state,
    fetchShares,
    fetchHistory,
    buyShare,
    sellShare,
    cancelTrade
  }
}
