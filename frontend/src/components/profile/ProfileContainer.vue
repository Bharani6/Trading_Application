<template>
  <div class="profilepage">
    <div class="header-section">
      <div>
        <h1 class="page-title">Account Settings</h1>
        <p class="page-subtitle">Manage your account details and personal information</p>
      </div>
      <div :class="['kyc-status', `status-${authStore.state.user?.role === 'admin' ? 'active' : (authStore.state.user?.status || 'pending')}`]">
        <i :class="(authStore.state.user?.role === 'admin' || authStore.state.user?.status === 'active') ? 'fas fa-check-circle' : 'fas fa-clock'"></i>
        {{ (authStore.state.user?.role === 'admin' || authStore.state.user?.status === 'active') ? 'Verified' : authStore.state.user?.status === 'rejected' ? 'Rejected' : 'Verification Pending' }}
      </div>
    </div>

    <!-- User Info Card -->
    <div class="user-info-card">
      <div class="avatar-wrap">
        <div class="avatar-circle">
          <span>{{ initials }}</span>
        </div>
        <div class="avatar-details">
          <h2 class="user-full-name">{{ authStore.state.user?.name || 'Not Provided' }}</h2>
          <p class="user-email">{{ authStore.state.user?.email || 'Not Provided' }}</p>
          <span :class="['role-badge', `role-${authStore.state.user?.role || 'user'}`]">
            {{ authStore.state.user?.role || 'user' }}
          </span>
        </div>
      </div>
      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">Account ID</span>
          <span class="info-value id-value">{{ authStore.state.user?.id?.slice(0, 8) || 'Not Provided' }}...</span>
        </div>
        <div class="info-item">
          <span class="info-label">Account Status</span>
          <span :class="['status-chip', `chip-${authStore.state.user?.role === 'admin' ? 'active' : (authStore.state.user?.status || 'pending')}`]">
            {{ authStore.state.user?.role === 'admin' ? 'active' : (authStore.state.user?.status || 'pending') }}
          </span>
        </div>
      </div>
    </div>

    <!-- Settings Section -->
    <div class="settings-section mt-4">
      <div class="settings-content">
        <form @submit.prevent="saveProfile" class="grid-form">
          <!-- PERSONAL DETAILS TAB -->
          <template v-if="activeTab === 'personal-details'">
            <div class="full-width section-title">
              <h3><i class="fas fa-user"></i> Personal Details</h3>
              <hr class="section-divider">
            </div>
            
            <div class="form-group">
              <label class="form-label">Full Name</label>
              <div class="static-value">{{ authStore.state.user?.name || 'Not Provided' }}</div>
            </div>
            <div class="form-group">
              <label class="form-label">Email Address</label>
              <div class="static-value">{{ authStore.state.user?.email || 'Not Provided' }}</div>
            </div>
            <div class="form-group">
              <div class="flex-row-between">
                <label class="form-label">Mobile Number</label>
                <i v-if="editingField !== 'mobile'" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = 'mobile'"></i>
                <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
              </div>
              <div v-if="editingField !== 'mobile'" class="static-value">
                {{ formatMobile(kycForm.mobile) }}
              </div>
              <div v-else class="input-with-prefix">
                <span class="prefix">+91</span>
                <input type="text" class="form-input prefix-input" v-model="kycForm.mobile" placeholder="10-digit number" />
              </div>
            </div>
            <div class="form-group">
              <div class="flex-row-between">
                <label class="form-label">Father's Name</label>
                <i v-if="editingField !== 'fatherName'" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = 'fatherName'"></i>
                <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
              </div>
              <div v-if="editingField !== 'fatherName'" class="static-value">{{ kycForm.fatherName || 'Not Provided' }}</div>
              <input v-else type="text" class="form-input" v-model="kycForm.fatherName" placeholder="Enter father's name" />
            </div>
            <div class="form-group">
              <div class="flex-row-between">
                <label class="form-label">Mother's Name</label>
                <i v-if="editingField !== 'motherName'" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = 'motherName'"></i>
                <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
              </div>
              <div v-if="editingField !== 'motherName'" class="static-value">{{ kycForm.motherName || 'Not Provided' }}</div>
              <input v-else type="text" class="form-input" v-model="kycForm.motherName" placeholder="Enter mother's name" />
            </div>
            <div class="form-group">
              <div class="flex-row-between">
                <label class="form-label">Date of Birth</label>
                <i v-if="editingField !== 'dob'" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = 'dob'"></i>
                <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
              </div>
              <div v-if="editingField !== 'dob'" class="static-value">{{ kycForm.dob || 'Not Provided' }}</div>
              <input v-else type="date" class="form-input" v-model="kycForm.dob" />
            </div>
            <div class="form-group">
              <div class="flex-row-between">
                <label class="form-label">Income Range</label>
                <i v-if="editingField !== 'incomeRange'" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = 'incomeRange'"></i>
                <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
              </div>
              <div v-if="editingField !== 'incomeRange'" class="static-value">{{ kycForm.incomeRange || 'Not Provided' }}</div>
              <select v-else class="form-input" v-model="kycForm.incomeRange">
                <option value="" disabled>Select Income Range</option>
                <option value="Below 1 Lakh">Below 1 Lakh</option>
                <option value="1-5 Lakhs">1-5 Lakhs</option>
                <option value="5-10 Lakhs">5-10 Lakhs</option>
                <option value="Above 10 Lakhs">Above 10 Lakhs</option>
              </select>
            </div>
            
            <div class="full-width mt-4"><h4 style="color: var(--text-muted); font-size: 14px; border-bottom: 1px dashed var(--border-color); padding-bottom: 4px;">Address Details</h4></div>
            
            <div class="form-group">
              <div class="flex-row-between">
                <label class="form-label">Pincode</label>
                <i v-if="editingField !== 'pincode'" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = 'pincode'"></i>
                <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
              </div>
              <div v-if="editingField !== 'pincode'" class="static-value">{{ kycForm.pincode || 'Not Provided' }}</div>
              <input v-else type="text" class="form-input" v-model="kycForm.pincode" placeholder="6-digit Pincode" maxlength="6" />
            </div>
            <div class="form-group">
              <label class="form-label">Country</label>
              <div class="static-value">{{ kycForm.country || 'Not Provided' }}</div>
            </div>
            <div class="form-group">
              <label class="form-label">State</label>
              <div class="static-value">{{ kycForm.state || 'Not Provided' }}</div>
            </div>
            <div class="form-group">
              <label class="form-label">City</label>
              <div class="static-value">{{ kycForm.city || 'Not Provided' }}</div>
            </div>
            <div class="form-group full-width">
              <div class="flex-row-between">
                <label class="form-label">Full Address</label>
                <i v-if="editingField !== 'address'" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = 'address'"></i>
                <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
              </div>
              <div v-if="editingField !== 'address'" class="static-value" style="min-height: 60px;">{{ kycForm.address || 'Not Provided' }}</div>
              <textarea v-else class="form-textarea" v-model="kycForm.address" placeholder="Enter your full residential address" rows="3"></textarea>
            </div>
          </template>

          <!-- CHANGE PASSWORD TAB -->
          <template v-if="activeTab === 'change-password'">
            <div class="full-width section-title">
              <h3><i class="fas fa-lock"></i> Change Password</h3>
              <hr class="section-divider">
            </div>
            <div class="form-group full-width" style="max-width: 400px; margin-bottom: 20px;">
              <label class="form-label">Current Password</label>
              <input type="password" class="form-input" placeholder="Enter current password" />
            </div>
            <div class="form-group full-width" style="max-width: 400px; margin-bottom: 20px;">
              <label class="form-label">New Password</label>
              <input type="password" class="form-input" placeholder="Enter new password" />
            </div>
            <div class="form-group full-width" style="max-width: 400px; margin-bottom: 20px;">
              <label class="form-label">Confirm New Password</label>
              <input type="password" class="form-input" placeholder="Confirm new password" />
            </div>
            <div class="full-width">
              <button type="button" class="btn-outline-small" style="padding: 10px 20px;" @click="toast.info('Password change functionality coming soon')">
                <i class="fas fa-key"></i> Update Password
              </button>
            </div>
          </template>

          <!-- BANK ACCOUNTS TAB -->
          <template v-if="activeTab === 'bank-accounts'">
            <div class="full-width section-title">
              <div class="title-with-action">
                <h3><i class="fas fa-university"></i> Bank Accounts</h3>
                <button v-if="kycForm.bankAccounts.length === 1" type="button" class="btn-outline-small" @click="addSecondaryBank">
                  <i class="fas fa-plus"></i> Add Secondary Bank
                </button>
              </div>
              <hr class="section-divider">
            </div>

            <template v-for="(bank, index) in kycForm.bankAccounts" :key="index">
              <div class="full-width bank-header">
                <h4>{{ bank.accountType === 'primary' ? 'Primary' : 'Secondary' }} Bank Account</h4>
                <button v-if="bank.accountType === 'secondary'" type="button" class="btn-text-red" @click="removeSecondaryBank(index)">
                  <i class="fas fa-trash"></i> Remove
                </button>
              </div>
              
              <div class="form-group">
                <label class="form-label">Account Holder Name</label>
                <div class="static-value">{{ authStore.state.user?.name || 'Not Provided' }}</div>
              </div>
              <div class="form-group">
                <div class="flex-row-between">
                  <label class="form-label">Account Number</label>
                  <i v-if="editingField !== `bank_${index}_acc`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `bank_${index}_acc`"></i>
                  <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                </div>
                <div v-if="editingField !== `bank_${index}_acc`" class="static-value">{{ bank.accountNumber || 'Not Provided' }}</div>
                <input v-else type="text" class="form-input" v-model="bank.accountNumber" @input="bank.accountNumber = $event.target.value.replace(/[^0-9]/g, '')" placeholder="9-18 digit account number" maxlength="18" />
              </div>
              <div class="form-group">
                <div class="flex-row-between">
                  <label class="form-label">IFSC Code</label>
                  <i v-if="editingField !== `bank_${index}_ifsc`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `bank_${index}_ifsc`"></i>
                  <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                </div>
                <div v-if="editingField !== `bank_${index}_ifsc`" class="static-value">{{ bank.ifsc || 'Not Provided' }}</div>
                <input v-else type="text" class="form-input" v-model="bank.ifsc" @input="bank.ifsc = $event.target.value.toUpperCase().replace(/[^A-Z0-9]/g, '')" placeholder="11-character IFSC code" maxlength="11" />
              </div>
              <div class="form-group">
                <label class="form-label">Bank Name</label>
                <div class="static-value">{{ bank.bankName || 'Not Provided' }}</div>
              </div>
              <div class="form-group">
                <label class="form-label">Branch Name</label>
                <div class="static-value">{{ bank.branch || 'Not Provided' }}</div>
              </div>
              <div class="full-width"><br/></div>
            </template>
          </template>

          <!-- NOMINEE DETAILS TAB -->
          <template v-if="activeTab === 'nominee-details'">
            <div class="full-width section-title">
              <div class="title-with-action">
                <h3><i class="fas fa-users"></i> Nominee Details</h3>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="kycForm.nomineeEnabled">
                  <span class="slider round"></span>
                </label>
              </div>
              <hr class="section-divider">
              <p v-if="!kycForm.nomineeEnabled" class="muted-text text-sm">I do not wish to nominate anyone at this time.</p>
              <div v-else class="flex-row-between mb-2">
                <p class="muted-text text-sm">Total Allocation: <strong :class="totalNomineePercentage === 100 ? 'text-green' : 'text-red'">{{ totalNomineePercentage }}%</strong> / 100%</p>
                <button type="button" class="btn-outline-small" v-if="kycForm.nominees.length < 3" @click="addNominee">
                  <i class="fas fa-plus"></i> Add Nominee
                </button>
              </div>
            </div>

            <template v-if="kycForm.nomineeEnabled">
              <template v-for="(nom, index) in kycForm.nominees" :key="'nom_'+index">
                <div class="full-width bank-header mt-2">
                  <h4>Nominee {{ index + 1 }}</h4>
                  <button v-if="kycForm.nominees.length > 1" type="button" class="btn-text-red" @click="removeNominee(index)">
                    <i class="fas fa-trash"></i> Remove
                  </button>
                </div>

                <div class="form-group">
                  <div class="flex-row-between">
                    <label class="form-label">Nominee Full Name</label>
                    <i v-if="editingField !== `nom_name_${index}`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `nom_name_${index}`"></i>
                    <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                  </div>
                  <div v-if="editingField !== `nom_name_${index}`" class="static-value">{{ nom.name || 'Not Provided' }}</div>
                  <input v-else type="text" class="form-input" v-model="nom.name" placeholder="Enter nominee name" />
                </div>
                <div class="form-group">
                  <div class="flex-row-between">
                    <label class="form-label">Relationship with Applicant</label>
                    <i v-if="editingField !== `nom_rel_${index}`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `nom_rel_${index}`"></i>
                    <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                  </div>
                  <div v-if="editingField !== `nom_rel_${index}`" class="static-value">{{ nom.relationship || 'Not Provided' }}</div>
                  <select v-else class="form-input" v-model="nom.relationship">
                    <option value="" disabled>Select Relationship</option>
                    <option value="Spouse">Spouse</option>
                    <option value="Son">Son</option>
                    <option value="Daughter">Daughter</option>
                    <option value="Father">Father</option>
                    <option value="Mother">Mother</option>
                    <option value="Brother">Brother</option>
                    <option value="Sister">Sister</option>
                    <option value="Other">Other</option>
                  </select>
                </div>
                <div class="form-group">
                  <div class="flex-row-between">
                    <label class="form-label">Nominee Date of Birth</label>
                    <i v-if="editingField !== `nom_dob_${index}`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `nom_dob_${index}`"></i>
                    <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                  </div>
                  <div v-if="editingField !== `nom_dob_${index}`" class="static-value">{{ nom.dob || 'Not Provided' }}</div>
                  <input v-else type="date" class="form-input" v-model="nom.dob" />
                </div>
                <div class="form-group">
                  <div class="flex-row-between">
                    <label class="form-label">Nominee PAN</label>
                    <i v-if="editingField !== `nom_pan_${index}`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `nom_pan_${index}`"></i>
                    <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                  </div>
                  <div v-if="editingField !== `nom_pan_${index}`" class="static-value">{{ nom.pan || 'Not Provided' }}</div>
                  <input v-else type="text" class="form-input" v-model="nom.pan" placeholder="Enter nominee PAN" />
                </div>
                <div class="form-group">
                  <div class="flex-row-between">
                    <label class="form-label">Allocation Percentage (%)</label>
                    <i v-if="editingField !== `nom_perc_${index}`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `nom_perc_${index}`"></i>
                    <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                  </div>
                  <div v-if="editingField !== `nom_perc_${index}`" class="static-value">{{ nom.percentage }}%</div>
                  <input v-else type="number" class="form-input" v-model="nom.percentage" min="1" max="100" placeholder="e.g. 100" />
                </div>

                <!-- GUARDIAN DETAILS (Only if Nominee is Minor) -->
                <template v-if="isNomineeMinor(index)">
                  <div class="full-width guardian-banner mt-2">
                    <i class="fas fa-info-circle"></i> Nominee is below 18 years of age. Guardian details are mandatory as per SEBI guidelines.
                  </div>
                  <div class="form-group">
                    <div class="flex-row-between">
                      <label class="form-label">Guardian Name</label>
                      <i v-if="editingField !== `g_name_${index}`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `g_name_${index}`"></i>
                      <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                    </div>
                    <div v-if="editingField !== `g_name_${index}`" class="static-value">{{ nom.guardianName || 'Not Provided' }}</div>
                    <input v-else type="text" class="form-input" v-model="nom.guardianName" placeholder="Enter guardian name" />
                  </div>
                  <div class="form-group">
                    <div class="flex-row-between">
                      <label class="form-label">Guardian Date of Birth</label>
                      <i v-if="editingField !== `g_dob_${index}`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `g_dob_${index}`"></i>
                      <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                    </div>
                    <div v-if="editingField !== `g_dob_${index}`" class="static-value">{{ nom.guardianDob || 'Not Provided' }}</div>
                    <input v-else type="date" class="form-input" v-model="nom.guardianDob" :max="eighteenYearsAgo" title="Guardian must be 18 years or older" />
                  </div>
                  <div class="form-group">
                    <div class="flex-row-between">
                      <label class="form-label">Guardian Relationship</label>
                      <i v-if="editingField !== `g_rel_${index}`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `g_rel_${index}`"></i>
                      <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                    </div>
                    <div v-if="editingField !== `g_rel_${index}`" class="static-value">{{ nom.guardianRelationship || 'Not Provided' }}</div>
                    <select v-else class="form-input" v-model="nom.guardianRelationship">
                      <option value="" disabled>Select Guardian Relationship</option>
                      <option value="Father">Father</option>
                      <option value="Mother">Mother</option>
                      <option value="Legal Guardian">Legal Guardian</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <div class="flex-row-between">
                      <label class="form-label">Guardian PAN</label>
                      <i v-if="editingField !== `g_pan_${index}`" class="fas fa-pencil-alt text-green cursor-pointer" @click="editingField = `g_pan_${index}`"></i>
                      <i v-else class="fas fa-check text-green cursor-pointer" @click="editingField = null"></i>
                    </div>
                    <div v-if="editingField !== `g_pan_${index}`" class="static-value">{{ nom.guardianPan || 'Not Provided' }}</div>
                    <input v-else type="text" class="form-input" v-model="nom.guardianPan" placeholder="Enter guardian PAN" />
                  </div>
                </template>
                <div class="full-width"><br/></div>
              </template>
            </template>
          </template>

          <div class="form-actions full-width mt-4 text-right" v-if="activeTab !== 'change-password'">
            <button type="submit" class="submit-btn" :disabled="loading || (kycForm.nomineeEnabled && activeTab === 'nominee-details' && totalNomineePercentage !== 100)">
              <span v-if="!loading"><i class="fas fa-save"></i> Save Changes</span>
              <i v-else class="fas fa-spinner fa-spin"></i>
            </button>
            <p v-if="kycForm.nomineeEnabled && activeTab === 'nominee-details' && totalNomineePercentage !== 100" class="text-red text-sm mt-2">
              <i class="fas fa-exclamation-circle"></i> Allocation must equal exactly 100% to save.
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import '../../assets/css/profile.css'
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useAuthStore } from '../../store'
import { useToast } from 'vue-toastification'
import { useRoute } from 'vue-router'
import { userApi } from '../../api/user.api'
import { utilsApi } from '../../api/utils.api'

