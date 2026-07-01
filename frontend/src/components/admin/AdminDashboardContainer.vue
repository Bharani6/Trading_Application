<template>
  <div class="admin-dashboard-page">
    <div class="header-section">
      <h1 class="page-title">Admin Dashboard</h1>
    </div>

    <!-- Quick Stats Grid -->
    <div class="stats-grid">
      <!-- Total Users -->
      <div class="stat-card">
        <div class="stat-icon-wrapper bg-primary-light text-primary">
          <i class="fas fa-users stat-icon"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">Total Users</div>
          <div v-if="!loading" class="stat-value">{{ stats.totalUsers }}</div>
          <div v-else class="stat-value loading"><i class="fas fa-circle-notch fa-spin"></i></div>
        </div>
      </div>

      <!-- Pending Approvals -->
      <div class="stat-card">
        <div class="stat-icon-wrapper bg-warning-light text-warning">
          <i class="fas fa-clock stat-icon"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">Pending Approvals</div>
          <div v-if="!loading" class="stat-value">{{ stats.pendingUsers }}</div>
          <div v-else class="stat-value loading"><i class="fas fa-circle-notch fa-spin"></i></div>
        </div>
      </div>

      <!-- Active Users -->
      <div class="stat-card">
        <div class="stat-icon-wrapper bg-success-light text-success">
          <i class="fas fa-check-circle stat-icon"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">Active Users</div>
          <div v-if="!loading" class="stat-value">{{ stats.activeUsers }}</div>
          <div v-else class="stat-value loading"><i class="fas fa-circle-notch fa-spin"></i></div>
        </div>
      </div>

      <!-- Blocked Users -->
      <div class="stat-card">
        <div class="stat-icon-wrapper bg-danger-light text-danger">
          <i class="fas fa-ban stat-icon"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">Rejected / Blocked</div>
          <div v-if="!loading" class="stat-value">{{ stats.blockedUsers }}</div>
          <div v-else class="stat-value loading"><i class="fas fa-circle-notch fa-spin"></i></div>
        </div>
      </div>
    </div>

    <!-- Analytics & Visualizations Grid Removed per user request -->

    <!-- Tables Grid -->
    <div class="tables-grid">
      <!-- Pending Approvals Quick-View -->
      <div class="table-container">
        <div class="table-header-flex">
          <h3 class="section-title">Pending Approvals</h3>
          <button @click="$router.push('/admin/users')" class="view-all-btn">View All</button>
        </div>
        <div class="table-responsive">
          <table class="admin-table">
            <thead>
              <tr>
                <th>User</th>
                <th>Applied Date</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="pendingUsersList.length === 0">
                <td colspan="3" class="empty-state">No pending approvals at the moment.</td>
              </tr>
              <tr v-for="user in pendingUsersList" :key="user.id">
                <td>
                  <div class="user-cell">
                    <div class="user-avatar">{{ user.name ? user.name.charAt(0).toUpperCase() : 'U' }}</div>
                    <div class="user-details">
                      <span class="user-name">{{ user.name || 'Unknown' }}</span>
                      <span class="user-email">{{ user.email }}</span>
                    </div>
                  </div>
                </td>
                <td>{{ user.created_at || 'Unknown' }}</td>
                <td class="actions-cell">
                  <button class="btn-action btn-approve" @click="quickApprove(user.id)" :disabled="actionLoading === user.id">
                    <i class="fas fa-check"></i>
                  </button>
                  <button class="btn-action btn-reject" @click="quickReject(user.id)" :disabled="actionLoading === user.id">
                    <i class="fas fa-times"></i>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Recent Registrations -->
      <div class="table-container">
        <h3 class="section-title">Recent Registrations</h3>
        <div class="table-responsive">
          <table class="admin-table">
            <thead>
              <tr>
                <th>User</th>
                <th>Joined Date</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="recentUsers.length === 0">
                <td colspan="3" class="empty-state">No users registered yet.</td>
              </tr>
              <tr v-for="user in recentUsers" :key="user.id">
                <td>
                  <div class="user-cell">
                    <div class="user-avatar">{{ user.name ? user.name.charAt(0).toUpperCase() : 'U' }}</div>
                    <div class="user-details">
                      <span class="user-name">{{ user.name || 'Unknown' }}</span>
                      <span class="user-email">{{ user.email }}</span>
                    </div>
                  </div>
                </td>
                <td>{{ user.created_at || 'Unknown' }}</td>
                <td>
                  <span :class="['status-badge', getStatusClass(user.status)]">{{ user.status }}</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import { adminApi } from '../../api/admin.api'
