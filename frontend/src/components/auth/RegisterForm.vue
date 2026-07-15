<template>
  <div class="register-split-page">
    <!-- Animated Background Particles -->
    <div class="bg-particles">
      <div class="particle p1"></div><div class="particle p2"></div><div class="particle p3"></div><div class="particle p4"></div>
      <div class="particle p5"></div><div class="particle p6"></div><div class="particle p7"></div><div class="particle p8"></div>
    </div>

    <div class="right-panel">
      <div class="auth-card" style="max-width: 600px; width: 100%;">
        <div class="auth-card-header">
          <div style="position: relative;">
            <button type="button" class="btn btn-outline" style="position: absolute; left: 0; top: 0; width: 40px; height: 40px; display: flex; justify-content: center; align-items: center; padding: 0; z-index: 10;" @click="goBack">
              <i class="fas fa-arrow-left"></i>
            </button>
            <div class="brand-header" style="justify-content: center; width: 100%;">
              <img src="../../assets/icons/favicon.svg" alt="MoneymakePro" class="auth-brand-logo" />
              <span class="auth-brand-name">MoneymakePro</span>
            </div>
          </div>
          <h1 class="auth-title" style="text-align: center; margin-top: 15px;">Create Account</h1>
          
          <!-- Stepper -->
          <div class="wizard-stepper">
            <div class="step" :class="{ active: currentStep === 1, completed: currentStep > 1 }">1. Personal</div>
            <div class="step-divider"></div>
            <div class="step" :class="{ active: currentStep === 2, completed: currentStep > 2 }">2. Identity Proof</div>
            <div class="step-divider"></div>
            <div class="step" :class="{ active: currentStep === 3 || currentStep === 4, completed: currentStep > 4 }">3. Bank</div>
            <div class="step-divider"></div>
            <div class="step" :class="{ active: currentStep === 5, completed: currentStep > 5 }">4. IPV</div>
          </div>
        </div>

        <form @submit.prevent="handleNextOrSubmit" class="auth-form" novalidate>
          <div v-if="errorMessage" class="error-banner">
            <i class="fa fa-exclamation-circle"></i>
            {{ errorMessage }}
          </div>

          <!-- STEP 1: Personal Information -->
          <div v-if="currentStep === 1" class="form-grid">
            <div class="form-group">
              <label class="form-label">Full name <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input type="text" class="form-input" v-model="form.name" placeholder="John Doe" autocomplete="name" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Email address <span class="required-star">*</span></label>
              <div class="input-wrapper" style="display: flex; gap: 10px;">
                <input type="email" class="form-input" v-model="form.email" placeholder="you@example.com" autocomplete="email" :disabled="emailOtpVerified" />
                <button type="button" class="submit-btn" style="white-space: nowrap; padding: 0 15px; width: auto;" @click="sendEmailOtp" :disabled="emailOtpVerified || isSendingEmailOtp || !form.email">
                  <span v-if="isSendingEmailOtp"><i class="fa fa-spinner fa-spin"></i></span>
                  <span v-else>{{ emailOtpSent ? 'Resend OTP' : 'Send OTP' }}</span>
                </button>
              </div>
            </div>

            <div class="form-group full-width" v-if="emailOtpSent && !emailOtpVerified">
              <label class="form-label">Enter Email OTP <span class="required-star">*</span></label>
              <div class="input-wrapper" style="display: flex; gap: 10px; max-width: 300px;">
                <input type="text" class="form-input" v-model="form.emailOtp" placeholder="123456" maxlength="6" @input="form.emailOtp = form.emailOtp.replace(/[^0-9]/g, '')" />
                <button type="button" class="submit-btn" style="white-space: nowrap; padding: 0 15px; width: auto;" @click="verifyEmailOtp" :disabled="isVerifyingEmailOtp || form.emailOtp.length !== 6">
                  <span v-if="isVerifyingEmailOtp"><i class="fa fa-spinner fa-spin"></i></span>
                  <span v-else>Verify</span>
                </button>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Mobile number <span class="required-star">*</span></label>
              <div class="input-wrapper" style="display: flex; gap: 10px;">
                <input type="tel" class="form-input" v-model="form.mobile" placeholder="9876543210" autocomplete="tel" maxlength="10" @input="form.mobile = form.mobile.replace(/[^0-9]/g, '')" :disabled="otpVerified" />
                <button type="button" class="submit-btn" style="white-space: nowrap; padding: 0 15px; width: auto;" @click="sendOtp" :disabled="otpVerified || isSendingOtp || form.mobile.length !== 10">
                  <span v-if="isSendingOtp"><i class="fa fa-spinner fa-spin"></i></span>
                  <span v-else>{{ otpSent ? 'Resend OTP' : 'Send OTP' }}</span>
                </button>
              </div>
            </div>

            <div class="form-group full-width" v-if="otpSent && !otpVerified">
              <label class="form-label">Enter OTP <span class="required-star">*</span></label>
              <div class="input-wrapper" style="display: flex; gap: 10px; max-width: 300px;">
                <input type="text" class="form-input" v-model="form.otp" placeholder="123456" maxlength="6" @input="form.otp = form.otp.replace(/[^0-9]/g, '')" />
                <button type="button" class="submit-btn" style="white-space: nowrap; padding: 0 15px; width: auto;" @click="verifyOtp" :disabled="isVerifyingOtp || form.otp.length !== 6">
                  <span v-if="isVerifyingOtp"><i class="fa fa-spinner fa-spin"></i></span>
                  <span v-else>Verify</span>
                </button>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Date of birth <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input type="date" class="form-input" v-model="form.dob" :max="maxDate" style="cursor: pointer" @click="openDatePicker" @focus="openDatePicker" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Password <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input :type="showPassword ? 'text' : 'password'" class="form-input" v-model="form.password" placeholder="Min. 8 characters" autocomplete="new-password" />
                <i :class="showPassword ? 'fa fa-eye-slash' : 'fa fa-eye'" class="password-toggle-icon" @click="showPassword = !showPassword"></i>
              </div>
              <span class="password-hint">Use at least 8 characters with uppercase, lowercase, number & special character</span>
            </div>

            <div class="form-group">
              <label class="form-label">Re-enter Password <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input :type="showConfirmPassword ? 'text' : 'password'" class="form-input" v-model="form.confirm_password" placeholder="Re-enter Password" autocomplete="new-password" />
                <i :class="showConfirmPassword ? 'fa fa-eye-slash' : 'fa fa-eye'" class="password-toggle-icon" @click="showConfirmPassword = !showConfirmPassword"></i>
              </div>
            </div>

            <div class="form-group full-width">
              <label class="form-label">Residential address <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <textarea class="form-input" v-model="form.address" placeholder="Enter your full residential address" rows="2"></textarea>
              </div>
            </div>
            
            <div class="form-group full-width">
              <label class="form-label">Income range <span class="required-star">*</span></label>
              <div class="income-chips">
                <button type="button" v-for="range in ['< ₹1L', '₹1L - ₹5L', '₹5L - ₹10L', '> ₹10L']" :key="range" :class="['income-chip', { active: form.income_range === range }]" @click="form.income_range = range">
                  {{ range }}
                </button>
              </div>
            </div>
            
            <div class="form-group full-width">
              <label class="form-label">Occupation <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <select class="form-input" v-model="form.occupation">
                  <option value="" disabled>Select Occupation</option>
                  <option value="Private Sector">Private Sector</option>
                  <option value="Public Sector">Public Sector</option>
                  <option value="Government Service">Government Service</option>
                  <option value="Business Professional">Business Professional</option>
                  <option value="Student">Student</option>
                  <option value="Others">Others</option>
                </select>
              </div>
            </div>
          </div>

          <!-- STEP 2: Identity Details -->
          <div v-if="currentStep === 2" class="form-grid">
            <div class="form-group full-width">
              <label class="form-label">PAN Number <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input type="text" class="form-input" v-model="form.pan" placeholder="ABCDE1234F" style="text-transform: uppercase" maxlength="10" @input="form.pan = form.pan.toUpperCase().replace(/[^A-Z0-9]/g, '')" />
              </div>
              <span v-if="panError" class="error-text" style="color: #f43f5e; font-size: 12px; margin-top: 4px; display: block;">{{ panError }}</span>
            </div>
            
            <div class="form-group full-width">
              <label class="form-label">Aadhaar Number <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input type="text" class="form-input" v-model="form.aadhaar" placeholder="12 Digit Aadhaar" maxlength="12" @input="form.aadhaar = form.aadhaar.replace(/[^0-9]/g, '')" />
              </div>
              <span v-if="aadhaarError" class="error-text" style="color: #f43f5e; font-size: 12px; margin-top: 4px; display: block;">{{ aadhaarError }}</span>
            </div>
          </div>

          <!-- STEP 3: Bank Details -->
          <div v-if="currentStep === 3" class="form-grid">
            <div class="form-group full-width">
              <label class="form-label">Account Holder Name <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input type="text" class="form-input" v-model="form.bank_accounts[0].account_holder_name" placeholder="Name as per bank" />
              </div>
            </div>
            
            <div class="form-group">
              <label class="form-label">IFSC Code <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input type="text" class="form-input" v-model="form.bank_accounts[0].ifsc" placeholder="HDFC0001234" style="text-transform: uppercase" maxlength="11" @input="handleIfscInput" />
              </div>
            </div>
            
            <div class="form-group">
              <label class="form-label">Bank Name <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input type="text" class="form-input" v-model="form.bank_accounts[0].bank_name" placeholder="Auto-fetched or manual" />
              </div>
            </div>
            
            <div class="form-group full-width">
              <label class="form-label">Branch <span class="text-muted">(Optional)</span></label>
              <div class="input-wrapper">
                <input type="text" class="form-input" v-model="form.bank_accounts[0].branch" placeholder="Auto-fetched or manual" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Account Number <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input type="password" class="form-input" v-model="form.bank_accounts[0].account_number" placeholder="Enter Account Number" @input="form.bank_accounts[0].account_number = form.bank_accounts[0].account_number.replace(/[^0-9]/g, '')" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">Confirm Account Number <span class="required-star">*</span></label>
              <div class="input-wrapper">
                <input type="text" class="form-input" v-model="form.bank_accounts[0].confirm_account_number" placeholder="Re-enter Account Number" @input="form.bank_accounts[0].confirm_account_number = form.bank_accounts[0].confirm_account_number.replace(/[^0-9]/g, '')" />
              </div>
            </div>
          </div>

          <!-- STEP 4: Nominee Details -->
          <div v-if="currentStep === 4" class="form-grid">
            <div class="form-group full-width">
              <label class="form-label">Do you want to add a nominee?</label>
              <div class="income-chips" style="justify-content: flex-start; gap: 10px;">
                <button type="button" class="income-chip" :class="{ active: form.addNominee }" @click="form.addNominee = true">Yes</button>
                <button type="button" class="income-chip" :class="{ active: !form.addNominee }" @click="form.addNominee = false">No</button>
              </div>
            </div>
            
            <template v-if="form.addNominee">
              <div class="form-group">
                <label class="form-label">Nominee Name <span class="required-star">*</span></label>
                <div class="input-wrapper">
                  <input type="text" class="form-input" v-model="form.nominees[0].name" placeholder="Nominee Full Name" />
                </div>
              </div>
              
              <div class="form-group">
                <label class="form-label">Relationship <span class="required-star">*</span></label>
                <div class="input-wrapper">
                  <select class="form-input" v-model="form.nominees[0].relationship">
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
              </div>
              
              <div class="form-group">
                <label class="form-label">Date of Birth <span class="required-star">*</span></label>
                <div class="input-wrapper">
                  <input type="date" class="form-input" v-model="form.nominees[0].dob" @click="openDatePicker" @focus="openDatePicker" style="cursor: pointer" />
                </div>
              </div>
              
              <div class="form-group">
                <label class="form-label">Phone Number <span class="text-muted">(Optional)</span></label>
                <div class="input-wrapper">
                  <input type="tel" class="form-input" v-model="form.nominees[0].phone" placeholder="9876543210" maxlength="10" @input="form.nominees[0].phone = form.nominees[0].phone.replace(/[^0-9]/g, '')" />
                </div>
              </div>

              <!-- Guardian Details if Minor -->
              <template v-if="form.nominees[0].dob && calculateAge(form.nominees[0].dob) < 18">
                <div class="form-group full-width" style="margin-top: 10px;">
                  <h4 style="color: var(--primary); font-size: 14px; margin-bottom: 5px;"><i class="fas fa-shield-alt mr-2"></i> Guardian Details (Required for Minor)</h4>
                  <p style="font-size: 12px; color: var(--text-muted);">Since the nominee is below 18 years of age, guardian details are mandatory.</p>
                </div>
                
                <div class="form-group">
                  <label class="form-label">Guardian Name <span class="required-star">*</span></label>
                  <div class="input-wrapper">
                    <input type="text" class="form-input" v-model="form.nominees[0].guardian_name" placeholder="Guardian Full Name" />
                  </div>
                </div>
                
                <div class="form-group">
                  <label class="form-label">Guardian DOB <span class="required-star">*</span></label>
                  <div class="input-wrapper">
                    <input type="date" class="form-input" v-model="form.nominees[0].guardian_dob" :max="maxDate" @click="openDatePicker" @focus="openDatePicker" style="cursor: pointer" />
                  </div>
                </div>
              </template>
            </template>
          </div>

          <!-- STEP 5: IPV (In-Person Verification) -->
          <div v-if="currentStep === 5" class="form-grid">
            <div class="form-group full-width">
              <label class="form-label">In-Person Verification <span class="required-star">*</span></label>
              <p style="font-size: 13px; color: var(--text-muted); margin-bottom: 10px;">Please allow camera and location access to complete your KYC. Ensure your face is clearly visible.</p>
              
              <div v-if="!form.ipv_photo" class="camera-container" style="display: flex; flex-direction: column; align-items: center; background: rgba(0,0,0,0.2); padding: 15px; border-radius: 8px;">
                <video ref="videoElement" autoplay playsinline style="width: 100%; max-width: 320px; border-radius: 8px; margin-bottom: 15px; background: #000; box-shadow: 0 4px 12px rgba(0,0,0,0.2);"></video>
                <button type="button" style="background: linear-gradient(135deg, #3b82f6, #2563eb); color: #ffffff; border: none; border-radius: 8px; padding: 12px 24px; font-weight: 600; font-size: 14px; cursor: pointer; transition: all 0.3s ease; box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3); display: flex; align-items: center; gap: 8px;" @click="capturePhotoAndLocation" :disabled="isCapturing" onmouseover="this.style.transform='translateY(-2px)'; this.style.boxShadow='0 6px 16px rgba(59, 130, 246, 0.4)';" onmouseout="this.style.transform='translateY(0)'; this.style.boxShadow='0 4px 12px rgba(59, 130, 246, 0.3)';">
                  <i class="fas fa-camera"></i> {{ isCapturing ? 'Capturing...' : 'Capture Photo & Location' }}
                </button>
                <canvas ref="canvasElement" style="display: none;"></canvas>
              </div>

              <div v-else class="preview-container" style="display: flex; flex-direction: column; align-items: center; background: rgba(0,0,0,0.2); padding: 15px; border-radius: 8px;">
                <img :src="form.ipv_photo" alt="Captured IPV" style="width: 100%; max-width: 320px; border-radius: 8px; margin-bottom: 15px; border: 2px solid #10b981;" />
                <div style="font-size: 12px; color: var(--text-muted); margin-bottom: 10px; text-align: center;">
                  <i class="fas fa-map-marker-alt text-green mr-1"></i> Location: {{ form.ipv_latitude }}, {{ form.ipv_longitude }}
                </div>
                <button type="button" class="btn-text-red" @click="retakePhoto">
                  <i class="fas fa-redo mr-2"></i> Retake Photo
                </button>
              </div>
            </div>
          </div>

          <p class="terms-note" v-if="currentStep === 5">
            By creating an account, you agree to our
            <a href="#">Terms of Service</a> and <a href="#">Privacy Policy</a>.
          </p>

          <div style="display: flex; gap: 15px; margin-top: 20px;">
            <button type="submit" class="submit-btn" style="flex: 1;" :disabled="loading">
              <span v-if="!loading">{{ currentStep < 5 ? 'Continue' : 'Create Account' }}</span>
              <i v-else class="fa fa-spinner fa-spin"></i>
            </button>
          </div>
        </form>

        <div class="login-prompt">
          Already have an account? 
          <router-link to="/login" class="login-link">Sign in</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import '../../assets/css/register.css'
import { reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import { authApi } from '../../api/auth.api'
import { utilsApi } from '../../api/utils.api'

const router = useRouter()
const toast = useToast()

const currentStep = ref(1)
const loading = ref(false)
const errorMessage = ref('')
const panError = ref('')
const aadhaarError = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)

const otpSent = ref(false)
const otpVerified = ref(false)
const isSendingOtp = ref(false)
const isVerifyingOtp = ref(false)

const emailOtpSent = ref(false)
const emailOtpVerified = ref(false)
const isSendingEmailOtp = ref(false)
const isVerifyingEmailOtp = ref(false)

const videoElement = ref(null)
const canvasElement = ref(null)
const mediaStream = ref(null)
const isCapturing = ref(false)

const form = reactive({
  name: '',
  email: '',
  emailOtp: '',
  mobile: '',
  otp: '',
  dob: '',
  pan: '',
  aadhaar: '',
  income_range: '',
  occupation: '',
  ipv_photo: '',
  ipv_latitude: '',
  ipv_longitude: '',
  password: '',
  confirm_password: '',
  address: '',
  bank_accounts: [{
    account_type: 'savings',
    ifsc: '',
    bank_name: '',
    branch: '',
    account_number: '',
    confirm_account_number: '',
    account_holder_name: ''
  }],
  addNominee: false,
  nominees: [{
    name: '',
    dob: '',
    pan: '',
    relationship: '',
    percentage: 100,
    guardian_name: '',
    guardian_dob: '',
    guardian_relationship: '',
    guardian_pan: ''
  }]
})