const route = useRoute()
const authStore = useAuthStore()
const toast = useToast()
const loading = ref(false)
const editingField = ref(null)

const activeTab = computed(() => route.query.tab || 'personal-details')

const newNomineeObj = () => ({
  name: '',
  dob: '',
  pan: '',
  relationship: '',
  percentage: 100,
  guardianName: '',
  guardianDob: '',
  guardianRelationship: '',
  guardianPan: ''
})

const kycForm = reactive({
  mobile: '',
  fatherName: '',
  motherName: '',
  dob: '',
  incomeRange: '',
  country: '',
  state: '',
  city: '',
  pincode: '',
  address: '',
  bankAccounts: [
    { accountType: 'primary', bankName: '', accountNumber: '', ifsc: '', branch: '', _lastIfsc: '' }
  ],
  nomineeEnabled: false,
  nominees: [newNomineeObj()]
})

const totalNomineePercentage = computed(() => {
  return kycForm.nominees.reduce((sum, nom) => sum + (Number(nom.percentage) || 0), 0)
})

const addNominee = () => {
  if (kycForm.nominees.length < 3) {
    kycForm.nominees.push(newNomineeObj())
    // Auto-adjust percentages roughly
    const val = Math.floor(100 / kycForm.nominees.length)
    kycForm.nominees.forEach((n, idx) => {
      n.percentage = (idx === kycForm.nominees.length - 1) ? 100 - (val * (kycForm.nominees.length - 1)) : val
    })
  }
}

