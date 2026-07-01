<template>
  <div class="tradepage">
    <div class="header-section">
      <h1 class="page-title">Market Watch</h1>
      <div class="filters-wrap">
        <select class="filter-select" v-model="selectedSegment">
          <option>All Segments</option>
          <option>NSE</option>
          <option>BSE</option>
        </select>
        <div class="search-box">
          <i class="fas fa-search search-icon"></i>
          <input type="text" class="search-input" v-model="searchQuery" placeholder="Search stocks (e.g. RELIANCE)" />
        </div>
      </div>
    </div>

    <div class="market-section">
      <div class="market-row market-header">
        <span>Symbol / Name</span>
        <span>Segment</span>
        <span>Market Price</span>
        <span>Available Qty</span>
        <span>Owned Shares</span>
        <span>Action</span>
      </div>

      <div class="market-row stock-item" v-for="stock in filteredShares" :key="stock.id">
        <div class="company-info">
          <div class="stock-symbol">{{ stock.symbol }}</div>
          <div class="stock-name">{{ stock.name }}</div>
        </div>
        <div class="segment">
          <span class="segment-badge">{{ stock.segment }}</span>
        </div>
        <div class="stock-price">₹{{ stock.price.toFixed(2) }}</div>
        <div class="stock-qty">{{ stock.available_shares.toLocaleString() }}</div>
        <div class="stock-owned">{{ getOwnedShares(stock.id).toLocaleString() }}</div>
        <div class="actions-wrap">
          <button 
            class="action-btn buy-btn" 
            @click="openTradeModal(stock, 'Buy')"
          >Buy</button>
          <button 
            class="action-btn sell-btn" 
            @click="openTradeModal(stock, 'Sell')"
          >Sell</button>
        </div>
      </div>
    </div>

    <!-- Custom Trade Modal -->
    <div v-if="showModal" class="trade-modal-overlay">
      <div class="trade-modal-card">
        <div class="modal-header">
          <h3>{{ modalAction }} {{ modalStock?.symbol }}</h3>
          <button class="close-btn" @click="closeModal"><i class="fas fa-times"></i></button>
        </div>
        <div class="modal-body">
          <div class="stock-summary">
            <div class="summary-row">
              <span>Market Price:</span>
              <span class="highlight-val">₹{{ modalStock?.price.toFixed(2) }}</span>
            </div>
            <div class="summary-row">
              <span>Available Shares:</span>
              <span>{{ modalStock?.available_shares.toLocaleString() }}</span>
            </div>
            <div class="summary-row" v-if="modalAction === 'Buy'">
              <span>Current Balance:</span>
              <span class="highlight-val">₹{{ walletStore.state.balance.available_balance.toFixed(2) }}</span>
            </div>
          </div>
          <div class="input-group">
            <label>Quantity</label>
            <input 
              type="number" 
              class="qty-input" 
              v-model="modalQuantity" 
              min="1" 
              placeholder="Enter number of shares"
              @keyup.enter="confirmTrade"
            />
          </div>
          <div class="total-est" v-if="modalQuantity > 0">
            <span>Estimated Total:</span>
            <span class="highlight-val">₹{{ (modalQuantity * modalStock.price).toFixed(2) }}</span>
          </div>
        </div>
        <div class="modal-footer">
          <button class="cancel-btn" @click="closeModal">Cancel</button>
          <button 
            class="confirm-btn" 
            :class="modalAction === 'Buy' ? 'bg-buy' : 'bg-sell'"
            :disabled="!modalQuantity || modalQuantity <= 0 || tradeStore.state.loading"
            @click="confirmTrade"
          >
            <i v-if="tradeStore.state.loading" class="fas fa-circle-notch fa-spin"></i>
            <span v-else>Confirm {{ modalAction }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import '../../assets/css/trade.css'
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'vue-toastification'
import { useTradeStore } from '../../store/trade'
import { useAuthStore } from '../../store'
import { useWalletStore } from '../../store/wallet'

const toast = useToast()
const tradeStore = useTradeStore()
const authStore = useAuthStore()
const walletStore = useWalletStore()

// Filter State
const searchQuery = ref('')
const selectedSegment = ref('All Segments')

// Modal State
const showModal = ref(false)
const modalAction = ref('')
const modalStock = ref(null)
const modalQuantity = ref('')

onMounted(async () => {
  await Promise.all([
    tradeStore.fetchShares(),
    tradeStore.fetchHistory()
  ])
})

const filteredShares = computed(() => {
  if (!tradeStore.state.shares) return []
  return tradeStore.state.shares.filter(stock => {
    const matchesSegment = selectedSegment.value === 'All Segments' || stock.segment.toUpperCase() === selectedSegment.value.toUpperCase()
    const matchesSearch = !searchQuery.value || 
      stock.symbol.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
      stock.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchesSegment && matchesSearch
  })
})

const getOwnedShares = (shareId) => {
  if (!tradeStore.state.history || !tradeStore.state.history.length) return 0;
  
  const trades = tradeStore.state.history.filter(t => 
    (t.ShareID === shareId || t.share_id === shareId) && 
    (t.Status === 'completed' || t.status === 'completed')
  );
  
  let owned = 0;
  trades.forEach(t => {
    const type = (t.Type || t.type).toLowerCase();
    const qty = t.Quantity || t.quantity || 0;
    if (type === 'buy') owned += qty;
    else if (type === 'sell') owned -= qty;
  });
  return owned;
}

const openTradeModal = async (stock, action) => {
  if (authStore.state.user?.status !== 'active') {
    toast.error('You must be approved by an Admin (KYC verification) to trade.')
    return
  }
  await walletStore.fetchBalance()
  modalStock.value = stock
  modalAction.value = action
  modalQuantity.value = ''
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  modalStock.value = null
  modalQuantity.value = ''
}

const confirmTrade = async () => {
  const qty = parseInt(modalQuantity.value)
  if (!qty || isNaN(qty) || qty <= 0) {
    toast.warning('Please enter a valid quantity')
    return
  }

  const stock = modalStock.value
  const action = modalAction.value

  let res;
  if (action === 'Buy') {
    res = await tradeStore.buyShare(stock.id, qty)
  } else {
    res = await tradeStore.sellShare(stock.id, qty)
  }

  if (res.success) {
    toast.success(`Successfully ${action === 'Buy' ? 'bought' : 'sold'} ${qty} shares of ${stock.symbol}`)
    closeModal()
  } else {
    toast.error(res.message)
  }
}
</script>

<style scoped>
/* Trade Modal CSS */
.trade-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.trade-modal-card {
  background: #131929;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.4);
  animation: modalPop 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

@keyframes modalPop {
  0% { transform: scale(0.9); opacity: 0; }
  100% { transform: scale(1); opacity: 1; }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.modal-header h3 {
  margin: 0;
  font-size: 20px;
  color: #f1f5f9;
  font-weight: 700;
}

.close-btn {
  background: none;
  border: none;
  color: #94a3b8;
  font-size: 18px;
  cursor: pointer;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #f1f5f9;
}

.modal-body {
  padding: 24px;
}

.stock-summary {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
  color: #94a3b8;
}

.summary-row:last-child {
  margin-bottom: 0;
}

.highlight-val {
  color: #f1f5f9;
  font-weight: 600;
}

.input-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #cbd5e1;
  margin-bottom: 8px;
}

.qty-input {
  width: 100%;
  background: #0b0f1a;
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #f1f5f9;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 16px;
  outline: none;
  transition: border-color 0.2s;
}

.qty-input:focus {
  border-color: #3b82f6;
}

.total-est {
  margin-top: 16px;
  display: flex;
  justify-content: space-between;
  font-size: 16px;
  color: #94a3b8;
  border-top: 1px dashed rgba(255, 255, 255, 0.1);
  padding-top: 16px;
}

.total-est .highlight-val {
  font-size: 18px;
  color: #3b82f6;
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.cancel-btn, .confirm-btn {
  flex: 1;
  padding: 12px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 15px;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.cancel-btn {
  background: transparent;
  color: #94a3b8;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.cancel-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #f1f5f9;
}

.confirm-btn {
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.confirm-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.confirm-btn.bg-buy {
  background: #10b981;
}

.confirm-btn.bg-buy:hover:not(:disabled) {
  background: #059669;
}

.confirm-btn.bg-sell {
  background: #f43f5e;
}

.confirm-btn.bg-sell:hover:not(:disabled) {
  background: #e11d48;
}
</style>