const today = new Date()
const maxDate = new Date(today.getFullYear() - 18, today.getMonth(), today.getDate()).toISOString().split('T')[0]

const openDatePicker = (e) => {
  try {
    e.target.showPicker()
  } catch (err) {
    // browser doesn't support showPicker, fallback to default behavior
  }
}

const showError = (msg) => {
  errorMessage.value = msg
  window.scrollTo({ top: 0, behavior: 'smooth' })
  loading.value = false
}

const calculateAge = (dobString) => {
  if (!dobString) return 0
  const dob = new Date(dobString)
  const today = new Date()
  let age = today.getFullYear() - dob.getFullYear()
  const m = today.getMonth() - dob.getMonth()
  if (m < 0 || (m === 0 && today.getDate() < dob.getDate())) {
    age--
  }
  return age
}

const sendOtp = async () => {
  if (form.mobile.length !== 10) return
  isSendingOtp.value = true
  errorMessage.value = ''
  try {
    const res = await authApi.sendOTP({ mobile: form.mobile })
    if (res.data.success) {
      toast.success('OTP sent to mobile number.')
      otpSent.value = true
    }
  } catch (error) {
    errorMessage.value = error.response?.data?.error ? formatBackendError(error.response.data.error) : 'Failed to send OTP.'
    toast.error(errorMessage.value)
  } finally {
    isSendingOtp.value = false
  }
}