const removeNominee = (index) => {
  kycForm.nominees.splice(index, 1)
  // Give remaining percentage to the first one
  if (kycForm.nominees.length > 0) {
    const sum = kycForm.nominees.reduce((s, n) => s + (Number(n.percentage) || 0), 0)
    if (sum !== 100) {
       kycForm.nominees[0].percentage += (100 - sum)
    }
  }
}

const formatMobile = (mobile) => {
  if (!mobile) return 'Not Provided'
  // Remove existing +91 if present to avoid duplication
  const cleanMobile = mobile.replace(/^\+91\s*/, '')
  return `+91 ${cleanMobile}`
}

// Auto-fetch Pincode using Watch
watch(() => kycForm.pincode, async (newVal) => {
  if (newVal && newVal.length === 6) {
    try {
      // Use centralized api
      const json = await utilsApi.getPincodeDetails(newVal)
      if (json && json.length > 0 && json[0].Status === "Success") {
        const postOffices = json[0].PostOffice
        if (postOffices && postOffices.length > 0) {
          kycForm.state = postOffices[0].State || ''
          kycForm.city = postOffices[0].District || ''
          kycForm.country = postOffices[0].Country || 'India'
          toast.success('Pincode details fetched automatically')
        }
      }
    } catch (err) {
      // Silently ignore so we don't bother user if API fails
    }
  }
})

