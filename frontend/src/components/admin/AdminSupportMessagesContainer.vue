<template>
  <div class="admin-support-page">
    <div class="header-section">
      <h1 class="page-title">Support Messages</h1>
      <button class="btn btn-outline" @click="fetchMessages" :disabled="loading">
        <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i> Refresh
      </button>
    </div>

    <div class="table-container">
      <div v-if="loading && messages.length === 0" class="loading-state">
        <i class="fas fa-circle-notch fa-spin"></i>
        <p>Loading messages...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <i class="fas fa-exclamation-triangle"></i>
        <p>{{ error }}</p>
      </div>

      <div v-else class="table-responsive">
        <table class="admin-table">
          <thead>
            <tr>
              <th>Sender</th>
              <th>Message</th>
              <th>Status</th>
              <th>Date</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="messages.length === 0">
              <td colspan="4" class="empty-state">No support messages found.</td>
            </tr>
            <tr v-for="msg in messages" :key="msg.id">
              <td>
                <div class="user-cell">
                  <div class="user-avatar">{{ msg.name ? msg.name.charAt(0).toUpperCase() : 'U' }}</div>
                  <div class="user-details">
                    <span class="user-name">{{ msg.name }}</span>
                    <span class="user-email">{{ msg.email }}</span>
                  </div>
                </div>
              </td>
              <td class="message-cell">
                <div class="message-content">{{ msg.message }}</div>
              </td>
              <td>
                <span :class="['status-badge', getStatusClass(msg.status)]">{{ msg.status || 'Open' }}</span>
              </td>
              <td>{{ formatDate(msg.createdAt) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminApi } from '../../api/admin.api'
import { useToast } from 'vue-toastification'

const messages = ref([])
const loading = ref(false)
const error = ref(null)
const toast = useToast()

const fetchMessages = async () => {
  loading.value = true
  error.value = null
  try {
    const response = await adminApi.getSupportMessages()
    messages.value = response.data.messages || []
  } catch (err) {
    console.error("Failed to fetch messages", err)
    error.value = err.response?.data?.error || "Failed to load messages."
    toast.error(error.value)
  } finally {
    loading.value = false
  }
}

const getStatusClass = (status) => {
  if (!status) return 'status-warning';
  switch (status.toLowerCase()) {
    case 'open': return 'status-warning'
    case 'resolved': return 'status-success'
    case 'closed': return 'status-danger'
    default: return 'status-default'
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'Unknown'
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit'
  }).format(date)
}

onMounted(() => {
  fetchMessages()
})
</script>

<style scoped>
.admin-support-page {
  animation: fadeIn 0.3s ease;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-title {
  margin: 0;
  font-size: 1.5rem;
  color: var(--text-primary, #fff);
}

.table-container {
  background: var(--bg-card, #1a1a24);
  border: 1px solid var(--border-color, rgba(255,255,255,0.05));
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.loading-state, .error-state, .empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: var(--text-secondary, #a0aec0);
}

.loading-state i, .error-state i {
  font-size: 2rem;
  margin-bottom: 1rem;
}

.error-state {
  color: var(--text-danger, #ef4444);
}

.admin-table {
  width: 100%;
  border-collapse: collapse;
}

.admin-table th, .admin-table td {
  padding: 1rem 1.5rem;
  text-align: left;
  border-bottom: 1px solid var(--border-color, rgba(255,255,255,0.05));
}

.admin-table th {
  background: rgba(0, 0, 0, 0.2);
  color: var(--text-secondary, #a0aec0);
  font-weight: 500;
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--primary-color, #3b82f6);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  flex-shrink: 0;
}

.user-details {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 500;
  color: var(--text-primary, #fff);
}

.user-email {
  font-size: 0.85rem;
  color: var(--text-secondary, #a0aec0);
}

.message-cell {
  max-width: 400px;
}

.message-content {
  white-space: pre-wrap;
  color: var(--text-secondary, #a0aec0);
  font-size: 0.95rem;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-warning { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }
.status-success { background: rgba(16, 185, 129, 0.1); color: #10b981; }
.status-danger { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.status-default { background: rgba(156, 163, 175, 0.1); color: #9ca3af; }

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