import Chart from 'chart.js/auto'

const router = useRouter()
const toast = useToast()
const loading = ref(true)
const actionLoading = ref(null)

// Stats & Users Arrays
const stats = ref({
  totalUsers: 0,
  activeUsers: 0,
  pendingUsers: 0,
  blockedUsers: 0
})
const allUsers = ref([])
const recentUsers = ref([])
const pendingUsersList = ref([])

// Chart & UI State
const chartCanvas = ref(null)
const chartInstance = ref(null)

// Mocked Audit Logs (Backend endpoint pending)
const auditLogs = ref([
  { id: 1, action: 'User "testing" registered an account', time: '10 mins ago', type: 'info' },
  { id: 2, action: 'Admin approved KYC for "Bharanidharan"', time: '1 hour ago', type: 'success' },
  { id: 3, action: 'System executed end-of-day settlement', time: '5 hours ago', type: 'warning' },
  { id: 4, action: 'Admin uploaded new stock list', time: '1 day ago', type: 'info' }
])

const getAuditIcon = (type) => {
  switch(type) {
    case 'success': return 'fas fa-check-circle'
    case 'warning': return 'fas fa-exclamation-triangle'
    case 'danger': return 'fas fa-times-circle'
    default: return 'fas fa-info-circle'
  }
}

const getStatusClass = (status) => {
  if (status === 'active') return 'status-active'
  if (status === 'pending' || status === 'pending_approval') return 'status-pending'
  return 'status-rejected'
}

const fetchStats = async () => {
  try {
    const res = await adminApi.getUsers()
    if (res.data.success) {
      const users = res.data.data || []
      allUsers.value = users
      
      stats.value.totalUsers = users.length
      stats.value.activeUsers = users.filter(u => u.status === 'active').length
      stats.value.pendingUsers = users.filter(u => u.status === 'pending' || u.status === 'pending_approval').length
      stats.value.blockedUsers = users.filter(u => u.status === 'blocked' || u.status === 'rejected').length

      // Setup Recent Users (Sort by CreatedAt roughly, assume array order or basic sort)
      recentUsers.value = [...users].reverse().slice(0, 5)

      // Setup Pending Users
      pendingUsersList.value = users.filter(u => u.status === 'pending' || u.status === 'pending_approval').slice(0, 5)

    }
  } catch (err) {
    console.error(err)
    toast.error('Failed to load dashboard statistics')
  } finally {
    loading.value = false
  }
}

const quickApprove = async (userId) => {
  actionLoading.value = userId
  try {
    const res = await adminApi.approveUser(userId)
    if (res.data.success) {
      toast.success('User approved successfully')
      fetchStats() // Refresh data
    }
  } catch (err) {
    toast.error('Failed to approve user')
  } finally {
    actionLoading.value = null
  }
}