const verifyOtp = async () => {
  if (form.otp.length !== 6) return
  isVerifyingOtp.value = true
  errorMessage.value = ''
  try {
    const res = await authApi.verifyOTP({ mobile: form.mobile, otp: form.otp })
    if (res.data.success) {
      toast.success('Mobile number verified.')
      otpVerified.value = true
    }
  } catch (error) {
    errorMessage.value = error.response?.data?.error ? formatBackendError(error.response.data.error) : 'Invalid OTP.'
    toast.error(errorMessage.value)
  } finally {
    isVerifyingOtp.value = false
  }
}

const sendEmailOtp = async () => {
  if (!form.email) return
  isSendingEmailOtp.value = true
  errorMessage.value = ''
  try {
    const res = await authApi.sendEmailOTP({ email: form.email })
    if (res.data.success) {
      toast.success('OTP sent to email address.')
      emailOtpSent.value = true
    }
  } catch (error) {
    errorMessage.value = error.response?.data?.error ? formatBackendError(error.response.data.error) : 'Failed to send OTP.'
    toast.error(errorMessage.value)
  } finally {
    isSendingEmailOtp.value = false
  }
}

const verifyEmailOtp = async () => {
  if (form.emailOtp.length !== 6) return
  isVerifyingEmailOtp.value = true
  errorMessage.value = ''
  try {
    const res = await authApi.verifyEmailOTP({ email: form.email, otp: form.emailOtp })
    if (res.data.success) {
      toast.success('Email address verified.')
      emailOtpVerified.value = true
    }
  } catch (error) {
    errorMessage.value = error.response?.data?.error ? formatBackendError(error.response.data.error) : 'Invalid OTP.'
    toast.error(errorMessage.value)
  } finally {
    isVerifyingEmailOtp.value = false
  }
}

