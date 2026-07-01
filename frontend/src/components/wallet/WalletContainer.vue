<template>
  <div class="walletpage">
    <div class="header-section">
      <h1 class="page-title">My Wallet</h1>
      <p class="page-subtitle">Manage your funds and view ledger history</p>
    </div>

    <div class="wallet-grid">
      <!-- Balance Card -->
      <div class="balance-card">
        <div class="card-header">
          <h2 class="card-title">Available Balance</h2>
          <i class="fas fa-wallet bg-icon"></i>
        </div>
        <div class="balance-amount" v-if="!walletStore.state.loading">₹{{ walletStore.state.balance?.available_balance?.toFixed(2) || '0.00' }}</div>
        <div class="balance-amount" v-else>Loading...</div>
        <div class="balance-meta">
          <span>Blocked for trades: ₹{{ walletStore.state.balance?.blocked_balance?.toFixed(2) || '0.00' }}</span>
          <span>Total: ₹{{ walletStore.state.balance?.wallet_balance?.toFixed(2) || '0.00' }}</span>
        </div>
        
        <div class="wallet-actions">
          <button 
            class="action-btn" 
            :class="actionType === 'deposit' ? 'btn-primary' : 'btn-secondary'" 
            @click="actionType = 'deposit'"
          >
            <i class="fas fa-plus fa-lg" style="margin-right: 4px;"></i> Add Funds
          </button>
          <button 
            class="action-btn" 
            :class="actionType === 'withdraw' ? 'btn-primary' : 'btn-secondary'" 
            @click="actionType = 'withdraw'"
          >
            <i class="fas fa-arrow-down fa-lg" style="margin-right: 4px;"></i> Withdraw
          </button>
        </div>
      </div>

      <!-- Quick Action Form -->
      <div class="action-form-card">
        <h2 class="form-title">{{ actionType === 'deposit' ? 'Quick Deposit' : 'Withdraw Funds' }}</h2>
        <form @submit.prevent="initiateTransaction" class="deposit-form">
          <div class="amount-presets">
            <button type="button" class="preset-btn" @click="amount = 100">₹100</button>
            <button type="button" class="preset-btn" @click="amount = 500">₹500</button>
            <button type="button" class="preset-btn" @click="amount = 1000">₹1,000</button>
          </div>
          
          <div class="form-group">
            <label class="form-label">Custom Amount</label>
            <div class="input-wrapper">
              <i class="fas fa-rupee-sign input-icon" style="font-size: 16px;"></i>
              <input type="number" class="form-input" v-model="amount" :min="actionType === 'deposit' ? 100 : 1" placeholder="Enter amount" />
            </div>
          </div>
          
          <button 
            type="submit" 
            class="submit-btn" 
            :class="actionType === 'deposit' ? 'submit-deposit' : 'submit-withdraw'" 
            :disabled="!amount || walletStore.state.loading"
          >
            <i v-if="walletStore.state.loading" class="fas fa-circle-notch fa-spin"></i>
            <span v-else>{{ actionType === 'deposit' ? 'Proceed to Payment' : 'Confirm Withdrawal' }}</span>
          </button>
        </form>
      </div>
    </div>

    <!-- Payment Confirmation Modal (Groww/Zerodha style flow) -->
    <div class="modal-overlay" v-if="showBankModal" @click.self="showBankModal = false">
      <div class="payment-modal">
        <div class="modal-header">
          <h3>Select Bank Account</h3>
          <button class="close-btn" @click="showBankModal = false"><i class="fas fa-times"></i></button>
        </div>
        
        <div class="modal-body" v-if="!isProcessingPayment">
          <p class="payment-amount">{{ actionType === 'deposit' ? 'Amount to add:' : 'Amount to withdraw:' }} <strong>₹{{ amount }}</strong></p>
          <p class="modal-subtitle">{{ actionType === 'deposit' ? 'Pay securely using your verified bank accounts.' : 'Receive funds securely into your verified bank account.' }}</p>

          <div class="bank-list" v-if="userBanks.length > 0">
            <label 
              v-for="(bank, index) in userBanks" 
              :key="index" 
              class="bank-option" 
              :class="{ 'selected': selectedBank === index }"
            >
              <input type="radio" :value="index" v-model="selectedBank" name="bank_select" class="hidden-radio" />
              <div class="bank-details">
                <div class="bank-icon">
                  <i class="fas fa-university"></i>
                </div>
                <div class="bank-info">
                  <span class="bank-name">{{ bank.bank_name || 'Verified Bank' }}</span>
                  <span class="bank-acc">•••• {{ String(bank.account_number).slice(-4) || 'XXXX' }}</span>
                </div>
              </div>
              <div class="bank-type-badge">{{ bank.account_type }}</div>
              <i class="fas fa-check-circle check-icon" v-if="selectedBank === index"></i>
            </label>
          </div>
          <div class="bank-list" v-else>
            <div class="no-banks-message text-center p-4">
              <i class="fas fa-exclamation-triangle text-warning" style="font-size: 32px; margin-bottom: 12px;"></i>
              <p style="font-size: 15px; font-weight: 500; color: var(--text-main);">You have not added any bank accounts.</p>
              <p style="font-size: 13px; color: var(--text-muted); margin-top: 4px;">Please add the bank details in your profile then proceed.</p>
              <button class="btn-outline-small mt-4" style="border-color: #f59e0b; color: #f59e0b;" @click="$router.push('/profile')">
                <i class="fas fa-edit mr-2"></i> Update KYC
              </button>
            </div>
          </div>

          <button class="pay-now-btn mt-4" @click="processPayment" :disabled="selectedBank === null || userBanks.length === 0">
            {{ actionType === 'deposit' ? 'Pay' : 'Withdraw' }} ₹{{ amount }}
          </button>
        </div>

        <div class="modal-body processing-state" v-else>
          <div class="spinner-wrapper">
            <div class="loading-spinner"></div>
          </div>
          <h3>{{ actionType === 'deposit' ? 'Processing Payment' : 'Processing Withdrawal' }}</h3>
          <p>Please do not close this window or press back.</p>
          <p class="small-text">Waiting for bank confirmation...</p>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import '../../assets/css/wallet.css'
