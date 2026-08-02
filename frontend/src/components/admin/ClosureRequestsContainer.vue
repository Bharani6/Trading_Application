<template>
  <div class="usermanagementpage">
    <div class="dashboard-header">
      <h1>Closure Requests</h1>
      <p>Manage account closure requests</p>
    </div>
    
    <div class="table-card">
      <div class="table-header">
        <div class="search-box">
          <i class="fas fa-search"></i>
          <input type="text" v-model="searchQuery" placeholder="Search by name, email, or PAN..." />
        </div>
      </div>

      <div class="table-responsive">
        <table>
          <thead>
            <tr>
              <th>User</th>
              <th>KYC Docs</th>
              <th>Joined Date</th>
              <th>Closure Request Date</th>
              <th>Wallet Balance</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="7" class="text-center py-8">
                <i class="fas fa-spinner fa-spin text-primary text-2xl"></i>
              </td>
            </tr>
            <tr v-else-if="filteredUsers.length === 0">
              <td colspan="7" class="text-center py-8 text-muted">
                No closure requests found.
              </td>
            </tr>
            <tr v-for="user in filteredUsers" :key="user.id" v-else>
              <td>
                <div class="user-cell">
                  <div class="user-avatar" :class="{'has-photo': user.ipv_photo}">
                    <img v-if="user.ipv_photo" :src="user.ipv_photo" alt="avatar" class="avatar-img" />
                    <span v-else>{{ user.name ? user.name.charAt(0).toUpperCase() : 'U' }}</span>
                  </div>
                  <div class="user-details">
                    <span class="user-name">{{ user.name }}</span>
                    <span class="user-email">{{ user.email }}</span>
                  </div>
                </div>
              </td>
              <td class="docs-cell">
                <div>PAN: {{ user.pan || 'N/A' }}</div>
                <div style="color: var(--text-muted);">AADHAAR: {{ user.aadhaar ? '****' + user.aadhaar.slice(-4) : 'N/A' }}</div>
              </td>
              <td>{{ user.created_at || 'N/A' }}</td>
              <td>{{ user.updated_at ? user.updated_at.split(' ')[0] : 'N/A' }}</td>
              <td class="font-medium">₹{{ (user.wallet_balance || 0).toFixed(2) }}</td>
              <td>
                <span :class="['status-badge', getStatusClass(user.status)]">
                  {{ formatStatus(user.status) }}
                </span>
              </td>
              <td class="actions-cell">
                <button 
                  class="btn-action btn-view"
                  @click="openDetailsModal(user.id)"
                >
                  <i class="fas fa-eye"></i> View
                </button>
                <button 
                  class="btn-action btn-approve"
                  @click="updateStatus(user.id, 'close_account')"
                  :disabled="actionLoading === user.id"
                >
                  <i class="fas fa-check"></i> Approve
                </button>
                <button 
                  class="btn-action btn-reject"
                  @click="updateStatus(user.id, 'reject_closure')"
                  :disabled="actionLoading === user.id"
                >
                  <i class="fas fa-times"></i> Reject
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <!-- User Details Modal -->
  <div v-if="showModal" class="modal-overlay" @click.self="closeDetailsModal">
    <div class="modal-content">
      <div class="modal-header">
        <h2>User KYC Details</h2>
        <button class="btn-close" @click="closeDetailsModal"><i class="fas fa-times"></i></button>
      </div>
      
      <div v-if="detailsLoading" class="modal-body loading-body">
        <i class="fas fa-spinner fa-spin fa-2x"></i>
        <p>Loading details...</p>
      </div>
      <div v-else-if="selectedUser" class="modal-body">
        <!-- Personal Info -->
        <section class="detail-section">
          <h3><i class="fas fa-user"></i> Personal Information</h3>
          <div class="detail-grid">
            <div class="detail-item"><strong>Name:</strong> {{ selectedUser.name }}</div>
            <div class="detail-item"><strong>Email:</strong> {{ selectedUser.email }}</div>
            <div class="detail-item"><strong>Mobile:</strong> {{ selectedUser.mobile }}</div>
            <div class="detail-item"><strong>DOB:</strong> {{ selectedUser.dob }}</div>
            <div class="detail-item"><strong>PAN:</strong> {{ selectedUser.pan }}</div>
            <div class="detail-item"><strong>Aadhaar:</strong> {{ selectedUser.aadhaar }}</div>
            <div class="detail-item"><strong>Father Name:</strong> {{ selectedUser.father_name || 'N/A' }}</div>
            <div class="detail-item"><strong>Mother Name:</strong> {{ selectedUser.mother_name || 'N/A' }}</div>
            <div class="detail-item"><strong>Income Range:</strong> {{ selectedUser.income_range }}</div>
            <div class="detail-item"><strong>Occupation:</strong> {{ selectedUser.occupation }}</div>
          </div>
          <div class="detail-item full-width mt-2"><strong>Address:</strong> {{ selectedUser.address }} <span v-if="selectedUser.city">, {{ selectedUser.city }}, {{ selectedUser.state }}, {{ selectedUser.country }} - {{ selectedUser.pincode }}</span></div>
        </section>

        <!-- Bank Details -->
        <section class="detail-section mt-4">
          <h3><i class="fas fa-university"></i> Bank Accounts</h3>
          <div v-if="selectedUser.bank_accounts && selectedUser.bank_accounts.length">
            <div v-for="(bank, index) in selectedUser.bank_accounts" :key="index" class="bank-card">
              <div class="detail-grid">
                <div class="detail-item"><strong>Bank Name:</strong> {{ bank.bank_name }}</div>
                <div class="detail-item"><strong>Account Number:</strong> {{ bank.account_number }}</div>
                <div class="detail-item"><strong>IFSC:</strong> {{ bank.ifsc }}</div>
                <div class="detail-item"><strong>Account Type:</strong> {{ bank.account_type }}</div>
                <div class="detail-item"><strong>Branch:</strong> {{ bank.branch || 'N/A' }}</div>
              </div>
            </div>
          </div>
          <div v-else class="text-muted">No bank details found.</div>
        </section>

        <!-- Nominee Details -->
        <section class="detail-section mt-4">
          <h3><i class="fas fa-users"></i> Nominees</h3>
          <div v-if="selectedUser.nominees && selectedUser.nominees.length">
            <div v-for="(nominee, index) in selectedUser.nominees" :key="index" class="nominee-card">
              <div class="detail-grid">
                <div class="detail-item"><strong>Name:</strong> {{ nominee.name }}</div>
                <div class="detail-item"><strong>Relationship:</strong> {{ nominee.relationship }}</div>
                <div class="detail-item"><strong>DOB:</strong> {{ nominee.dob }}</div>
                <div class="detail-item"><strong>PAN:</strong> {{ nominee.pan || 'N/A' }}</div>
                <div class="detail-item"><strong>Percentage:</strong> {{ nominee.percentage }}%</div>
              </div>
            </div>
          </div>
          <div v-else class="text-muted">No nominees found.</div>
        </section>

        <!-- IPV Details -->
        <section class="detail-section mt-4">
          <h3><i class="fas fa-video"></i> IPV Details</h3>
          <div class="ipv-container">
            <div class="ipv-photo">
              <strong>Photo:</strong>
              <div v-if="selectedUser.ipv_photo" class="photo-wrapper mt-2">
                <img :src="selectedUser.ipv_photo" alt="IPV Photo" />
              </div>
              <div v-else class="text-muted mt-2">No photo available.</div>
            </div>
            <div class="ipv-location">
              <strong>Location:</strong>
              <div class="mt-2">
                <div>Lat: {{ selectedUser.ipv_latitude || 'N/A' }}</div>
                <div>Lng: {{ selectedUser.ipv_longitude || 'N/A' }}</div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import '../../assets/css/user-management.css'
