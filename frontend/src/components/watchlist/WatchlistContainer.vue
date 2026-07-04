<template>
  <div class="tradepage">
    <div class="header-section">
      <h1 class="page-title">My Watchlist</h1>
      <div class="filters-wrap">
        <select class="filter-select" v-model="sortBy">
          <option value="name">Sort by Name</option>
          <option value="price">Sort by Price</option>
        </select>
        <div class="search-box">
          <i class="fas fa-search search-icon"></i>
          <input type="text" class="search-input" v-model="searchQuery" placeholder="Search watchlist..." />
        </div>
      </div>
    </div>

    <div class="market-section" v-if="filteredWatchlist.length > 0">
      <div class="market-row watchlist-row market-header">
        <span>Symbol / Name</span>
        <span>Market Price</span>
        <span>1D Change</span>
        <span>Action</span>
      </div>

      <div class="market-row watchlist-row stock-item" v-for="item in filteredWatchlist" :key="item.id">
        <div class="company-info">
          <div class="stock-symbol">{{ item.symbol.split('.')[0] }}</div>
          <div class="stock-name">{{ item.stockName }}</div>
        </div>
        <div class="stock-price">
          ₹{{ item.currentPrice.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2}) }}
        </div>
        <div class="stock-change" :class="getChangeClass(item.currentPrice, item.previousPrice)">
          {{ getChangeText(item.currentPrice, item.previousPrice) }}
        </div>
        <div class="actions-wrap">
          <button class="action-btn buy-btn" @click="openModal('Buy', item)">B</button>
          <button class="action-btn sell-btn" @click="openModal('Sell', item)">S</button>
          <button class="action-btn remove-btn" @click="removeItem(item.id)">
            <i class="fas fa-trash-alt"></i>
          </button>
        </div>
      </div>
    </div>
    
    <div class="empty-state" v-else>
      <div class="empty-content">
        <i class="fas fa-chart-line empty-icon"></i>
        <h3>No Stocks in Watchlist</h3>
        <p>You haven't added any stocks to your watchlist yet. Go to Market to add some.</p>
        <router-link to="/trade" class="browse-market-btn">Browse Market</router-link>
      </div>
    </div>

    <!-- Trade Modal -->
    <div class="trade-modal-overlay" v-if="showModal" @click.self="closeModal">
      <div class="trade-modal-card">
        <div class="modal-header">
          <h3>{{ modalAction }} {{ modalStock?.symbol.split('.')[0] }}</h3>
          <button class="close-btn" @click="closeModal"><i class="fas fa-times"></i></button>
        </div>
        <div class="modal-body">
          <div class="stock-summary">
            <div class="summary-row">
              <span>Market Price:</span>
              <span class="highlight-val">₹{{ modalStock?.currentPrice.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2}) }}</span>
            </div>
            <div class="summary-row" v-if="modalAction === 'Buy' && walletStore.state.balance">
              <span>Current Balance:</span>
              <span class="highlight-val">₹{{ walletStore.state.balance.available_balance.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2}) }}</span>
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
            <span class="highlight-val">₹{{ (modalQuantity * modalStock.currentPrice).toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2}) }}</span>
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
import { useWatchlistStore } from '../../store/watchlist'
import { useTradeStore } from '../../store/trade'
import { useWalletStore } from '../../store/wallet'

const toast = useToast()
const watchlistStore = useWatchlistStore()
const tradeStore = useTradeStore()
const walletStore = useWalletStore()

const searchQuery = ref('')
const sortBy = ref('name')

// Trade Modal State
const showModal = ref(false)
const modalAction = ref('')
const modalStock = ref(null)
const modalQuantity = ref('')

onMounted(async () => {
  await watchlistStore.fetchWatchlist()
  await walletStore.fetchBalance()
})

const filteredWatchlist = computed(() => {
  let list = [...watchlistStore.state.watchlist]
  
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(item => 
      item.stockName.toLowerCase().includes(q) || 
      item.symbol.toLowerCase().includes(q)
    )
  }
  
  list.sort((a, b) => {
    if (sortBy.value === 'price') {
      return b.currentPrice - a.currentPrice
    } else {
      return a.stockName.localeCompare(b.stockName)
    }
  })
  
  return list
})

const getChangeClass = (current, previous) => {
  if (!previous) return 'text-gray-500'
  return current >= previous ? 'text-green-500' : 'text-red-500'
}

const getChangeText = (current, previous) => {
  if (!previous) return 'N/A'
  const diff = current - previous
  const percent = (diff / previous) * 100
  const sign = diff >= 0 ? '+' : ''
  return `${sign}${diff.toFixed(2)} (${sign}${percent.toFixed(2)}%)`
}

const removeItem = async (id) => {
  const res = await watchlistStore.removeFromWatchlist(id)
  if (res.success) {
    toast.success('Removed from watchlist')
  } else {
    toast.error(res.message)
  }
}

// Modal Actions
const openModal = (action, item) => {
  modalAction.value = action
  modalStock.value = item
  modalQuantity.value = ''
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  modalStock.value = null
  modalQuantity.value = ''
}

const confirmTrade = async () => {
  if (!modalQuantity.value || modalQuantity.value <= 0) return
  
  try {
    if (modalAction.value === 'Buy') {
      const res = await tradeStore.buyShare(modalStock.value.stockId, modalQuantity.value)
      if (res.success) {
        toast.success(`Successfully placed order to buy ${modalQuantity.value} shares of ${modalStock.value.symbol}`)
        closeModal()
        walletStore.fetchBalance()
      } else {
        toast.error(res.message)
      }
    } else {
      const res = await tradeStore.sellShare(modalStock.value.stockId, modalQuantity.value)
      if (res.success) {
        toast.success(`Successfully placed order to sell ${modalQuantity.value} shares of ${modalStock.value.symbol}`)
        closeModal()
        walletStore.fetchBalance()
      } else {
        toast.error(res.message)
      }
    }
  } catch (err) {
    toast.error(err.message || 'An error occurred')
  }
}
</script>

<style scoped>
.text-green-500 {
  color: #10b981;
}
.text-red-500 {
  color: #ef4444;
}
.text-gray-500 {
  color: #6b7280;
}
.stock-change {
  font-weight: 500;
  display: flex;
  align-items: center;
}
.remove-btn {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
}
.remove-btn:hover {
  background: rgba(239, 68, 68, 0.2);
}
.watchlist-row {
  grid-template-columns: 2fr 1fr 1.5fr 1.5fr !important;
}

@media (max-width: 768px) {
  .watchlist-row {
    grid-template-columns: 2fr 1fr 1fr 1.2fr !important;
  }
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: 12px;
  margin-top: 20px;
}
.empty-content {
  text-align: center;
}
.empty-icon {
  font-size: 48px;
  color: #4b5563;
  margin-bottom: 16px;
}
.empty-content h3 {
  color: #f3f4f6;
  font-size: 24px;
  margin-bottom: 8px;
}
.empty-content p {
  color: #9ca3af;
  margin-bottom: 24px;
}
.browse-market-btn {
  display: inline-block;
  background: #3b82f6;
  color: white;
  padding: 10px 20px;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 500;
  transition: background 0.2s;
}
.browse-market-btn:hover {
  background: #2563eb;
}
</style>
