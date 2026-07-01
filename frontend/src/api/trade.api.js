import api from './config'

export const tradeApi = {
  getShares: () => {
    return api.get('/shares')
  },
  
  getHistory: () => {
    return api.get('/trades/history')
  },

  buyShare: (shareId, quantity) => {
    return api.post('/trades/buy', { share_id: shareId, quantity: parseInt(quantity) })
  },

  sellShare: (shareId, quantity) => {
    return api.post('/trades/sell', { share_id: shareId, quantity: parseInt(quantity) })
  },

  cancelTrade: (tradeId) => {
    return api.post(`/trades/${tradeId}/cancel`)
  }
}