import { ref, onMounted, computed } from 'vue'
import { useToast } from 'vue-toastification'
import { adminApi } from '../../api/admin.api'

const toast = useToast()
const users = ref([])
const loading = ref(true)
const actionLoading = ref(null)
const searchQuery = ref('')
const showModal = ref(false)
const detailsLoading = ref(false)
const selectedUser = ref(null)

const fetchUsers = async () => {
  try {
    const res = await adminApi.getUsers()
    if (res.data.success) {
      users.value = res.data.data || []
    }
  } catch (err) {
    toast.error('Failed to load users')
    console.error(err)
  } finally {
    loading.value = false
  }
}

const filteredUsers = computed(() => {
  const closureUsers = users.value.filter(u => u.status === 'closure_requested')
  if (!searchQuery.value) return closureUsers
  const query = searchQuery.value.toLowerCase()
  return closureUsers.filter(u => 
    (u.name && u.name.toLowerCase().includes(query)) || 
    (u.email && u.email.toLowerCase().includes(query)) ||
    (u.pan && u.pan.toLowerCase().includes(query))
  )
})

const formatStatus = (status) => {
  if (!status) return 'Unknown'
  return status.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')
}

const getStatusClass = (status) => {
  switch(status) {
    case 'active': return 'status-active'
    case 'pending':
    case 'pending_approval': return 'status-pending'
    case 'blocked':
    case 'closure_requested':
    case 'rejected': return 'status-rejected'
    default: return ''
  }
}