let ifscTimeout = null
const handleIfscInput = () => {
  const bank = form.bank_accounts[0]
  bank.ifsc = bank.ifsc.toUpperCase().replace(/[^A-Z0-9]/g, '')
  if (ifscTimeout) clearTimeout(ifscTimeout)
  if (bank.ifsc.length === 11) {
    ifscTimeout = setTimeout(async () => {
      try {
        const response = await utilsApi.getIfscDetails(bank.ifsc)
        if (response.data.success && response.data.data) {
          bank.bank_name = response.data.data.BANK || ''
          bank.branch = response.data.data.BRANCH || ''
        }
      } catch (err) {
        // Silently fail, user can enter manually
      }
    }, 500)
  }
}

const validateStep = (step) => {
  errorMessage.value = ''
  if (step === 1) {
    if (!form.name) return 'Full name is required.'
    if (!/^[a-zA-Z\s]+$/.test(form.name)) return 'Name must contain only letters and spaces.'
    if (!form.email) return 'Email address is required.'
    if (!emailOtpVerified.value) return 'Please verify your email address with OTP.'
    if (!form.mobile || form.mobile.length !== 10) return 'Valid 10-digit mobile number is required.'
    if (!form.dob) return 'Date of birth is required.'
    if (calculateAge(form.dob) < 18) return 'You must be at least 18 years old.'
    if (!otpVerified.value) return 'Please verify your mobile number with OTP.'
    
    if (!form.password) return 'Password is required.'
    if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}/.test(form.password)) {
      return 'Password must be at least 8 chars long with uppercase, lowercase, numeric, and special characters.'
    }
    if (form.password !== form.confirm_password) return 'Passwords do not match.'

    if (!form.address) return 'Residential address is required.'
    if (!form.income_range) return 'Income range is required.'
    if (!form.occupation) return 'Occupation is required.'
    return true
  }
  
  if (step === 2) {
    panError.value = ''
    aadhaarError.value = ''
    let valid = true

    if (!form.pan) {
      panError.value = 'PAN Number is required.'
      valid = false
    } else if (!/^[A-Z]{5}[0-9]{4}[A-Z]{1}$/.test(form.pan)) {
      panError.value = 'PAN format must be ABCDE1234F'
      valid = false
    }

    if (!form.aadhaar || form.aadhaar.length !== 12) {
      aadhaarError.value = '12-digit Aadhaar Number is required.'
      valid = false
    }

    if (!valid) return false
    return true
  }
  
  if (step === 3) {
    const bank = form.bank_accounts[0]
    if (!bank.account_holder_name) return 'Account Holder Name is required.'
    if (!bank.ifsc || bank.ifsc.length !== 11) return 'Valid 11-character IFSC code is required.'
    if (!bank.bank_name) return 'Bank Name is required.'
    if (!bank.account_number) return 'Account Number is required.'
    if (bank.account_number !== bank.confirm_account_number) return 'Account Numbers do not match.'
    return true
  }
  
  if (step === 4) {
    if (form.addNominee) {
      const nom = form.nominees[0]
      if (!nom.name) return 'Nominee Name is required.'
      if (!nom.relationship) return 'Nominee Relationship is required.'
      if (!nom.dob) return 'Nominee Date of Birth is required.'
      if (calculateAge(nom.dob) < 18) {
        if (!nom.guardian_name) return 'Guardian Name is required for minor nominee.'
        if (!nom.guardian_dob) return 'Guardian Date of Birth is required.'
        if (calculateAge(nom.guardian_dob) < 18) return 'Guardian must be at least 18 years old.'
      }
    }
    return true
  }
  
  if (step === 5) {
    if (!form.ipv_photo || !form.ipv_latitude || !form.ipv_longitude) {
      return 'Please complete the IPV step by capturing your photo and location.'
    }
    return true
  }
  
  return true
}