const quickReject = async (userId) => {
  actionLoading.value = userId
  try {
    const res = await adminApi.rejectUser(userId)
    if (res.data.success) {
      toast.success('User rejected successfully')
      fetchStats() // Refresh data
    }
  } catch (err) {
    toast.error('Failed to reject user')
  } finally {
    actionLoading.value = null
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.admin-dashboard-page {
  padding: 24px;
  max-width: 1400px;
  position: relative;
}

.header-section { margin-bottom: 32px; }
.page-title { font-size: 28px; font-weight: 700; color: #f1f5f9; margin-bottom: 8px; }
.page-subtitle { color: #64748b; font-size: 16px; }

/* Stats Grid */
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 24px; margin-bottom: 32px; }
.stat-card { background-color: var(--surface-card, #131929); border: 1px solid rgba(255,255,255,0.07); border-radius: 16px; padding: 24px; display: flex; align-items: center; transition: all 0.3s; }
.stat-card:hover { transform: translateY(-4px); box-shadow: 0 12px 24px rgba(0,0,0,0.2); }
.stat-icon-wrapper { width: 56px; height: 56px; border-radius: 12px; display: flex; align-items: center; justify-content: center; margin-right: 16px; flex-shrink: 0; font-size: 24px; }
.stat-content { display: flex; flex-direction: column; }
.stat-label { font-size: 12px; font-weight: 600; text-transform: uppercase; color: #94a3b8; margin-bottom: 4px; }
.stat-value { font-size: 28px; font-weight: 700; color: #f1f5f9; }

/* Analytics Grid */
.analytics-grid { display: grid; grid-template-columns: 2fr 1fr; gap: 24px; margin-bottom: 32px; }
.chart-container, .audit-log-container { background: var(--surface-card, #131929); border: 1px solid rgba(255,255,255,0.07); border-radius: 16px; padding: 24px; }
.section-title { font-size: 18px; font-weight: 600; color: #f1f5f9; margin-bottom: 20px; }
.chart-container canvas { width: 100% !important; height: 300px !important; }

/* Audit Log */
.audit-feed { display: flex; flex-direction: column; gap: 16px; max-height: 300px; overflow-y: auto; padding-right: 8px; }
.audit-item { display: flex; gap: 12px; padding-bottom: 16px; border-bottom: 1px solid rgba(255,255,255,0.05); }
.audit-item:last-child { border-bottom: none; }
.audit-icon { width: 32px; height: 32px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 14px; flex-shrink: 0; }
.audit-info { background: rgba(59,130,246,0.1); color: #3b82f6; }
.audit-success { background: rgba(16,185,129,0.1); color: #10b981; }
.audit-warning { background: rgba(245,158,11,0.1); color: #f59e0b; }
.audit-content p { margin: 0; font-size: 14px; color: #e2e8f0; line-height: 1.4; }
.audit-time { font-size: 12px; color: #64748b; }

/* Tables Grid */
.tables-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 24px; }
.table-container { background: var(--surface-card, #131929); border: 1px solid rgba(255,255,255,0.07); border-radius: 16px; padding: 24px; }
.table-header-flex { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.table-header-flex .section-title { margin-bottom: 0; }
.view-all-btn { background: transparent; border: none; color: #3b82f6; cursor: pointer; font-weight: 600; font-size: 14px; }
.view-all-btn:hover { text-decoration: underline; }

.table-responsive { overflow-x: auto; }
.admin-table { width: 100%; border-collapse: collapse; text-align: left; }
.admin-table th { padding: 12px 16px; font-size: 12px; text-transform: uppercase; color: #64748b; font-weight: 700; border-bottom: 1px solid rgba(255,255,255,0.07); }
.admin-table td { padding: 16px; border-bottom: 1px solid rgba(255,255,255,0.03); color: #f1f5f9; font-size: 14px; vertical-align: middle; }
.admin-table tr:hover td { background-color: rgba(255,255,255,0.02); }

.user-cell { display: flex; align-items: center; gap: 12px; }
.user-avatar { width: 36px; height: 36px; border-radius: 50%; background: #1a2236; display: flex; align-items: center; justify-content: center; font-weight: 700; color: #f1f5f9; border: 1px solid rgba(255,255,255,0.07); }
.user-details { display: flex; flex-direction: column; }
.user-name { font-weight: 600; }
.user-email { font-size: 12px; color: #64748b; }

.status-badge { padding: 4px 10px; border-radius: 20px; font-size: 12px; font-weight: 700; text-transform: uppercase; }
.status-active { background: rgba(16,185,129,0.15); color: #10b981; border: 1px solid rgba(16,185,129,0.3); }
.status-pending { background: rgba(234,179,8,0.15); color: #eab308; border: 1px solid rgba(234,179,8,0.3); }
.status-rejected { background: rgba(244,63,94,0.15); color: #f43f5e; border: 1px solid rgba(244,63,94,0.3); }

.actions-cell { display: flex; gap: 8px; }
.btn-action { padding: 6px 12px; border: none; border-radius: 6px; cursor: pointer; transition: 0.2s; }
.btn-approve { background: rgba(16,185,129,0.1); color: #10b981; border: 1px solid rgba(16,185,129,0.2); }
.btn-approve:hover { background: #10b981; color: #000; }
.btn-reject { background: rgba(244,63,94,0.1); color: #f43f5e; border: 1px solid rgba(244,63,94,0.2); }
.btn-reject:hover { background: #f43f5e; color: #fff; }

.empty-state { text-align: center; color: #64748b; padding: 32px !important; }

@media (max-width: 1024px) {
  .analytics-grid, .tables-grid { grid-template-columns: 1fr; }
}
</style>