const updateStatus = async (userId, action) => {
  actionLoading.value = userId
  try {
    const res = await adminApi.updateUserAction(userId, action)
    if (res.data.success) {
      toast.success(`User successfully ${action}d`)
      // Optimistically update the list
      const u = users.value.find(x => x.id === userId)
      if (u) {
        if (action === 'approve') u.status = 'active'
        else if (action === 'reject') u.status = 'rejected'
        else if (action === 'block') u.status = 'blocked'
        else if (action === 'close_account') u.status = 'closed'
        else if (action === 'reject_closure') u.status = 'active'
      }
    }
  } catch (err) {
    toast.error(err.response?.data?.message || `Failed to ${action} user`)
  } finally {
    actionLoading.value = null
  }
}

const openDetailsModal = async (userId) => {
  showModal.value = true
  detailsLoading.value = true
  selectedUser.value = null
  
  try {
    const res = await adminApi.getUserDetails(userId)
    if (res.data.success) {
      selectedUser.value = res.data.data
    }
  } catch (err) {
    toast.error('Failed to load user details')
    showModal.value = false
  } finally {
    detailsLoading.value = false
  }
}

const closeDetailsModal = () => {
  showModal.value = false
  selectedUser.value = null
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.btn-view {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
  border-color: rgba(59, 130, 246, 0.2);
  margin-right: 0.5rem;
}
.btn-view:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.2);
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.user-avatar.has-photo {
  background: transparent;
  overflow: hidden;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.modal-content {
  background: var(--bg-card, #1a1a24);
  border: 1px solid var(--border-color, rgba(255,255,255,0.1));
  border-radius: 12px;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 10px 25px rgba(0,0,0,0.5);
  animation: slideUp 0.3s ease;
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color, rgba(255,255,255,0.1));
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.25rem;
  color: var(--text-primary, #fff);
}

.btn-close {
  background: none;
  border: none;
  color: var(--text-secondary, #a0aec0);
  font-size: 1.25rem;
  cursor: pointer;
  padding: 0.5rem;
  transition: color 0.2s;
}

.btn-close:hover {
  color: var(--text-primary, #fff);
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
}

.loading-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  color: var(--text-secondary, #a0aec0);
}

.loading-body i {
  margin-bottom: 1rem;
  color: var(--primary-color, #3b82f6);
}

.detail-section {
  background: rgba(0,0,0,0.2);
  border-radius: 8px;
  padding: 1.25rem;
}

.detail-section h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  font-size: 1.1rem;
  color: var(--text-primary, #fff);
  border-bottom: 1px solid rgba(255,255,255,0.05);
  padding-bottom: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.detail-section h3 i {
  color: var(--primary-color, #3b82f6);
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
}

.detail-item {
  color: var(--text-primary, #fff);
  font-size: 0.9rem;
}

.detail-item strong {
  color: var(--text-secondary, #a0aec0);
  display: block;
  margin-bottom: 0.25rem;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-item.full-width {
  grid-column: 1 / -1;
}

.mt-2 { margin-top: 0.5rem; }
.mt-4 { margin-top: 1.5rem; }
.text-muted { color: var(--text-secondary, #a0aec0); }

.bank-card, .nominee-card {
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.05);
  border-radius: 6px;
  padding: 1rem;
  margin-bottom: 1rem;
}
.bank-card:last-child, .nominee-card:last-child {
  margin-bottom: 0;
}

.ipv-container {
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
}

.photo-wrapper {
  width: 200px;
  height: 200px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid rgba(255,255,255,0.1);
  background: #000;
}

.photo-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
