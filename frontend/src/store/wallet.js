import { reactive } from 'vue'
import api from '../services/api'

const state = reactive({
  balance: {
    wallet_balance: 0,
    blocked_balance: 0,
    available_balance: 0
  },
  transactions: [],
  loading: false
})

export const useWalletStore = () => {
  const fetchBalance = async () => {
    try {
      const res = await api.get('/wallet/balance')
      if (res.data.success) {
        state.balance = res.data.data
      }
    } catch (error) {
      console.error('Failed to fetch wallet balance:', error)
    }
  }

  const fetchTransactions = async () => {
    try {
      const res = await api.get('/wallet/transactions')
      if (res.data.success) {
        state.transactions = res.data.data || []
      }
    } catch (error) {
      console.error('Failed to fetch transactions:', error)
    }
  }

  const addFunds = async (amount) => {
    state.loading = true
    try {
      const res = await api.post('/wallet/add-fund', { amount: parseFloat(amount) })
      if (res.data.success) {
        await fetchBalance()
        await fetchTransactions()
        return { success: true }
      }
      return { success: false, message: 'Deposit failed' }
    } catch (error) {
      return { success: false, message: error.response?.data?.message || 'Deposit error' }
    } finally {
      state.loading = false
    }
  }

  const withdrawFunds = async (amount) => {
    state.loading = true
    try {
      const res = await api.post('/wallet/withdraw', { amount: parseFloat(amount) })
      if (res.data.success) {
        await fetchBalance()
        await fetchTransactions()
        return { success: true }
      }
      return { success: false, message: 'Withdrawal failed' }
    } catch (error) {
      return { success: false, message: error.response?.data?.message || 'Withdrawal error' }
    } finally {
      state.loading = false
    }
  }

  return {
    state,
    fetchBalance,
    fetchTransactions,
    addFunds,
    withdrawFunds
  }
}