const stopCamera = () => {
  if (mediaStream.value) {
    mediaStream.value.getTracks().forEach(track => track.stop())
    mediaStream.value = null
  }
}

const startCamera = async () => {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: 'user' } })
    mediaStream.value = stream
    if (videoElement.value) {
      videoElement.value.srcObject = stream
    }
  } catch (err) {
    toast.error('Camera access denied or unavailable.')
  }
}

watch(currentStep, async (newStep) => {
  if (newStep === 5) {
    await startCamera()
  } else {
    stopCamera()
  }
})

const goBack = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  } else {
    router.push('/login')
  }
}

const capturePhotoAndLocation = () => {
  isCapturing.value = true
  
  if (!navigator.geolocation) {
    isCapturing.value = false
    toast.error("Geolocation is not supported by this browser.")
    return
  }

  navigator.geolocation.getCurrentPosition(
    (position) => {
      form.ipv_latitude = position.coords.latitude.toString()
      form.ipv_longitude = position.coords.longitude.toString()
      
      // Capture Photo
      if (videoElement.value && canvasElement.value) {
        const video = videoElement.value
        const canvas = canvasElement.value
        canvas.width = video.videoWidth
        canvas.height = video.videoHeight
        canvas.getContext('2d').drawImage(video, 0, 0, canvas.width, canvas.height)
        
        // Convert to base64
        form.ipv_photo = canvas.toDataURL('image/jpeg', 0.8)
        stopCamera()
        toast.success("IPV Captured successfully!")
      }
      isCapturing.value = false
    },
    (error) => {
      isCapturing.value = false
      toast.error("Location access denied or unavailable.")
    },
    { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 }
  )
}

