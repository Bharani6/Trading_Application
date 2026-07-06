<template>
  <div class="landing-page">
    <!-- Navbar -->
    <nav class="navbar">
      <div class="logo">
        <img src="../../assets/icons/favicon.svg" alt="MoneymakePro Logo" class="brand-logo" /> MoneymakePro
      </div>

      <div class="nav-auth">
        <router-link to="/login" class="btn btn-outline">Log In</router-link>
        <router-link to="/register" class="btn btn-primary">Sign Up</router-link>
      </div>
    </nav>

    <!-- Hero Section -->
    <header class="hero">
      <div class="hero-glow">
        <svg class="hero-bg-svg" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1440 800" preserveAspectRatio="none">
          <path fill="none" stroke="#10B981" stroke-width="3" stroke-opacity="0.4" d="M-100 400 Q 200 100, 600 400 T 1500 400" filter="url(#glow)"></path>
          <path fill="none" stroke="#00F0FF" stroke-width="2" stroke-opacity="0.3" d="M-100 450 Q 300 200, 700 450 T 1500 450"></path>
          <path fill="none" stroke="#10B981" stroke-width="1" stroke-opacity="0.2" d="M-100 500 Q 400 300, 800 500 T 1500 500"></path>
          <!-- Glow Filter -->
          <defs>
            <filter id="glow" x="-20%" y="-20%" width="140%" height="140%">
              <feGaussianBlur stdDeviation="3" result="blur" />
              <feComposite in="SourceGraphic" in2="blur" operator="over" />
            </filter>
          </defs>
          <!-- Geometric glowing dots -->
          <circle cx="250" cy="300" r="4" fill="#10B981" filter="url(#glow)"><animate attributeName="opacity" values="0.2;1;0.2" dur="3s" repeatCount="indefinite"/></circle>
          <circle cx="750" cy="400" r="3" fill="#00F0FF" filter="url(#glow)"><animate attributeName="opacity" values="0.1;1;0.1" dur="4s" repeatCount="indefinite"/></circle>
          <circle cx="1150" cy="350" r="5" fill="#10B981" filter="url(#glow)"><animate attributeName="opacity" values="0.3;1;0.3" dur="2s" repeatCount="indefinite"/></circle>
        </svg>
      </div>
      <h1 class="hero-title">Features Designed for<br/>Your Profit Growth</h1>
      <p class="hero-subtitle">Experience lightning-fast trades, deep liquidity, and institutional-grade analytics.</p>
      
      <!-- 3D Isometric Centerpiece -->
      <div class="isometric-container">
        <!-- Left floating card -->
        <div class="iso-card card-left">
          <div class="iso-content">
            <div class="iso-row"><span class="label">HDFC Bank</span><span class="value green">+2.4%</span></div>
            <div class="iso-row"><span class="label">Cupid</span><span class="value green">+22.4%</span></div>
            <div class="iso-row"><span class="label">Ashok Leyland</span><span class="value red">-1.2%</span></div>
            <div class="iso-row"><span class="label">Gold ETF</span><span class="value green">+0.8%</span></div>
            <div class="iso-chart"></div>
          </div>
        </div>
        
        <!-- Right floating card -->
        <div class="iso-card card-right">
          <div class="iso-content">
            <div class="iso-header">Portfolio Value</div>
            <div class="iso-balance">$142,850.00</div>
            <div class="iso-line-chart"></div>
          </div>
        </div>

        <!-- Center Glowing Logo Base -->
        <div class="iso-base">
          <div class="iso-logo"><i class="fas fa-chart-line"></i></div>
        </div>
      </div>
    </header>

    <!-- Stats Section -->
    <section class="stats-section">
      <div class="stats-grid">
        <div class="stat-card" v-reveal>
          <h3>$10B+</h3>
          <p>Daily Volume</p>
        </div>
        <div class="stat-card" v-reveal>
          <h3>1M+</h3>
          <p>Active Investors</p>
        </div>
        <div class="stat-card" v-reveal>
          <h3>0.00ms</h3>
          <p>Trade Latency</p>
        </div>
        <div class="stat-card" v-reveal>
          <h3>Zero</h3>
          <p>Opening Fees</p>
        </div>
      </div>
    </section>

    <!-- 3 Easy Steps Section -->
    <section class="steps-section" id="steps">
      <div class="steps-header" v-reveal>
        <h2>Invest, Trade & Save with MoneymakePro</h2>
        <p class="steps-subtitle">In 3 Easy Steps</p>
      </div>
      
      <div class="steps-container">
        <!-- Step 1 -->
        <div class="step-card" v-reveal>
          <div class="step-number">01</div>
          <div class="step-icon"><i class="fas fa-mobile-alt"></i></div>
          <h3>Enter Mobile Number and verify OTP</h3>
        </div>
        
        <!-- Connecting Line -->
        <div class="step-connector"></div>

        <!-- Step 2 -->
        <div class="step-card" v-reveal>
          <div class="step-number">02</div>
          <div class="step-icon"><i class="fas fa-id-card"></i></div>
          <h3>Verify KYC and Bank details</h3>
        </div>

        <!-- Connecting Line -->
        <div class="step-connector"></div>

        <!-- Step 3 -->
        <div class="step-card" v-reveal>
          <div class="step-number">03</div>
          <div class="step-icon"><i class="fas fa-file-signature"></i></div>
          <h3>Verify and eSign your form</h3>
        </div>
      </div>
    </section>

    <section class="pricing-section" id="pricing">
      <div class="section-header" v-reveal>
        <h2>Transparent Pricing</h2>
        <p class="section-subtitle">No hidden charges. No nasty surprises. Ever.</p>
      </div>
      <div class="pricing-grid">
        <div class="pricing-card" v-reveal>
          <div class="pricing-price">₹0</div>
          <h3>Delivery Trades</h3>
          <p>Invest for the long term without paying a dime in brokerage. Completely free.</p>
        </div>
        <div class="pricing-card featured" v-reveal>
          <div class="pricing-price">₹20</div>
          <h3>Intraday & F&O</h3>
          <p>Flat ₹20 or 0.05% (whichever is lower) per executed order across all segments.</p>
        </div>
      </div>
    </section>

    <!-- Demat Details Section -->
    <section class="demat-section" id="demat-details">
      <div class="demat-header" v-reveal>
        <h2>Details required to<br/>Open a Demat Account</h2>
        <p class="demat-subtitle">Fully digital and paperless onboarding</p>
      </div>

      <div class="demat-grid">
        <!-- Identity Proof -->
        <div class="demat-card" v-reveal>
          <div class="demat-card-img">
            <img src="../../assets/images/pan_proof.png" alt="PAN Card - Identity Proof" />
          </div>
          <div class="demat-content">
            <h3>Identity proof</h3>
            <p>Your PAN is required as identity proof and to fetch your ITR Name</p>
          </div>
        </div>

        <!-- Address Proof -->
        <div class="demat-card" v-reveal>
          <div class="demat-card-img">
            <img src="../../assets/images/aadhaar_proof.png" alt="Aadhaar Card - Address Proof" />
          </div>
          <div class="demat-content">
            <h3>Address proof</h3>
            <p>Your Aadhaar required as address proof. Ensure it is linked to a mobile number</p>
          </div>
        </div>

        <!-- Income Proof -->
        <div class="demat-card" v-reveal>
          <div class="demat-card-img">
            <img src="../../assets/images/income_proof.png" alt="Income Proof Document" />
          </div>
          <div class="demat-content">
            <h3>Income proof</h3>
            <p>Required only if you wish to activate F&O, Commodities segments</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Feature Section 1 -->
    <section class="feature-section left-align" id="features">
      <div class="feature-content" v-reveal>
        <div class="feature-badge">Smart Trading</div>
        <h2>Buy & Sell Stocks With One Tap</h2>
        <p>Place market and limit orders instantly across NSE & BSE. MoneymakePro gives you a fast, clean trading experience so you never miss a market move — whether you're a first-time investor or a seasoned trader.</p>
        <ul class="feature-list">
          <li><i class="fas fa-check-circle"></i> Instant Buy & Sell on NSE / BSE</li>
          <li><i class="fas fa-check-circle"></i> Real-time live stock price updates</li>
          <li><i class="fas fa-check-circle"></i> Full trade history & P&L tracking</li>
        </ul>
      </div>
      <div class="feature-graphic" v-reveal>
        <div class="iso-card float-anim">
          <div class="iso-content">
             <div class="iso-header">Your Portfolio Today</div>
             <div class="iso-balance">₹1,42,850</div>
             <div class="iso-row"><span class="label">HDFC Bank</span><span class="value green">+2.4%</span></div>
             <div class="iso-row"><span class="label">Ashok Leyland</span><span class="value red">-1.2%</span></div>
             <div class="iso-row"><span class="label">Gold ETF</span><span class="value green">+0.8%</span></div>
             <div class="iso-line-chart"></div>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer CTA -->
    <footer class="footer-cta" v-reveal>
      <div class="cta-container">
        <h2>Ready to Start Your Journey?</h2>
        <p>Join thousands of investors growing their wealth with MoneymakePro.</p>
        <router-link to="/register" class="btn btn-primary btn-lg">Open an Account</router-link>
      </div>
      <div class="footer-bottom">
        <div class="logo"><img src="../../assets/icons/favicon.svg" alt="MoneymakePro Logo" class="brand-logo" /> MoneymakePro</div>
        <div class="footer-links">
          <a href="#" @click.prevent="openDisclaimer('Terms of Service')">Terms of Service</a>
          <a href="#" @click.prevent="openDisclaimer('Privacy Policy')">Privacy Policy</a>
          <a href="#" @click.prevent="showSupportModal = true">Support</a>
        </div>
      </div>

      <div class="footer-copyright">
        &copy; 2026 MoneymakePro. All Rights Reserved. MoneymakePro is an online brand of MoneymakePro.
      </div>
    </footer>

    <SupportModal :isOpen="showSupportModal" @close="showSupportModal = false" />
    <DisclaimerModal :isOpen="showDisclaimerModal" :title="disclaimerTitle" @close="showDisclaimerModal = false" />
  </div>
</template>

<script setup>
import '../../assets/css/landing.css'
import { ref, onMounted } from 'vue'
import SupportModal from '../common/SupportModal.vue'
import DisclaimerModal from '../common/DisclaimerModal.vue'

const showSupportModal = ref(false)
const showDisclaimerModal = ref(false)
const disclaimerTitle = ref('Disclaimer')

const openDisclaimer = (title) => {
  disclaimerTitle.value = title
  showDisclaimerModal.value = true
}

const vReveal = {
  mounted: (el) => {
    el.classList.add('reveal-item')
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          entry.target.classList.add('reveal-active')
          observer.unobserve(entry.target)
        }
      })
    }, { threshold: 0.15 })
    observer.observe(el)
  }
}
</script>
