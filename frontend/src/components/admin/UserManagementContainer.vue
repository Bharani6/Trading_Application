<template>
  <div class="usermanagementpage">
    <div class="dashboard-header">
      <h1>User Approvals</h1>
      <p>Manage user accounts and KYC applications</p>
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
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody v-if="loading">
            <tr>
              <td colspan="5" class="empty-state">
                <i class="fas fa-spinner fa-spin"></i> Loading users...
              </td>
            </tr>
          </tbody>
          <tbody v-else-if="filteredUsers.length === 0">
            <tr>
              <td colspan="5" class="empty-state">
                No users found matching your criteria.
              </td>
            </tr>
          </tbody>
          <tbody v-else>
            <tr v-for="user in filteredUsers" :key="user.id">
              <td>
                <div class="user-cell">
                  <div class="user-avatar">
                    {{ user.name ? user.name.charAt(0).toUpperCase() : 'U' }}
                  </div>
                  <div class="user-details">
                    <span class="user-name">{{ user.name || 'Unknown User' }}</span>
                    <span class="user-email">{{ user.email }}</span>
                  </div>
                </div>
              </td>
              <td class="docs-cell">
                <div>PAN: {{ user.pan || 'N/A' }}</div>
                <div style="color: var(--text-muted);">AADHAAR: {{ user.aadhaar ? '****' + user.aadhaar.slice(-4) : 'N/A' }}</div>
              </td>
              <td>{{ user.created_at || 'Unknown' }}</td>
              <td>
                <span :class="['status-badge', getStatusClass(user.status)]">
                  {{ user.status }}
                </span>
              </td>
              <td class="actions-cell">
                <button 
                  v-if="user.status === 'pending_approval' || user.status === 'pending' || user.status === 'rejected' || user.status === 'blocked'"
                  class="btn-action btn-approve"
                  @click="updateStatus(user.id, 'approve')"
                  :disabled="actionLoading === user.id"
                >
                  <i class="fas fa-check"></i> Approve
                </button>
                <button 
                  v-if="user.status === 'pending_approval' || user.status === 'pending'"
                  class="btn-action btn-reject"
                  @click="updateStatus(user.id, 'reject')"
                  :disabled="actionLoading === user.id"
                >
                  <i class="fas fa-times"></i> Reject
                </button>
                <span v-if="user.status === 'active'" style="color: var(--text-muted); font-size: 12px;">
                  No actions
                </span>
              </td>
            </tr>
          </tbody>
        </table>
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
  if (!searchQuery.value) return users.value
  const query = searchQuery.value.toLowerCase()
  return users.value.filter(u => 
    (u.name && u.name.toLowerCase().includes(query)) || 
    (u.email && u.email.toLowerCase().includes(query)) ||
    (u.pan && u.pan.toLowerCase().includes(query))
  )
})

const getStatusClass = (status) => {
  switch(status) {
    case 'active': return 'status-active'
    case 'pending':
    case 'pending_approval': return 'status-pending'
    case 'blocked':
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
        u.status = action === 'approve' ? 'active' : 'rejected'
      }
    }
  } catch (err) {
    toast.error(err.response?.data?.message || `Failed to ${action} user`)
  } finally {
    actionLoading.value = null
  }
}

onMounted(() => {
  fetchUsers()
})
</script>
