<template>
  <div class="stock-upload-page">
    <div class="header-section">
      <h1 class="page-title">Stock Management</h1>
      <p class="page-subtitle">Bulk upload stocks to the market via CSV</p>
    </div>

    <div class="content-section">
      <div class="custom-card">
        <h3 class="card-title">Upload CSV File</h3>
        
        <div class="info-alert">
          <i class="fas fa-info-circle"></i>
          <span>CSV Header must be exactly: <code>Name, Series, Segment, Price, Changes</code></span>
        </div>

        <div class="form-group">
          <label class="form-label">Select CSV file</label>
          <div class="file-input-wrapper">
            <input 
              type="file" 
              accept=".csv" 
              @change="handleFileChange" 
              class="file-input"
              id="csv-file"
            />
            <label for="csv-file" class="file-label">
              <i class="fas fa-upload"></i>
              <span>{{ file ? file.name : 'Choose a file...' }}</span>
            </label>
          </div>
        </div>

        <div class="form-actions">
          <button 
            class="primary-btn" 
            :disabled="!file || uploading" 
            @click="uploadCsv"
          >
            <i class="fas fa-spinner fa-spin" v-if="uploading"></i>
            <span v-else>Upload & Process</span>
          </button>
        </div>
      </div>

      <div class="custom-card danger-card">
        <h3 class="card-title text-danger">Danger Zone</h3>
        <p class="card-text">This will permanently delete all stocks from the market. Use this to reset the system.</p>
        
        <button 
          class="danger-btn" 
          :disabled="deleting" 
          @click="deleteAllStocks"
        >
          <i class="fas fa-spinner fa-spin" v-if="deleting"></i>
          <i class="fas fa-trash" v-else></i>
          <span>Delete All Stocks</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useToast } from 'vue-toastification'
import { adminApi } from '../../api/admin.api'

const toast = useToast()
const file = ref(null)
const uploading = ref(false)
const deleting = ref(false)

const handleFileChange = (e) => {
  if (e.target.files.length > 0) {
    file.value = e.target.files[0]
  } else {
    file.value = null
  }
}

const deleteAllStocks = async () => {
  if (!confirm('Are you sure you want to permanently delete all stocks? This action cannot be undone.')) return
  
  deleting.value = true
  try {
    const res = await adminApi.deleteShares()
    if (res.data.success) {
      toast.success('All stocks have been deleted.')
    } else {
      toast.error(res.data.message || 'Failed to delete stocks')
    }
  } catch (err) {
    console.error(err)
    toast.error('An error occurred while deleting stocks')
  } finally {
    deleting.value = false
  }
}

const uploadCsv = async () => {
  if (!file.value) return

  uploading.value = true
  const formData = new FormData()
  formData.append('file', file.value)

  try {
    const res = await adminApi.uploadShares(formData)
    if (res.data.success) {
      toast.success('Stocks uploaded successfully!')
      file.value = null
    } else {
      toast.error(res.data.message || 'Upload failed')
    }
  } catch (err) {
    console.error(err)
    toast.error('Failed to upload CSV')
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
.stock-upload-page {
  padding: 24px;
  max-width: 800px;
}

.header-section {
  margin-bottom: 32px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #f1f5f9;
  margin-bottom: 8px;
}

.page-subtitle {
  color: #64748b;
  font-size: 16px;
}

.content-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.custom-card {
  background-color: var(--surface-color, #131929);
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: 16px;
  padding: 24px;
}

.card-title {
  font-size: 20px;
  font-weight: 600;
  color: #f1f5f9;
  margin-bottom: 24px;
}

.text-danger {
  color: #f43f5e;
}

.card-text {
  color: #94a3b8;
  margin-bottom: 24px;
}

.info-alert {
  background-color: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  color: #60a5fa;
  margin-bottom: 24px;
}

.info-alert code {
  background-color: rgba(0, 0, 0, 0.2);
  padding: 2px 6px;
  border-radius: 4px;
  color: #f1f5f9;
}

.form-group {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  color: #f1f5f9;
  margin-bottom: 8px;
  font-weight: 500;
}

.file-input {
  display: none;
}

.file-label {
  display: flex;
  align-items: center;
  gap: 12px;
  background-color: rgba(255, 255, 255, 0.03);
  border: 1px dashed rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 16px 24px;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.2s ease;
}

.file-label:hover {
  background-color: rgba(255, 255, 255, 0.05);
  border-color: var(--primary-color, #3b82f6);
  color: #f1f5f9;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.primary-btn {
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 12px 24px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.primary-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.primary-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.danger-btn {
  background-color: transparent;
  color: #f43f5e;
  border: 1px solid #f43f5e;
  border-radius: 8px;
  padding: 12px 24px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.danger-btn:hover:not(:disabled) {
  background-color: rgba(244, 63, 94, 0.1);
}

.danger-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