// Auto-fetch IFSC using Watch
watch(() => kycForm.bankAccounts, (newVal) => {
  newVal.forEach(async (bank) => {
    if (bank.ifsc && bank.ifsc.length === 11 && bank.ifsc !== bank._lastIfsc) {
      bank._lastIfsc = bank.ifsc // Prevent loop/duplicate calls
      try {
        const response = await utilsApi.getIfscDetails(bank.ifsc)
        if (response.data && response.data.data) {
          bank.bankName = response.data.data.BANK
          bank.branch = response.data.data.BRANCH
          toast.success(`Fetched branch: ${response.data.data.BRANCH}`)
        }
      } catch (err) {
        // Silently ignore
      }
    }
  })
}, { deep: true })

const eighteenYearsAgo = computed(() => {
  const d = new Date()
  d.setFullYear(d.getFullYear() - 18)
  return d.toISOString().split('T')[0]
})

const initials = computed(() => {
  const name = authStore.state.user?.name || 'U'
  return name.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()
})

const populateForm = (user) => {
  if (!user) return
  kycForm.mobile = user.mobile ? user.mobile.replace(/^\+91\s*/, '') : ''
  kycForm.dob = user.dob || ''
  kycForm.address = user.address || ''
  kycForm.fatherName = user.fatherName || user.father_name || ''
  kycForm.motherName = user.motherName || user.mother_name || ''
  kycForm.country = user.country || ''
  kycForm.state = user.state || ''
  kycForm.city = user.city || ''
  kycForm.pincode = user.pincode || ''
  kycForm.incomeRange = user.income_range || user.incomeRange || ''
  kycForm.pan = user.pan || ''
  kycForm.aadhaar = user.aadhaar || ''
  
  if (user.bank_accounts && user.bank_accounts.length > 0) {
     kycForm.bankAccounts = user.bank_accounts.map(b => ({
       accountType: b.account_type,
       bankName: b.bank_name,
       accountNumber: b.account_number,
       ifsc: b.ifsc,
       branch: b.branch,
       _lastIfsc: b.ifsc
     }))
  }
  
  if (user.nominees && user.nominees.length > 0) {
    kycForm.nomineeEnabled = true
    kycForm.nominees = user.nominees.map(n => ({
      name: n.name,
      dob: n.dob,
      pan: n.pan,
      relationship: n.relationship,
      percentage: n.percentage,
      guardianName: n.guardian_name,
      guardianDob: n.guardian_dob,
      guardianRelationship: n.guardian_relationship,
      guardianPan: n.guardian_pan
    }))
  }
}

