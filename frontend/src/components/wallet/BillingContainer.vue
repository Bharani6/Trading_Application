<template>
  <div class="billingpage">
    <div class="header-section">
      <h1 class="page-title">Billing & History</h1>
      <button class="export-btn" @click="downloadCSV"><i class="fas fa-download"></i> Export CSV</button>
    </div>

    <div class="ledger-section">
      <div class="ledger-filters">
        <select class="filter-select" v-model="filterType">
          <option>All Transactions</option>
          <option>Deposits</option>
          <option>Withdrawals</option>
          <option>Trades</option>
        </select>
        <input type="date" class="date-picker" v-model="filterDate" />
      </div>

      <table class="ledger-table">
        <thead>
          <tr>
            <th>Date</th>
            <th>Reference ID</th>
            <th>Type</th>
            <th>Description</th>
            <th>Status</th>
            <th class="amount-col">Amount</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="tx in filteredLedger" :key="tx.id || tx.ID">
            <td>{{ new Date(tx.CreatedAt || tx.created_at).toLocaleString() }}</td>
            <td class="ref-id">{{ tx.ReferenceID || tx.reference_id || 'N/A' }}</td>
            <td><span :class="['type-badge', `type-${(tx.Type || tx.type || '').toLowerCase()}`]">{{ (tx.Type || tx.type || '').toUpperCase() }}</span></td>
            <td>{{ tx.Description || tx.description }}</td>
            <td><span class="status-badge">{{ (tx.Status || tx.status || '').toUpperCase() }}</span></td>
            <td :class="['amount-col', (tx.Amount || tx.amount || 0) > 0 ? 'amount-positive' : 'amount-negative']">
              {{ (tx.Amount || tx.amount || 0) > 0 ? '+' : '' }}₹{{ Math.abs(tx.Amount || tx.amount || 0).toFixed(2) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import '../../assets/css/billing.css'
import { ref, computed, onMounted } from 'vue'
import { useWalletStore } from '../../store/wallet'
import { useToast } from 'vue-toastification'

const walletStore = useWalletStore()
const toast = useToast()

const filterType = ref('All Transactions')
const filterDate = ref('')

const filteredLedger = computed(() => {
  if (!walletStore.state.transactions) return []
  return walletStore.state.transactions.filter(tx => {
    const type = (tx.Type || tx.type || '').toLowerCase()
    const dateStr = (tx.CreatedAt || tx.created_at)

    // Filter by type
    let typeMatches = true
    if (filterType.value === 'Deposits') typeMatches = type.includes('add_fund') || type.includes('deposit')
    else if (filterType.value === 'Withdrawals') typeMatches = type.includes('withdraw')
    else if (filterType.value === 'Trades') typeMatches = type.includes('trade')
    
    // Filter by date
    let dateMatches = true
    if (filterDate.value && dateStr) {
      const txDate = new Date(dateStr).toISOString().split('T')[0]
      dateMatches = txDate === filterDate.value
    }
    
    return typeMatches && dateMatches
  })
})

const downloadCSV = () => {
  if (!filteredLedger.value || filteredLedger.value.length === 0) {
    // Requires toast to be imported. We don't have toast in Billing right now. Let's import it.
    toast.error('No records found to export')
    return
  }
  
  const headers = ['Date', 'Reference ID', 'Type', 'Description', 'Status', 'Amount']
  const rows = filteredLedger.value.map(tx => {
    return [
      new Date(tx.CreatedAt || tx.created_at).toLocaleString(),
      tx.ReferenceID || tx.reference_id || 'N/A',
      (tx.Type || tx.type || '').toUpperCase(),
      tx.Description || tx.description || '',
      (tx.Status || tx.status || '').toUpperCase(),
      (tx.Amount || tx.amount || 0).toFixed(2)
    ].map(v => `"${v}"`).join(',') // wrap in quotes to prevent comma breaks
  })
  
  const csvContent = "data:text/csv;charset=utf-8," + [headers.join(','), ...rows].join('\n')
  const encodedUri = encodeURI(csvContent)
  const link = document.createElement("a")
  link.setAttribute("href", encodedUri)
  link.setAttribute("download", "transactions_export.csv")
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

onMounted(async () => {
  await walletStore.fetchTransactions()
})
</script>