import { ref, onMounted, computed } from 'vue'
import { useToast } from 'vue-toastification'
import { useWalletStore } from '../../store/wallet'
import { useAuthStore } from '../../store'

const toast = useToast()
const amount = ref(null)
const actionType = ref('deposit') // 'deposit' or 'withdraw'
const walletStore = useWalletStore()
const authStore = useAuthStore()

const showBankModal = ref(false)
const selectedBank = ref(0) // Default to first bank
const isProcessingPayment = ref(false)

const userBanks = computed(() => {
  return authStore.state.user?.bank_accounts || []
})

onMounted(async () => {
  await walletStore.fetchBalance()
})

const initiateTransaction = async () => {
  if (!amount.value || amount.value <= 0) {
    toast.error('Please enter a valid amount greater than 0')
    return;
  }
  
  if (actionType.value === 'deposit' && amount.value < 100) {
    toast.error('Minimum deposit amount is ₹100')
    return;
  }
  
  if (actionType.value === 'withdraw') {
    if (amount.value > (walletStore.state.balance?.available_balance || 0)) {
      toast.error('Insufficient available balance for withdrawal')
      return
    }
  }

  // Open payment/withdrawal modal
  selectedBank.value = 0 
  showBankModal.value = true
}

const processPayment = async () => {
  isProcessingPayment.value = true
  
  // Simulate network/gateway delay (2 seconds)
  await new Promise(resolve => setTimeout(resolve, 2000))
  
  const val = amount.value
  let res;
  
  if (actionType.value === 'deposit') {
    res = await walletStore.addFunds(val)
  } else {
    res = await walletStore.withdrawFunds(val)
  }
  
  isProcessingPayment.value = false
  showBankModal.value = false
  
  if (res.success) {
    const actionText = actionType.value === 'deposit' ? 'added to' : 'withdrawn from'
    toast.success(`Successfully ${actionText} wallet: ₹${val}`)
    amount.value = null
  } else {
    toast.error(res.message)
  }
}
</script>

<style scoped>
.submit-withdraw {
  background: linear-gradient(135deg, #f43f5e 0%, #e11d48 100%) !important;
  box-shadow: 0 8px 24px rgba(244, 63, 94, 0.3) !important;
}
.submit-withdraw:hover:not(:disabled) {
  filter: brightness(1.1) !important;
  box-shadow: 0 8px 30px rgba(244, 63, 94, 0.5) !important;
}

/* Modal CSS */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.payment-modal {
  background: var(--surface-card, #131929);
  border: 1px solid var(--border-color, rgba(255,255,255,0.1));
  border-radius: var(--border-radius-lg, 12px);
  width: 100%;
  max-width: 450px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.4);
  overflow: hidden;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color, rgba(255,255,255,0.1));
  background: rgba(0,0,0,0.1);
}
.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--text-main);
}
.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 20px;
  cursor: pointer;
}
.close-btn:hover {
  color: var(--text-main);
}

.modal-body {
  padding: 24px;
}

.payment-amount {
  font-size: 18px;
  color: var(--text-main);
  margin-bottom: 4px;
}
.payment-amount strong {
  color: var(--primary-color);
  font-size: 22px;
}
.modal-subtitle {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 24px;
}

.bank-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.bank-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: rgba(255,255,255,0.03);
  border: 1px solid var(--border-color, rgba(255,255,255,0.1));
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}
.bank-option:hover {
  background: rgba(255,255,255,0.05);
}
.bank-option.selected {
  border-color: var(--primary-color);
  background: rgba(16, 185, 129, 0.05);
}

.hidden-radio {
  display: none;
}

.bank-details {
  display: flex;
  align-items: center;
  gap: 16px;
}
.bank-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(255,255,255,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-main);
  font-size: 18px;
}
.bank-info {
  display: flex;
  flex-direction: column;
}
.bank-name {
  font-weight: 600;
  font-size: 15px;
  color: var(--text-main);
}
.bank-acc {
  font-size: 13px;
  color: var(--text-muted);
}
.bank-type-badge {
  font-size: 10px;
  text-transform: uppercase;
  background: rgba(255,255,255,0.1);
  padding: 2px 6px;
  border-radius: 4px;
  color: var(--text-muted);
  position: absolute;
  top: 16px;
  right: 45px;
}

.check-icon {
  color: var(--primary-color);
  font-size: 20px;
}

.pay-now-btn {
  width: 100%;
  padding: 14px;
  background: var(--primary-color);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.pay-now-btn:hover:not(:disabled) {
  background: var(--primary-dim);
}
.pay-now-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.processing-state {
  text-align: center;
  padding: 40px 24px;
}
.processing-state h3 {
  margin-top: 24px;
  margin-bottom: 8px;
  color: var(--text-main);
}
.processing-state p {
  color: var(--text-muted);
  font-size: 14px;
}
.small-text {
  font-size: 12px !important;
  opacity: 0.7;
  margin-top: 16px;
}

.spinner-wrapper {
  display: flex;
  justify-content: center;
}
.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255,255,255,0.1);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
.mt-4 { margin-top: 24px; }
.mt-2 { margin-top: 8px; }
.text-warning { color: #f59e0b; }
.help-text { font-size: 12px; }
</style>