const retakePhoto = async () => {
  form.ipv_photo = ''
  form.ipv_latitude = ''
  form.ipv_longitude = ''
  await startCamera()
}

const formatBackendError = (errObj) => {
  if (typeof errObj === 'string') return errObj
  return 'Invalid input parameters'
}

const handleNextOrSubmit = async () => {
  const isValid = validateStep(currentStep.value)
  if (isValid !== true) {
    showError(isValid)
    return
  }

  if (currentStep.value < 5) {
    currentStep.value++
    return
  }

  // Step 5 Complete: Submit Payload
  loading.value = true
  errorMessage.value = ''

  try {
    const payload = {
      name: form.name,
      email: form.email,
      mobile: form.mobile,
      dob: form.dob,
      pan: form.pan.toUpperCase(),
      aadhaar: form.aadhaar,
      income_range: form.income_range,
      occupation: form.occupation,
      ipv_photo: form.ipv_photo,
      ipv_latitude: form.ipv_latitude,
      ipv_longitude: form.ipv_longitude,
      password: form.password,
      address: form.address,
      bank_accounts: [{
        account_type: form.bank_accounts[0].account_type,
        ifsc: form.bank_accounts[0].ifsc.toUpperCase(),
        bank_name: form.bank_accounts[0].bank_name,
        branch: form.bank_accounts[0].branch,
        account_number: form.bank_accounts[0].account_number
      }],
      nominees: form.addNominee ? [{
        name: form.nominees[0].name,
        dob: form.nominees[0].dob,
        relationship: form.nominees[0].relationship,
        percentage: 100,
        guardian_name: form.nominees[0].guardian_name,
        guardian_dob: form.nominees[0].guardian_dob,
      }] : []
    }

    const res = await authApi.register(payload)
    if (res.data.success) {
      toast.success('Registration successful! Please log in.')
      router.push('/login')
    }
  } catch (error) {
    if (error.response?.data) {
      const { message, error: details } = error.response.data
      errorMessage.value = details ? formatBackendError(details) : (message || 'Registration failed.')
    } else {
      errorMessage.value = 'Registration failed. Please try again.'
    }
    toast.error(errorMessage.value)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.wizard-stepper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 1.5rem 0 2rem;
  font-size: 13px;
  font-weight: 500;
}
.step {
  color: var(--text-muted);
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s;
}
.step.active {
  color: var(--primary);
  background: rgba(59, 130, 246, 0.1);
  font-weight: 600;
}
.step.completed {
  color: #10b981;
}
.step-divider {
  flex: 1;
  height: 1px;
  background: rgba(255, 255, 255, 0.1);
  margin: 0 10px;
}
</style>
