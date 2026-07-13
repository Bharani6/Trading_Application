<template>
  <div class="dashboardpage">


    <!-- SUB NAVIGATION TABS -->
    <div class="dashboard-tabs">
      <router-link :to="{ query: { tab: 'explore' } }" class="tab-item" :class="{ active: activeTab === 'explore' }">Explore</router-link>
      <router-link :to="{ query: { tab: 'holdings' } }" class="tab-item" :class="{ active: activeTab === 'holdings' }">Holdings</router-link>
      <router-link :to="{ query: { tab: 'orders' } }" class="tab-item" :class="{ active: activeTab === 'orders' }">Orders</router-link>
      <router-link :to="{ query: { tab: 'watchlist' } }" class="tab-item" :class="{ active: activeTab === 'watchlist' }">Watchlist</router-link>
    </div>

    <!-- MARKET INDICES BANNER -->
    <div class="market-indices-banner" v-if="marketIndices.length">
      <div class="marquee-content">
        <div v-for="(idx, i) in marketIndices" :key="i" class="index-item">
          <span class="index-name">{{ idx.name }}:</span>
          <span class="index-val">{{ formatNumber(idx.current) }}</span>
          <span :class="['index-change', idx.change >= 0 ? 'text-green' : 'text-red']">
            ({{ idx.change >= 0 ? '+' : '' }}{{ formatNumber(idx.change) }}, 
            {{ idx.change >= 0 ? '+' : '' }}{{ formatNumber(idx.change_pct) }}%)
          </span>
        </div>
      </div>
    </div>

    <!-- 2 COLUMN GRID -->
    <div class="dashboard-grid-2col">
      <!-- LEFT COLUMN -->
      <div class="main-column">
        
        <!-- EXPLORE TAB -->
        <template v-if="activeTab === 'explore'">
          <div class="stats-section">
            <div class="statistics-card" @click="goToHoldings" style="cursor: pointer; width: 100%; display: flex; flex-direction: column;">
              <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
                <h3 class="statistics-label" style="margin: 0; display: flex; align-items: center; gap: 0.5rem;">
                  <div class="icon-wrap bg-purple" style="width: 32px; height: 32px; border-radius: 8px; display: flex; justify-content: center; align-items: center;">
                    <i class="fas fa-chart-pie" style="font-size: 14px;"></i>
                  </div>
                  Your Investments <i class="fas fa-chevron-right" style="font-size: 12px; color: #888;"></i>
                </h3>
              </div>
              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
                <div>
                  <div class="text-muted" style="font-size: 12px; margin-bottom: 4px;">Current</div>
                  <div style="font-size: 18px; font-weight: 600;">₹{{ formatNumber(portfolioStats.currentVal) }}</div>
                </div>
                <div>
                  <div class="text-muted" style="font-size: 12px; margin-bottom: 4px;">Invested</div>
                  <div style="font-size: 16px; font-weight: 500;">₹{{ formatNumber(portfolioStats.invested) }}</div>
                </div>
                <div>
                  <div class="text-muted" style="font-size: 12px; margin-bottom: 4px;">1D Returns</div>
                  <div :class="portfolioStats.oneDayRet >= 0 ? 'text-green' : 'text-red'" style="font-size: 14px; font-weight: 500;">
                    {{ portfolioStats.oneDayRet >= 0 ? '+' : '' }}₹{{ formatNumber(portfolioStats.oneDayRet) }} ({{ portfolioStats.oneDayRet >= 0 ? '+' : '' }}{{ formatNumber(portfolioStats.oneDayRetPct) }}%)
                  </div>
                </div>
                <div>
                  <div class="text-muted" style="font-size: 12px; margin-bottom: 4px;">Total Returns</div>
                  <div :class="portfolioStats.totalRet >= 0 ? 'text-green' : 'text-red'" style="font-size: 14px; font-weight: 500;">
                    {{ portfolioStats.totalRet >= 0 ? '+' : '' }}₹{{ formatNumber(portfolioStats.totalRet) }} ({{ portfolioStats.totalRet >= 0 ? '+' : '' }}{{ formatNumber(portfolioStats.totalRetPct) }}%)
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="charts-section">
            <div class="chart-card">
              <div class="card-header" style="display: flex; justify-content: space-between; align-items: center;">
                <h2 class="card-title">Portfolio Analytics</h2>
                <div class="chart-toggles">
                  <button :class="['chart-toggle', timeFrame === '1D' ? 'active' : '']" @click="setTimeFrame('1D')">1D</button>
                  <button :class="['chart-toggle', timeFrame === '1W' ? 'active' : '']" @click="setTimeFrame('1W')">1W</button>
                  <button :class="['chart-toggle', timeFrame === '1M' ? 'active' : '']" @click="setTimeFrame('1M')">1M</button>
                  <button :class="['chart-toggle', timeFrame === '1Y' ? 'active' : '']" @click="setTimeFrame('1Y')">1Y</button>
                  <button :class="['chart-toggle', timeFrame === 'ALL' ? 'active' : '']" @click="setTimeFrame('ALL')">ALL</button>
                </div>
              </div>
              <div class="card-body">
                <canvas id="portfolioChart"></canvas>
              </div>
            </div>

            <div class="chart-card side-card">
              <div class="card-header">
                <h2 class="card-title">Recent Activity</h2>
                <router-link to="/billing" class="view-all-link">View All</router-link>
              </div>
              <div class="activity-list">
                <div v-for="tx in recentTransactions" :key="tx.id" class="activity-item">
                  <div :class="['activity-icon', tx.type.includes('sell') || tx.type === 'add_fund' ? 'icon-green' : 'icon-red']">
                    <i :class="tx.type === 'add_fund' ? 'fas fa-plus' : tx.type.includes('sell') ? 'fas fa-arrow-up' : 'fas fa-arrow-down'"></i>
                  </div>
                  <div class="activity-info">
                    <p class="activity-desc">{{ tx.description }}</p>
                    <span class="activity-date">{{ formatDate(tx.created_at) }}</span>
                  </div>
                  <div :class="['activity-amount', tx.type === 'add_fund' || tx.type === 'trade_sell' ? 'amount-pos' : 'amount-neg']">
                    {{ tx.type === 'add_fund' || tx.type === 'trade_sell' ? '+' : '-' }}₹{{ formatNumber(tx.amount) }}
                  </div>
                </div>
                <div v-if="!recentTransactions.length" class="empty-state">
                  <i class="fas fa-receipt"></i>
                  <p>No transactions yet</p>
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- HOLDINGS TAB -->
        <template v-if="activeTab === 'holdings'">
          <div class="holdings-section kite-theme">
            <!-- Holdings Summary Header -->
            <div class="holdings-summary-card">
              <div class="holdings-summary-top">
                <div class="holdings-summary-title-group">
                  <span class="holdings-summary-title">HOLDINGS ({{ portfolioStats.holdingsList.length }}) <i class="fas fa-caret-down"></i></span>
                  <div class="holdings-summary-current-val">{{ showValues ? '₹' + formatNumber(portfolioStats.currentVal) : '***' }}</div>
                </div>
                <div class="holdings-summary-actions">
                  <button class="holdings-summary-icon-btn" @click="toggleValues"><i class="fas" :class="showValues ? 'fa-eye' : 'fa-eye-slash'"></i></button>
                  <button class="holdings-summary-btn" @click="toggleAnalyse" :class="{ 'active': showAnalyse }"><i class="fas fa-chart-bar"></i> Analyse</button>
                </div>
              </div>
              <div class="holdings-summary-bottom">
                <div class="holdings-summary-stat">
                  <span class="holdings-summary-label">Invested value</span>
                  <span class="holdings-summary-val">{{ showValues ? '₹' + formatNumber(portfolioStats.invested) : '***' }}</span>
                </div>
                <div class="holdings-summary-stat text-center">
                  <span class="holdings-summary-label">1D returns</span>
                  <span class="holdings-summary-val" :class="portfolioStats.oneDayRet >= 0 ? 'text-green' : 'text-red'">
                    <template v-if="showValues">{{ portfolioStats.oneDayRet >= 0 ? '+' : '' }}₹{{ formatNumber(portfolioStats.oneDayRet) }} ({{ portfolioStats.oneDayRet >= 0 ? '+' : '' }}{{ formatNumber(portfolioStats.oneDayRetPct) }}%)</template>
                    <template v-else>***</template>
                  </span>
                </div>
                <div class="holdings-summary-stat text-right">
                  <span class="holdings-summary-label">Total returns</span>
                  <span class="holdings-summary-val" :class="portfolioStats.totalRet >= 0 ? 'text-green' : 'text-red'">
                    <template v-if="showValues">{{ portfolioStats.totalRet >= 0 ? '+' : '' }}₹{{ formatNumber(portfolioStats.totalRet) }} ({{ portfolioStats.totalRet >= 0 ? '+' : '' }}{{ formatNumber(portfolioStats.totalRetPct) }}%)</template>
                    <template v-else>***</template>
                  </span>
                </div>
              </div>

              <!-- Analyse Graph Section -->
              <div class="analyse-section" v-if="showAnalyse">
                <canvas id="analyseChart" style="height: 200px; width: 100%;"></canvas>
              </div>
            </div>

            <!-- Holdings List -->
            <div class="holdings-list-card">
              <div class="holdings-list-header">
                <div class="holdings-list-col holdings-list-company">Company <i class="fas fa-caret-down"></i></div>
                <div class="holdings-list-col holdings-list-chart"></div>
                <div class="holdings-list-col holdings-list-market text-right">Market price (1D%) <i class="fas fa-caret-down"></i></div>
                <div class="holdings-list-col holdings-list-returns text-right">Returns (%) <i class="fas fa-caret-down"></i></div>
                <div class="holdings-list-col holdings-list-current text-right">Current (Invested) <i class="fas fa-caret-down"></i></div>
              </div>
              <div class="holdings-list-body">
                <div v-for="h in portfolioStats.holdingsList" :key="h.id" class="holdings-list-row clickable-row" @click="openTradePanel(h)">
                  <div class="holdings-list-col holdings-list-company">
                    <div class="holdings-list-name">{{ h.name }}</div>
                    <div class="holdings-list-sub">{{ h.qty }} shares • Avg. ₹{{ formatNumber(h.avgPrice) }}</div>
                  </div>
                  <div class="holdings-list-col holdings-list-chart">
                    <!-- SVG sparkline placeholder -->
                    <svg viewBox="0 0 100 30" class="sparkline" :class="h.returns >= 0 ? 'spark-green' : 'spark-red'">
                      <polyline fill="none" stroke-width="1.5" points="0,15 10,25 20,10 30,20 40,5 50,15 60,10 70,20 80,15 90,25 100,10" />
                    </svg>
                  </div>
                  <div class="holdings-list-col holdings-list-market text-right">
                    <div class="holdings-list-val">₹{{ formatNumber(h.currentPrice) }}</div>
                    <div class="holdings-list-sub" :class="h.oneDayRet >= 0 ? 'text-green' : 'text-red'">
                      {{ h.oneDayRet >= 0 ? '+' : '' }}₹{{ formatNumber(h.oneDayRet) }} ({{ h.oneDayRet >= 0 ? '+' : '' }}{{ formatNumber(h.oneDayRetPct) }}%)
                    </div>
                  </div>
                  <div class="holdings-list-col holdings-list-returns text-right">
                    <div class="holdings-list-val" :class="h.returns >= 0 ? 'text-green' : 'text-red'">
                      {{ h.returns >= 0 ? '+' : '' }}₹{{ formatNumber(h.returns) }}
                    </div>
                    <div class="holdings-list-sub" :class="h.returnsPct >= 0 ? 'text-green' : 'text-red'">
                      {{ h.returnsPct >= 0 ? '+' : '' }}{{ formatNumber(h.returnsPct) }}%
                    </div>
                  </div>
                  <div class="holdings-list-col holdings-list-current text-right">
                    <div class="holdings-list-val">₹{{ formatNumber(h.currentValue) }}</div>
                    <div class="holdings-list-sub">₹{{ formatNumber(h.invested) }}</div>
                  </div>
                </div>
                <div v-if="!portfolioStats.holdingsList.length" class="empty-state-td">
                  No active holdings.
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- ORDERS TAB -->
        <template v-if="activeTab === 'orders'">
          <div class="orders-section">
            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>Time</th>
                    <th>Type</th>
                    <th>Company</th>
                    <th class="text-right">Qty</th>
                    <th class="text-right">Price</th>
                    <th class="text-right">Status</th>
                    <th class="text-right">Action</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="t in normalizedHistory" :key="t.id">
                    <td>{{ formatDate(t.created_at) }}</td>
                    <td>
                      <span :class="['type-badge', t.type === 'buy' ? 'badge-blue' : 'badge-red']">
                        {{ t.type.toUpperCase() }}
                      </span>
                    </td>
                    <td>
                      <div>{{ t.name }}</div>
                     
                    </td>
                    <td class="text-right">{{ t.quantity }}</td>
                    <td class="text-right">₹{{ formatNumber(t.price) }}</td>
                    <td class="text-right">
                      <span :class="['status-badge', `status-${t.status}`]">{{ t.status }}</span>
                    </td>
                    <td class="text-right">
                      <button v-if="t.status === 'pending'" class="btn-text-red" @click="handleCancelTrade(t.id)">Cancel</button>
                      <span v-else class="muted-text">-</span>
                    </td>
                  </tr>
                  <tr v-if="!normalizedHistory.length">
                    <td colspan="7" class="text-center empty-state-td">No orders placed yet.</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </template>
        
        <!-- WATCHLIST TAB -->
        <template v-if="activeTab === 'watchlist'">
          <Watchlist />
        </template>
        
      </div>
    </div>

    <!-- Trade Side Panel Overlay & Panel -->
    <div class="panel-overlay" :class="{ 'is-open': showTradePanel }" @click="closeTradePanel"></div>
    <div class="trade-side-panel" :class="{ 'is-open': showTradePanel }">
      <div class="panel-header" v-if="activeStock">
        <div>
          <h2>{{ activeStock.name }}</h2>
          <p>{{ activeStock.symbol }} &bull; ₹{{ formatNumber(activeStock.currentPrice) }}</p>
        </div>
        <button class="panel-close-btn" @click="closeTradePanel"><i class="fas fa-times"></i></button>
      </div>

      <div class="panel-body" v-if="activeStock">
        <div class="panel-tabs">
          <button class="panel-tab buy" :class="{ active: tradeType === 'buy' }" @click="tradeType = 'buy'">BUY</button>
          <button class="panel-tab sell" :class="{ active: tradeType === 'sell' }" @click="tradeType = 'sell'">SELL</button>
        </div>

        <div class="trade-form-group">
          <label>Qty</label>
          <input type="number" class="trade-input" v-model.number="tradeQty" min="1" />
        </div>

        <div class="trade-form-group">
          <label>Price</label>
          <input type="text" class="trade-input" value="At market" readonly />
        </div>

        <div class="trade-info-box">
          <div class="trade-info-row">
            <span class="text-muted">Balance</span>
            <span>₹{{ formatNumber(walletStore.state.balance?.available_balance) }}</span>
          </div>
          <div class="trade-info-row">
            <span class="text-muted">Approx req.</span>
            <span>₹{{ formatNumber(estimatedTradeValue) }}</span>
          </div>
          <div class="trade-info-row" v-if="tradeType === 'sell'">
            <span class="text-muted">Available Qty</span>
            <span>{{ activeStock.qty }}</span>
          </div>
        </div>
      </div>

      <div class="panel-footer" v-if="activeStock">
        <button class="btn-submit-trade" :class="tradeType" @click="executeTrade">
          {{ tradeType.toUpperCase() }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import '../../assets/css/dashboard.css'
import { onMounted, computed, watch, nextTick, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Chart from 'chart.js/auto'
import { useWalletStore } from '../../store/wallet'
import { useTradeStore } from '../../store/trade'
import { useAuthStore } from '../../store'
import { useToast } from 'vue-toastification'
import { formatNumber, formatDate } from '../../utils/formatters'
import Watchlist from '../../views/Watchlist.vue'
import api from '../../services/api'

const route = useRoute()
const router = useRouter()
const walletStore = useWalletStore()

const goToHoldings = () => {
  router.push({ query: { tab: 'holdings' } })
}
const tradeStore = useTradeStore()
const authStore = useAuthStore()
const toast = useToast()

let chartInstance = null

const activeTab = computed(() => route.query.tab || 'explore')
const timeFrame = ref('ALL')
const marketIndices = ref([
  { name: 'NIFTY', current: 0, change: 0, change_pct: 0 },
  { name: 'SENSEX', current: 0, change: 0, change_pct: 0 },
  { name: 'BANKNIFTY', current: 0, change: 0, change_pct: 0 },
  { name: 'MIDCPNIFTY', current: 0, change: 0, change_pct: 0 },
  { name: 'FINNIFTY', current: 0, change: 0, change_pct: 0 }
])

const fetchMarketIndices = async () => {
  try {
    const res = await api.get('/market/indices')
    if (res.data && res.data.data) {
      marketIndices.value = res.data.data
    }
  } catch (error) {
    console.error('Failed to fetch market indices:', error)
  }
}

const showValues = ref(true)
const showAnalyse = ref(false)
let analyseChartInstance = null

const toggleValues = () => { showValues.value = !showValues.value }
const toggleAnalyse = () => {
  showAnalyse.value = !showAnalyse.value
  if(showAnalyse.value) {
    nextTick(() => { buildAnalyseChart() })
  }
}

const buildAnalyseChart = () => {
  const ctx = document.getElementById('analyseChart')
  if (!ctx) return
  if (analyseChartInstance) analyseChartInstance.destroy()

  const labels = []
  const dataInvested = []
  const dataCurrent = []

  const history = tradeStore.state.history || []
  const sortedHistory = [...history]
    .sort((a, b) => new Date(a.CreatedAt || a.created_at) - new Date(b.CreatedAt || b.created_at))
    .filter(t => (t.Status || t.status || '').toLowerCase() === 'completed')
  
  let runningInv = 0
  
  sortedHistory.forEach(t => {
    const type = (t.Type || t.type || '').toLowerCase()
    const qty = t.Quantity || t.quantity || 0
    const price = t.Price || t.price || 0
    
    if (type === 'buy') runningInv += (qty * price)
    else if (type === 'sell') runningInv -= (qty * price)
    if (runningInv < 0) runningInv = 0
    
    labels.push(formatDate(t.CreatedAt || t.created_at))
    dataInvested.push(parseFloat(runningInv.toFixed(2)))
    dataCurrent.push(0) // placeholder
  })

  // Add final point
  labels.push('Now')
  dataInvested.push(portfolioStats.value.invested)
  dataCurrent.push(portfolioStats.value.currentVal)

  // Interpolate Current Value divergence so the chart looks realistic
  const finalPL = portfolioStats.value.currentVal - portfolioStats.value.invested
  const totalPoints = dataCurrent.length
  for (let i = 0; i < totalPoints - 1; i++) {
    const divergence = finalPL * ((i + 1) / totalPoints)
    dataCurrent[i] = parseFloat((dataInvested[i] + divergence).toFixed(2))
  }

  const isProfit = portfolioStats.value.totalRet >= 0
  const colorHex = isProfit ? '#10b981' : '#f43f5e'
  const fillRgba = isProfit ? 'rgba(16, 185, 129, 0.1)' : 'rgba(244, 63, 94, 0.1)'

  analyseChartInstance = new Chart(ctx, {
    type: 'line',
    data: {
      labels,
      datasets: [
        {
          label: 'Current Value',
          data: dataCurrent,
          borderColor: colorHex,
          backgroundColor: fillRgba,
          borderWidth: 2,
          fill: true,
          tension: 0.4,
          pointRadius: 3,
          pointBackgroundColor: colorHex,
        },
        {
          label: 'Invested',
          data: dataInvested,
          borderColor: '#3b82f6',
          borderDash: [5, 5],
          borderWidth: 2,
          fill: false,
          tension: 0.4,
          pointRadius: 0,
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: true, labels: { color: '#e0e0e0', usePointStyle: true, boxWidth: 8 } },
        tooltip: {
          backgroundColor: '#121212',
          titleColor: '#9b9b9b',
          bodyColor: '#e0e0e0',
          borderColor: 'rgba(255,255,255,0.08)',
          borderWidth: 1,
          padding: 12,
          callbacks: {
            label: function(context) { return context.dataset.label + ': ₹' + Number(context.raw).toLocaleString(); }
          }
        }
      },
      scales: {
        x: { display: true, ticks: { color: '#9b9b9b', maxTicksLimit: 5 }, grid: { display: false } },
        y: { display: true, ticks: { color: '#9b9b9b', callback: val => '₹' + val.toLocaleString() }, grid: { color: 'rgba(255,255,255,0.05)' } }
      },
      interaction: {
        mode: 'index',
        intersect: false,
      }
    }
  })
}

const setTimeFrame = (tf) => {
  timeFrame.value = tf
  buildChart()
}

// Trade Panel State
const showTradePanel = ref(false)
const activeStock = ref(null)
const tradeType = ref('buy')
const tradeQty = ref(1)

const estimatedTradeValue = computed(() => {
  if (!activeStock.value) return 0
  return activeStock.value.currentPrice * (tradeQty.value || 0)
})

const openTradePanel = (holding) => {
  activeStock.value = holding
  tradeType.value = 'buy'
  tradeQty.value = 1
  showTradePanel.value = true
}

const closeTradePanel = () => {
  showTradePanel.value = false
  activeStock.value = null
}

const handleCancelTrade = async (tradeId) => {
  if (confirm("Are you sure you want to cancel this order?")) {
    const res = await tradeStore.cancelTrade(tradeId)
    if (res.success) {
      toast.success(res.message)
      await walletStore.fetchBalance()
    } else {
      toast.error(res.message)
    }
  }
}

const executeTrade = async () => {
  if (!activeStock.value || tradeQty.value <= 0) return
  
  if (tradeType.value === 'sell' && tradeQty.value > activeStock.value.qty) {
    toast.error(`You only have ${activeStock.value.qty} shares to sell.`)
    return
  }

  try {
    const req = {
      share_id: activeStock.value.id,
      quantity: tradeQty.value
    }

    if (tradeType.value === 'buy') {
      const res = await tradeStore.buyShare(activeStock.value.id, tradeQty.value)
      if (!res.success) throw new Error(res.message)
      toast.success(res.message || `Successfully bought ${tradeQty.value} shares of ${activeStock.value.symbol}`)
    } else {
      const res = await tradeStore.sellShare(activeStock.value.id, tradeQty.value)
      if (!res.success) throw new Error(res.message)
      toast.success(res.message || `Successfully sold ${tradeQty.value} shares of ${activeStock.value.symbol}`)
    }

    closeTradePanel()
    // Refresh data
    await Promise.all([
      walletStore.fetchBalance(),
      walletStore.fetchTransactions(),
      tradeStore.fetchHistory()
    ])
  } catch (err) {
    toast.error(err.response?.data?.message || err.message || 'Trade failed')
  }
}


const portfolioStats = computed(() => {
  const history = tradeStore.state.history || []
  const shares = tradeStore.state.shares || []
  
  const holdings = {}
  
  // Sort history chronologically
  const sortedHistory = [...history].sort((a, b) => new Date(a.CreatedAt || a.created_at) - new Date(b.CreatedAt || b.created_at))
  
  const todayStr = new Date().toDateString()
  
  sortedHistory.forEach(t => {
    const type = (t.Type || t.type || '').toLowerCase()
    const status = (t.Status || t.status || '').toLowerCase()
    const qty = t.Quantity || t.quantity || 0
    const price = t.Price || t.price || 0
    const shareId = t.ShareID || t.share_id || (t.Share ? (t.Share.ID || t.Share.id) : null)
    const isToday = new Date(t.CreatedAt || t.created_at).toDateString() === todayStr
    
    if (status === 'completed' && shareId) {
      if (!holdings[shareId]) {
        holdings[shareId] = { qty: 0, totalCost: 0, qtyToday: 0, costToday: 0 }
      }
      
      const holding = holdings[shareId]
      if (type === 'buy') {
        holding.qty += qty
        holding.totalCost += (qty * price)
        if (isToday) {
           holding.qtyToday += qty
           holding.costToday += (qty * price)
        }
      } else if (type === 'sell') {
        if (holding.qty > 0) {
          const avgPrice = holding.totalCost / holding.qty
          holding.totalCost -= (qty * avgPrice)
          holding.qty -= qty
          
          if (holding.qty < holding.qtyToday) {
             if (holding.qtyToday > 0) {
                 const avgToday = holding.costToday / holding.qtyToday
                 holding.qtyToday = holding.qty
                 holding.costToday = holding.qtyToday * avgToday
             }
          }
          
          if (holding.qty <= 0) {
            holding.qty = 0
            holding.totalCost = 0
            holding.qtyToday = 0
            holding.costToday = 0
          }
        }
      }
    }
  })
  
  let currentVal = 0
  let totalInvested = 0
  let totalOneDayRet = 0
  
  const holdingsList = Object.entries(holdings).filter(([id, h]) => h.qty > 0).map(([id, h]) => {
     const share = shares.find(s => (s.ID || s.id) === id)
     const currentPrice = parseFloat((share ? (share.Price || share.price) : (h.totalCost / h.qty)).toFixed(2))
     const prevPrice = parseFloat((share ? (share.PreviousPrice || share.previous_price || share.Price || share.price) : (h.totalCost / h.qty)).toFixed(2))
     const avgPrice = parseFloat((h.totalCost / h.qty).toFixed(2))
     const currentValue = parseFloat((h.qty * currentPrice).toFixed(2))
     const invested = parseFloat(h.totalCost.toFixed(2))
     const returns = parseFloat((currentValue - invested).toFixed(2))
     const returnsPct = invested > 0 ? parseFloat(((returns / invested) * 100).toFixed(2)) : 0
     
     const qtyBeforeToday = h.qty - (h.qtyToday || 0)
     const oneDayRetBeforeToday = (currentPrice - prevPrice) * qtyBeforeToday
     const oneDayRetToday = (currentPrice * (h.qtyToday || 0)) - (h.costToday || 0)
     const oneDayRetForHolding = oneDayRetBeforeToday + oneDayRetToday
     const oneDayRet = parseFloat(oneDayRetForHolding.toFixed(2))
     const prevVal = currentValue - oneDayRet
     const oneDayRetPct = prevVal > 0 ? parseFloat(((oneDayRet / prevVal) * 100).toFixed(2)) : 0
     
     currentVal += currentValue
     totalInvested += invested
     totalOneDayRet += oneDayRetForHolding
     
     return {
       id,
       symbol: share?.Symbol || share?.symbol || 'Unknown',
       name: share?.Name || share?.name || 'Unknown',
       qty: h.qty,
       invested,
       avgPrice,
       currentPrice,
       currentValue,
       returns,
       returnsPct,
       oneDayRet,
       oneDayRetPct
     }
  })
  
  const totalRet = parseFloat((currentVal - totalInvested).toFixed(2))
  const totalRetPct = totalInvested > 0 ? parseFloat(((totalRet / totalInvested) * 100).toFixed(2)) : 0
  
  const prevTotalVal = currentVal - totalOneDayRet
  const oneDayRet = parseFloat(totalOneDayRet.toFixed(2))
  const oneDayRetPct = prevTotalVal > 0 ? parseFloat(((totalOneDayRet / prevTotalVal) * 100).toFixed(2)) : 0
  
  return {
    invested: totalInvested,
    currentVal: parseFloat(currentVal.toFixed(2)),
    totalRet,
    totalRetPct,
    oneDayRet,
    oneDayRetPct,
    holdingsList
  }
})

const normalizedHistory = computed(() => {
  const history = tradeStore.state.history || []
  return history.map(t => ({
    id: t.ID || t.id,
    type: (t.Type || t.type || '').toLowerCase(),
    status: (t.Status || t.status || '').toLowerCase(),
    quantity: t.Quantity || t.quantity,
    price: t.Price || t.price,
    created_at: t.CreatedAt || t.created_at,
    symbol: t.Share?.Symbol || t.Share?.symbol || t.share?.Symbol || t.share?.symbol || 'Unknown',
    name: t.Share?.Name || t.Share?.name || t.share?.Name || t.share?.name || 'Unknown'
  })).sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
})

const normalizedTransactions = computed(() => {
  if (!walletStore.state.transactions) return []
  return walletStore.state.transactions.map(tx => ({
    id: tx.ID || tx.id,
    type: (tx.Type || tx.type || '').toLowerCase(),
    description: tx.Description || tx.description || '',
    amount: tx.Amount || tx.amount || 0,
    created_at: tx.CreatedAt || tx.created_at
  }))
})

const recentTransactions = computed(() => {
  return normalizedTransactions.value.slice(0, 5)
})

const buildChart = () => {
  const ctx = document.getElementById('portfolioChart')
  if (!ctx) return
  if (chartInstance) chartInstance.destroy()

  const labels = []
  const dataInvested = []
  const dataCurrent = []

  const history = tradeStore.state.history || []
  const sortedHistory = [...history]
    .sort((a, b) => new Date(a.CreatedAt || a.created_at) - new Date(b.CreatedAt || b.created_at))
    .filter(t => (t.Status || t.status || '').toLowerCase() === 'completed')
  
  let runningInv = 0
  
  let cutoff = null
  if (timeFrame.value !== 'ALL') {
    const now = new Date()
    cutoff = new Date()
    if (timeFrame.value === '1D') cutoff.setDate(now.getDate() - 1)
    else if (timeFrame.value === '1W') cutoff.setDate(now.getDate() - 7)
    else if (timeFrame.value === '1M') cutoff.setMonth(now.getMonth() - 1)
    else if (timeFrame.value === '1Y') cutoff.setFullYear(now.getFullYear() - 1)
  }

  sortedHistory.forEach(t => {
    const type = (t.Type || t.type || '').toLowerCase()
    const qty = t.Quantity || t.quantity || 0
    const price = t.Price || t.price || 0
    
    if (type === 'buy') runningInv += (qty * price)
    else if (type === 'sell') runningInv -= (qty * price)
    if (runningInv < 0) runningInv = 0
    
    if (!cutoff || new Date(t.CreatedAt || t.created_at) >= cutoff) {
      labels.push(formatDate(t.CreatedAt || t.created_at))
      dataInvested.push(parseFloat(runningInv.toFixed(2)))
      dataCurrent.push(0) // placeholder
    }
  })

  // Add final point
  labels.push('Now')
  dataInvested.push(portfolioStats.value.invested)
  dataCurrent.push(portfolioStats.value.currentVal)

  // Interpolate Current Value divergence so the chart looks realistic
  const finalPL = portfolioStats.value.currentVal - portfolioStats.value.invested
  const totalPoints = dataCurrent.length
  for (let i = 0; i < totalPoints - 1; i++) {
    const divergence = finalPL * ((i + 1) / totalPoints)
    dataCurrent[i] = parseFloat((dataInvested[i] + divergence).toFixed(2))
  }

  const isProfit = portfolioStats.value.totalRet >= 0
  const colorHex = isProfit ? '#10b981' : '#f43f5e'
  const fillRgba = isProfit ? 'rgba(16, 185, 129, 0.1)' : 'rgba(244, 63, 94, 0.1)'

  chartInstance = new Chart(ctx, {
    type: 'line',
    data: {
      labels,
      datasets: [
        {
          label: 'Current Value (₹)',
          data: dataCurrent,
          borderColor: colorHex,
          backgroundColor: fillRgba,
          borderWidth: 2,
          fill: true,
          tension: 0.4,
          pointRadius: 3,
          pointBackgroundColor: colorHex,
          pointBorderColor: '#131929',
          pointBorderWidth: 1,
        },
        {
          label: 'Invested (₹)',
          data: dataInvested,
          borderColor: '#3b82f6',
          borderDash: [5, 5],
          borderWidth: 2,
          fill: false,
          tension: 0.4,
          pointRadius: 0,
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: true, labels: { color: '#e0e0e0', usePointStyle: true, boxWidth: 8 } },
        tooltip: {
          backgroundColor: '#1a2236',
          titleColor: '#f1f5f9',
          bodyColor: '#64748b',
          borderColor: 'rgba(255,255,255,0.07)',
          borderWidth: 1,
          padding: 12,
          callbacks: {
            title: function(tooltipItems) {
              return 'Time: ' + tooltipItems[0].label;
            },
            label: function(context) {
              if (context.datasetIndex === 1) return null;
              const currentVal = context.parsed.y;
              const investedVal = context.chart.data.datasets[1].data[context.dataIndex];
              const returns = currentVal - investedVal;
              const sign = returns >= 0 ? '+' : '';
              return [
                'Price: ₹' + currentVal.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2}),
                'Total Returns: ' + sign + '₹' + returns.toLocaleString(undefined, {minimumFractionDigits: 2, maximumFractionDigits: 2})
              ];
            }
          }
        }
      },
      scales: {
        y: {
          grid: { color: 'rgba(255,255,255,0.04)' },
          ticks: { color: '#64748b', callback: val => '₹' + val.toLocaleString() }
        },
        x: {
          grid: { display: false },
          ticks: { color: '#64748b', maxTicksLimit: 6 }
        }
      }
    }
  })
}