watch(() => authStore.state.user, (newVal) => {
  populateForm(newVal)
}, { immediate: true })

onMounted(() => {
  if (authStore.state.user) {
    populateForm(authStore.state.user)
  }
})

const addSecondaryBank = () => {
  if (kycForm.bankAccounts.length < 2) {
    kycForm.bankAccounts.push({
      accountType: 'secondary',
      bankName: '',
      accountNumber: '',
      ifsc: '',
      branch: '',
      _lastIfsc: ''
    })
  }
}

const removeSecondaryBank = (index) => {
  kycForm.bankAccounts.splice(index, 1)
}

const isNomineeMinor = (index) => {
  const dobStr = kycForm.nominees[index].dob
  if (!dobStr) return false
  const dob = new Date(dobStr)
  const today = new Date()
  let age = today.getFullYear() - dob.getFullYear()
  const m = today.getMonth() - dob.getMonth()
  if (m < 0 || (m === 0 && today.getDate() < dob.getDate())) {
    age--
  }
  return age < 18
}

const formatBackendError = (errObj) => {
  if (typeof errObj === 'string') {
    const match = errObj.match(/Field validation for '([^']+)' failed on the '([^']+)' tag/)
    if (match) {
      return `${match[1]} is required or invalid (${match[2]}).`
    }
    return errObj
  }
  return 'Invalid input parameters'
}

