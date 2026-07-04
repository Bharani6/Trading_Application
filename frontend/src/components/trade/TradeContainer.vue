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
        <div class="countdown-timer" style="display: flex; align-items: center; color: var(--text-muted); font-size: 14px; margin-left: 15px;">
          <i class="fas fa-sync-alt" :class="{'fa-spin': isRefreshing}" style="margin-right: 6px;"></i>
          Updates in {{ countdown }}s
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
          <div class="stock-symbol">{{ stock.symbol.split('.')[0] }}</div>
          <div class="stock-name">{{ stock.name }}</div>
        </div>
        <div class="segment">
          <span class="segment-badge">{{ stock.segment }}</span>
        </div>
        <div class="stock-price">
          ₹{{ stock.price.toFixed(2) }}
          <span :class="getPriceChange(stock).class" style="font-size: 0.85em; margin-left: 4px;">
            {{ getPriceChange(stock).text }}
          </span>
        </div>
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
          <button 
            class="watchlist-icon-btn"
            @click="toggleWatchlist(stock)"
            :title="isInWatchlist(stock.id) ? 'Remove from Watchlist' : 'Add to Watchlist'"
          >
            <i class="fas fa-heart" :class="isInWatchlist(stock.id) ? 'text-red-500' : 'text-gray-500'"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Custom Trade Modal -->
    <div v-if="showModal" class="trade-modal-overlay">
      <div class="trade-modal-card">
        <div class="modal-header">
          <h3>{{ modalAction }} {{ modalStock?.symbol.split('.')[0] }}</h3>
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
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useToast } from 'vue-toastification'
import { useTradeStore } from '../../store/trade'
import { useAuthStore } from '../../store'
import { useWalletStore } from '../../store/wallet'
import { useWatchlistStore } from '../../store/watchlist'

const toast = useToast()
const tradeStore = useTradeStore()
const authStore = useAuthStore()
const walletStore = useWalletStore()
const watchlistStore = useWatchlistStore()

// Filter State
const searchQuery = ref('')
const selectedSegment = ref('All Segments')

// Modal State
const showModal = ref(false)
const modalAction = ref('')
const modalStock = ref(null)
const modalQuantity = ref('')

// Auto-refresh State
const countdown = ref(60)
const isRefreshing = ref(false)
let refreshInterval = null

const fetchMarketData = async () => {
  isRefreshing.value = true
  await tradeStore.fetchShares(searchQuery.value)
  isRefreshing.value = false
}

let searchTimeout = null
watch(searchQuery, () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchMarketData()
  }, 300)
})

onMounted(async () => {
  await Promise.all([
    fetchMarketData(),
    tradeStore.fetchHistory(),
    watchlistStore.fetchWatchlist()
  ])

  refreshInterval = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      countdown.value = 60
      fetchMarketData()
    }
  }, 1000)
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
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

const getPriceChange = (stock) => {
  if (!stock.previous_price || stock.previous_price === 0) return { text: '', class: '' }
  const diff = stock.price - stock.previous_price
  const pct = (diff / stock.previous_price) * 100
  const sign = diff >= 0 ? '+' : ''
  const color = diff >= 0 ? 'text-green' : 'text-red'
  return {
    text: `(1d ${sign}${pct.toFixed(3)}%)`,
    class: color
  }
}

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

const isInWatchlist = (stockId) => {
  return watchlistStore.state.watchlist.some(w => w.stockId === stockId)
}

const toggleWatchlist = async (stock) => {
  const watchlistItem = watchlistStore.state.watchlist.find(w => w.stockId === stock.id)
  if (watchlistItem) {
    const res = await watchlistStore.removeFromWatchlist(watchlistItem.id)
    if (res.success) {
      toast.success('Removed from watchlist')
    } else {
      toast.error(res.message)
    }
  } else {
    const res = await watchlistStore.addToWatchlist(stock.id)
    if (res.success) {
      toast.success('Added to watchlist')
    } else {
      toast.error(res.message)
    }
  }
}
</script>

<style scoped>

.watchlist-icon-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #94a3b8;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  padding: 0;
}

.watchlist-icon-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #f1f5f9;
}

.text-red-500 {
  color: #ef4444;
}

.text-gray-500 {
  color: #6b7280;
}
</style>