onMounted(async () => {
  await Promise.all([
    walletStore.fetchBalance(),
    walletStore.fetchTransactions(),
    tradeStore.fetchShares(),
    tradeStore.fetchHistory()
  ])
  fetchMarketIndices()
  setInterval(fetchMarketIndices, 60000)
  if (activeTab.value === 'explore') {
    buildChart()
  }
})

watch(() => walletStore.state.transactions, () => {
  if (activeTab.value === 'explore') {
    buildChart()
  }
}, { deep: true })

watch(activeTab, async (newTab) => {
  if (newTab === 'explore') {
    await nextTick()
    buildChart()
  }
})
</script>

<style scoped>
.chart-toggles {
  display: flex;
  gap: 4px;
  background: rgba(255, 255, 255, 0.05);
  padding: 4px;
  border-radius: 8px;
}

.market-indices-banner {
  background: rgba(15, 23, 42, 0.7);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding: 10px 0;
  overflow: hidden;
  white-space: nowrap;
  position: relative;
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.marquee-content {
  display: inline-flex;
  gap: 30px;
  animation: marquee 25s linear infinite;
}

.market-indices-banner:hover .marquee-content {
  animation-play-state: paused;
}

.index-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
}

.index-name {
  color: #94a3b8;
}

.index-val {
  color: #f1f5f9;
  font-weight: 600;
}

.index-change {
  font-size: 12px;
}

@keyframes marquee {
  0% { transform: translateX(100vw); }
  100% { transform: translateX(-100%); }
}
.chart-toggles {
  display: flex;
  gap: 4px;
  background: rgba(255, 255, 255, 0.05);
  padding: 4px;
  border-radius: 8px;
}
.chart-toggle {
  background: transparent;
  border: none;
  color: #64748b;
  font-size: 11px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}
.chart-toggle:hover {
  color: #f1f5f9;
}
.chart-toggle.active {
  background: #3b82f6;
  color: #ffffff;
}

.analyse-section {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px dashed rgba(255,255,255,0.1);
}
.holdings-summary-btn.active {
  background-color: rgba(59, 130, 246, 0.2);
  border-color: #3b82f6;
  color: #3b82f6;
}
</style>