// Ensure SEBI valid PAN format
const validatePAN = (pan) => {
  const regex = /^[A-Z]{5}[0-9]{4}[A-Z]{1}$/
  return regex.test(pan)
}

const saveProfile = async () => {
  editingField.value = null // Close any open edit fields
  
  if (kycForm.nomineeEnabled) {
    if (totalNomineePercentage.value !== 100) {
      toast.error('Total nominee percentage must equal exactly 100%')
      return
    }
    for (let i = 0; i < kycForm.nominees.length; i++) {
      const nom = kycForm.nominees[i]
      if (nom.pan && !validatePAN(nom.pan)) {
        toast.error(`Invalid PAN format for Nominee ${i+1}. Must be ABCDE1234F`)
        return
      }
      if (isNomineeMinor(i) && nom.guardianPan && !validatePAN(nom.guardianPan)) {
        toast.error(`Invalid Guardian PAN format for Nominee ${i+1}. Must be ABCDE1234F`)
        return
      }
    }
  }

  loading.value = true
  try {
    let nomineesPayload = []
    if (kycForm.nomineeEnabled) {
      nomineesPayload = kycForm.nominees.map((nom, i) => {
        const payload = {
          name: nom.name,
          dob: nom.dob,
          pan: nom.pan,
          relationship: nom.relationship,
          percentage: Number(nom.percentage)
        }
        if (isNomineeMinor(i)) {
          payload.guardian_name = nom.guardianName
          payload.guardian_dob = nom.guardianDob
          payload.guardian_relationship = nom.guardianRelationship
          payload.guardian_pan = nom.guardianPan
        }
        return payload
      })
    }

    const payload = {
      mobile: formatMobile(kycForm.mobile), // Append +91 before saving to backend
      father_name: kycForm.fatherName,
      mother_name: kycForm.motherName,
      date_of_birth: kycForm.dob,
      income_range: kycForm.incomeRange,
      country: kycForm.country,
      state: kycForm.state,
      city: kycForm.city,
      pincode: kycForm.pincode,
      address: kycForm.address,
      bank_accounts: kycForm.bankAccounts
        .filter(b => b.accountNumber || b.ifsc || b.bankName)
        .map(b => ({
          account_type: b.accountType,
          ifsc: b.ifsc,
          bank_name: b.bankName,
          branch: b.branch,
          account_number: b.accountNumber
        })),
      nominees: nomineesPayload
    }

    await userApi.updateKyc(payload)
    
    // Fetch updated user and update store/localStorage
    const userRes = await userApi.getProfile()
    if (userRes.data && userRes.data.data) {
      authStore.login(userRes.data.data, authStore.state.token)
    }
    
    toast.success('Account settings saved successfully!')
  } catch (err) {
    const apiError = err.response?.data?.error
    const apiMessage = err.response?.data?.message
    
    if (apiError && typeof apiError === 'string') {
      toast.error(formatBackendError(apiError))
    } else {
      toast.error(apiMessage || 'Failed to submit account settings')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.flex-row-between {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.cursor-pointer {
  cursor: pointer;
  padding: 4px;
}
.cursor-pointer:hover {
  opacity: 0.8;
}
.text-green { color: #10b981; }
.text-red { color: #f43f5e; }
.mb-2 { margin-bottom: 12px; }

.static-value {
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.02);
  border-radius: var(--border-radius-md);
  color: var(--text-main);
  border: 1px solid transparent;
  font-weight: 500;
}

.input-with-prefix {
  position: relative;
  display: flex;
  align-items: center;
}
.prefix {
  position: absolute;
  left: 14px;
  color: var(--text-muted);
  font-weight: 600;
  z-index: 10;
}
.prefix-input {
  padding-left: 50px !important;
}

/* Inherit existing form-group styles from user-layout and common CSS */
.section-title {
  margin-top: 24px;
  scroll-margin-top: 80px; /* Offset for sticky header */
}
.section-title h3 {
  font-size: var(--font-size-lg);
  color: var(--primary-color);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.section-divider {
  border: 0;
  height: 1px;
  background: var(--border-color);
  margin-bottom: 16px;
}
.title-with-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.btn-outline-small {
  background: transparent;
  border: 1px solid var(--primary-color);
  color: var(--primary-color);
  padding: 6px 12px;
  border-radius: var(--border-radius-md);
  font-size: var(--font-size-xs);
  cursor: pointer;
  transition: all var(--transition-fast);
}
.btn-outline-small:hover {
  background: var(--primary-dim);
}
.bank-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--surface-elevated, #1e2740);
  padding: 10px 16px;
  border-radius: var(--border-radius-md);
  margin-bottom: 12px;
}
.bank-header h4 {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--text-main);
}
.btn-text-red {
  background: none;
  border: none;
  color: var(--loss-color);
  font-size: var(--font-size-xs);
  cursor: pointer;
}
.btn-text-red:hover {
  text-decoration: underline;
}
.help-text {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 4px;
  display: block;
}
.mt-4 { margin-top: 32px; }
.mt-2 { margin-top: 16px; }
.text-right { text-align: right; }

/* Toggle Switch CSS */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}
.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}
.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--surface-elevated, #1e2740);
  transition: .4s;
}
.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: .4s;
}
input:checked + .slider {
  background-color: var(--primary-color);
}
input:checked + .slider:before {
  transform: translateX(20px);
}
.slider.round {
  border-radius: 34px;
}
.slider.round:before {
  border-radius: 50%;
}
.muted-text {
  color: var(--text-muted);
}
.text-sm {
  font-size: var(--font-size-sm);
}

.guardian-banner {
  background-color: rgba(245, 158, 11, 0.1);
  border: 1px solid var(--yellow-border, rgba(245, 158, 11, 0.3));
  color: var(--yellow-color, #f59e0b);
  padding: 12px;
  border-radius: var(--border-radius-md);
  font-size: var(--font-size-sm);
  margin-bottom: 12px;
}
</style>
