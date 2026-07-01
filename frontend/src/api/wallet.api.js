import api from './config'

export const walletApi = {
  getBalance: () => {
    return api.get('/wallet/balance')
  },
  
  getTransactions: () => {
    return api.get('/wallet/transactions')
  },

  addFunds: (amount) => {
    return api.post('/wallet/add-fund', { amount: parseFloat(amount) })
  },

  withdrawFunds: (amount) => {
    return api.post('/wallet/withdraw', { amount: parseFloat(amount) })
  }
}
